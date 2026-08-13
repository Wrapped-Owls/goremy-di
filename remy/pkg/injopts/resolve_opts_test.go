package injopts

import "testing"

func TestResolveConfOption_Is(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		opts  ResolveConfOption
		check ResolveConfOption
		want  bool
	}{
		{
			name:  "single flag is set",
			opts:  ResolveOptDuckTyping,
			check: ResolveOptDuckTyping,
			want:  true,
		},
		{
			name:  "another flag is not set",
			opts:  ResolveOptDuckTyping,
			check: ResolveOptIsolated,
			want:  false,
		},
		{
			name:  "combined flags report each one",
			opts:  ResolveOptDuckTyping | ResolveOptTracePath,
			check: ResolveOptTracePath,
			want:  true,
		},
		{
			name:  "combined check needs every flag",
			opts:  ResolveOptDuckTyping,
			check: ResolveOptDuckTyping | ResolveOptTracePath,
			want:  false,
		},
		{
			name:  "none has no flag",
			opts:  ResolveOptNone,
			check: ResolveOptDuckTyping,
			want:  false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := testCase.opts.Is(testCase.check); got != testCase.want {
				t.Fatalf(
					"Is(%b) on %b = %v, want %v",
					testCase.check,
					testCase.opts,
					got,
					testCase.want,
				)
			}
		})
	}
}

// the three flags must occupy distinct bits, or one would imply another
func TestResolveConfOption_FlagsAreDistinct(t *testing.T) {
	t.Parallel()

	flags := []ResolveConfOption{ResolveOptDuckTyping, ResolveOptIsolated, ResolveOptTracePath}
	for outer := range flags {
		for inner := range flags {
			if outer == inner {
				continue
			}
			if flags[outer].Is(flags[inner]) {
				t.Errorf("flag %b overlaps %b", flags[outer], flags[inner])
			}
		}
	}
}
