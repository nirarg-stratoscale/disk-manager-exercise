package auth

import "net/http"

// RoundTripper sets user information into requests headers
func RoundTripper(user *User, inner http.RoundTripper) http.RoundTripper {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &roundTripper{
		user:  user,
		inner: inner,
	}
}

type roundTripper struct {
	user  *User
	inner http.RoundTripper
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.user.toHeaders(req.Header)
	return r.inner.RoundTrip(req)
}
