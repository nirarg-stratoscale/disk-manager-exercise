package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

const (
	KeyRoles  = "X-Auth-Roles"
	RoleScope = "roles"
)

func Policy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := middleware.MatchedRouteFrom(r)
		if route != nil {
			for _, v := range route.Authenticators {
				if err := checkRoles(v.Scopes, r.Header); err != nil {
					http.Error(w, err.Error(), http.StatusForbidden)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func checkRoles(scopes map[string][]string, headers http.Header) error {
	routeRoles := scopes[RoleScope]
	if len(routeRoles) == 0 {
		return nil
	}
	routeRolesMap := mapRoles(scopes[RoleScope])
	userRoles := strings.Split(headers.Get(KeyRoles), ",")
	if len(userRoles) == 0 {
		return fmt.Errorf("no roles provided")
	}

	for _, userRole := range userRoles {
		if routeRolesMap[userRole] {
			return nil
		}
	}
	return fmt.Errorf("insufficient roles: %v, need: %v", userRoles, routeRoles)
}

func mapRoles(roles []string) map[string]bool {
	m := make(map[string]bool, len(roles))
	for _, r := range roles {
		m[r] = true
	}
	return m
}
