package wishilist

import (
	"encoding/json"
	"net/http"

	"github.com/wrapped-owls/goremy-di/remy"
)

type Controller struct {
	Injector remy.DependencyRetriever
}

func (h Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uc, err := remy.GetWithContext[UseCase](h.Injector, r.Context())
	if err != nil {
		http.Error(w, "dependency retrieval failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := remy.GetWithContext[User](h.Injector, r.Context())
	if err != nil {
		http.Error(w, "user retrieval failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"user":  user.Username,
		"items": uc.ReadMyWishlist(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
