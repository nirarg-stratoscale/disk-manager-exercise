package auth

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"
)

const rolesScope = "roles"

// Middleware is middleware for user authorization
// It gets user information from the headers, injects them into the context
// and enforces role based policies according to the definition in the swagger
func Middleware(log logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return &auth{
			Log:   log,
			Inner: inner,
		}
	}
}

type auth struct {
	Log   logrus.FieldLogger
	Inner http.Handler
}

func (a *auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFrom(r)
	if route == nil {
		a.Log.Error("No route matches the request %s %s", r.Method, r.URL)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Extract user from headers
	user := fromHeader(r.Header)

	// check permissions for the route for the user in the request
	if !a.allowedRoles(user, route) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Add the user information to the context, so it can be used in the business logic.
	// In the business logic the user can be extracted with the `FromContext` function.
	r = r.WithContext(user.ToContext(r.Context()))
	a.Inner.ServeHTTP(w, r)
}

// allowedRoles test if a user is allowed to access roles of a route
func (a *auth) allowedRoles(user *User, route *middleware.MatchedRoute) bool {
	var (
		allowedRoles []string
		log          = a.Log.WithFields(logrus.Fields{
			"user-id":   user.ID,
			"operation": route.Operation.ID,
		})
	)

	for _, policy := range route.Authenticators {
		roles := policy.Scopes[rolesScope]
		if len(roles) == 0 {
			log.Debugf("Authorized user: All roles are allowed")
			return true
		}
		for _, role := range roles {
			allowedRoles = append(allowedRoles, role)
			if user.Allowed(role) {
				log.Debugf("Authorized user: role=%q", role)
				return true
			}
		}
	}
	log.Debugf("Unauthorized user, insufficient roles: user-roles=%q allowed-roles=%q", user.Roles, allowedRoles)
	return false
}
