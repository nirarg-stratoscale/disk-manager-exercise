package httputil

import (
	"crypto/tls"
	"fmt"
	"net"
	net_http "net/http"
	"time"
)

// GetCustomHTTPClient return a custom HTTP client with specific source IP and un/secure SSL
func GetCustomHTTPClient(ip string, secure bool, keepAlive bool) (*net_http.Client, error) {
	if ip == "" {
		ip = "0.0.0.0"
	}
	a := net.ParseIP(ip)
	if a == nil {
		return nil, fmt.Errorf("Invalid IP %s", ip)
	}
	ka := 30 * time.Second
	if !keepAlive {
		ka = 0
	}
	transport := &net_http.Transport{
		Proxy: net_http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: ka,
			LocalAddr: &net.TCPAddr{IP: a},
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !secure}

	client := &net_http.Client{
		Transport: transport,
	}
	return client, nil
}
