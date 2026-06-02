package handlers

import (
	"net"
	"net/http"
)

// LocalNetworkOnly blocks requests that do not come from localhost or a private LAN IP.
func LocalNetworkOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}

		ip := net.ParseIP(host)
		if ip == nil || (!ip.IsLoopback() && !ip.IsPrivate()) {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
