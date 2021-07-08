package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	// update {REALM} with keycloak.json realm
	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/{REALM}")
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     "", // fill with keycloak.json resource
		ClientSecret: "", // fill with keycloak.json credentials.secret
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:9000/auth/keycloak/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		state := "fake_state"

		http.SetCookie(w, &http.Cookie{
			Name:     "state",
			Value:    state,
			MaxAge:   int(time.Hour.Seconds()),
			Secure:   false,
			HttpOnly: true,
		})

		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/keycloak/callback", func(w http.ResponseWriter, r *http.Request) {
		state, err := r.Cookie("state")
		if err != nil {
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("state") != state.Value {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err != nil {
			http.Error(w, "could not get user info: "+err.Error(), http.StatusInternalServerError)
		}

		resp := struct {
			OAuth2Token *oauth2.Token
			UserInfo    *oidc.UserInfo
		}{token, user}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, string(data))
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
