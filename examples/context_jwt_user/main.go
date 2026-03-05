package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/wrapped-owls/goremy-di/examples/context_jwt_user/internal/middleware"
	"github.com/wrapped-owls/goremy-di/examples/context_jwt_user/internal/middleware/auth"
	"github.com/wrapped-owls/goremy-di/examples/context_jwt_user/internal/wishilist"
	"github.com/wrapped-owls/goremy-di/remy"
)

func NewUserFromContext(ctx context.Context) (wishilist.User, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return wishilist.User{}, err
	}

	userID, _ := claims["userID"].(string)
	username, _ := claims["username"].(string)
	if strings.TrimSpace(userID) == "" {
		return wishilist.User{}, errors.New("userID claim missing or invalid")
	}

	return wishilist.User{ID: userID, Username: username}, nil
}

func main() {
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})

	remy.RegisterConstructor(
		inj, remy.Singleton[wishilist.InMemoryWishlistRepo], wishilist.NewInMemoryWishlistRepo,
	)
	remy.RegisterConstructorArgs1Err(inj, remy.Factory[wishilist.User], NewUserFromContext)
	remy.RegisterConstructorArgs2(
		inj, remy.Factory[wishilist.UseCase], wishilist.NewWishlistUseCase,
	)

	endpoint := middleware.AuthMiddleware(wishilist.Controller{Injector: inj})

	requestWithToken := func(label, token string) {
		req := httptest.NewRequest(http.MethodGet, "http://example.local/wishlist", nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		recorder := httptest.NewRecorder()
		endpoint.ServeHTTP(recorder, req)
		log.Printf(
			"%s -> status=%d body=%s",
			label,
			recorder.Code,
			strings.TrimSpace(recorder.Body.String()),
		)
	}

	aliceToken := auth.BuildClaimsPayload(auth.Claims{
		"userID":   "u-alice",
		"username": "alice",
	})
	bobToken := auth.BuildClaimsPayload(auth.Claims{
		"userID":   "u-bob",
		"username": "bob",
	})

	requestWithToken("alice token", aliceToken)
	requestWithToken("bob token", bobToken)
	requestWithToken("missing token", "")
}
