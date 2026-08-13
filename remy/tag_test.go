package remy

import (
	"strings"
	"testing"
)

type (
	scopeAuth    struct{}
	scopeBilling struct{}
)

func TestNewTag(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		first    Tag
		second   Tag
		wantSame bool
	}{
		{
			name:     "same name anchored on different scopes never collides",
			first:    NewTag[scopeAuth]("primary"),
			second:   NewTag[scopeBilling]("primary"),
			wantSame: false,
		},
		{
			name:     "same name and scope is stable within the process",
			first:    NewTag[scopeAuth]("primary"),
			second:   NewTag[scopeAuth]("primary"),
			wantSame: true,
		},
		{
			name:     "different names on the same scope differ",
			first:    NewTag[scopeAuth]("primary"),
			second:   NewTag[scopeAuth]("replica"),
			wantSame: false,
		},
		{
			name:     "a scoped tag never equals the bare name",
			first:    NewTag[scopeAuth]("primary"),
			second:   "primary",
			wantSame: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if same := testCase.first == testCase.second; same != testCase.wantSame {
				t.Fatalf(
					"%q == %q is %v, want %v",
					testCase.first,
					testCase.second,
					same,
					testCase.wantSame,
				)
			}
			if !strings.HasPrefix(string(testCase.first), "primary@") {
				t.Errorf("tag %q should keep the readable name as prefix", testCase.first)
			}
		})
	}
}

func TestNewTag_ResolvesIndependently(t *testing.T) {
	t.Parallel()

	inj := NewInjector()
	authPrimary := NewTag[scopeAuth]("primary")
	billingPrimary := NewTag[scopeBilling]("primary")

	RegisterInstance(inj, "auth-db", authPrimary)
	RegisterInstance(inj, "billing-db", billingPrimary)
	// a bare literal must stay a third, independent slot
	RegisterInstance(inj, "legacy-db", "primary")

	testCases := []struct {
		name string
		tag  Tag
		want string
	}{
		{name: "auth scope", tag: authPrimary, want: "auth-db"},
		{name: "billing scope", tag: billingPrimary, want: "billing-db"},
		{name: "untyped literal", tag: "primary", want: "legacy-db"},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := Get[string](inj, testCase.tag)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != testCase.want {
				t.Fatalf("Get = %q, want %q", got, testCase.want)
			}
		})
	}
}
