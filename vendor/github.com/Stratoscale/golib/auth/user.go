package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKeyType string

const (
	contextKey contextKeyType = "user"
	rolesSep                  = ";"
	AdminRole                 = "admin"
)

// HTTP header keys that holds the user information
const (
	KeyUserID    = "X-Auth-User-ID"
	KeyRoles     = "X-Auth-Roles"
	KeyProjectID = "X-Auth-Project-ID"
	KeyDomainID  = "X-Auth-Project-Domain-ID"
)

// FromContext returns the user information from the context.
// If no credentials information exists in the context, a nil will be returned.
// This function should be used in the business logic after the Middleware have
// extracted the user from the headers and injected it into the context.
func FromContext(ctx context.Context) *User {
	user := ctx.Value(contextKey)
	if user == nil {
		return nil
	}
	return user.(*User)
}

// User is the user information
type User struct {
	ID      string
	Project string
	Domain  string
	Roles   []string
}

// Allowed checks if a user is allowed for a certain role
func (u *User) Allowed(wantRole string) bool {
	if u == nil {
		return false
	}
	for _, gotRole := range u.Roles {
		if wantRole == gotRole {
			return true
		}
	}
	return false
}

// toHeaders updates headers with the user infomration
func (u *User) toHeaders(h http.Header) {
	h.Set(KeyUserID, u.ID)
	h.Set(KeyProjectID, u.Project)
	h.Set(KeyDomainID, u.Domain)
	h.Set(KeyRoles, strings.Join(u.Roles, rolesSep))
}

// ToContext sets credentials information the context
func (u *User) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, u)
}

// IsAdmin returns whether the user has an admin role
func (u *User) IsAdmin() bool {
	for _, role := range u.Roles {
		if role == AdminRole {
			return true
		}
	}
	return false
}

// fromHeaders extract user from requests headers
func fromHeader(h http.Header) *User {
	return &User{
		ID:      h.Get(KeyUserID),
		Domain:  h.Get(KeyDomainID),
		Project: h.Get(KeyProjectID),
		Roles:   strings.Split(h.Get(KeyRoles), rolesSep),
	}
}
