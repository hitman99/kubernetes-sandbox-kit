package admin

import (
	"encoding/base64"
	"github.com/hitman99/kubernetes-sandbox/internal/config"
	"net/http"
	"strings"
	"sync"
)

type authMiddleware struct {
	adminToken string
	_m         sync.Mutex
}

func NewAuthMiddleware() *authMiddleware {
	cfg, updates := config.Get()
	am := &authMiddleware{
		adminToken: cfg.AdminToken,
		_m:         sync.Mutex{},
	}
	am.syncConfig(updates)
	return am
}

func (a *authMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		parts := strings.Split(strings.TrimSpace(token), " ")
		if len(parts) == 2 && parts[0] == "Bearer" && len(parts[1]) > 0 {
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				http.Error(w, "Access Denied", http.StatusUnauthorized)
				return
			} else {
				a._m.Lock()
				defer a._m.Unlock()
				if string(decoded) != a.adminToken {
					http.Error(w, "Access Denied", http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			}
		} else {
			http.Error(w, "Access Denied", http.StatusUnauthorized)
		}
	})
}

func (a *authMiddleware) syncConfig(updates <-chan config.Config) {
	go func() {
		for c := range updates {
			a._m.Lock()
			a.adminToken = c.AdminToken
			a._m.Unlock()
		}
	}()
}
