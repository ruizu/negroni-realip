package realip

import (
	"net"
	"net/http"
	"strings"
)

type Middleware struct {
	header         string
	trustedProxies []*net.IPNet
}

func New(ip []string) *Middleware {
	return NewCustomHeader("X-Forwarded-For", ip)
}

func NewCustomHeader(h string, ip []string) *Middleware {
	cidr := make([]*net.IPNet, len(ip))
	for i, v := range ip {
		if !strings.Contains(v, "/") {
			v += "/32"
		}
		_, c, err := net.ParseCIDR(v)
		if err != nil {
			return nil
		}
		cidr[i] = c
	}
	return &Middleware{header: h, trustedProxies: cidr}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h := r.Header.Get("X-Forwarded-For")
	if h == "" {
		next(w, r)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	if m.checkIP(ip) {
		r.RemoteAddr = strings.TrimSpace(strings.Split(h, ",")[0])
	}

	next(w, r)
}

func (m *Middleware) checkIP(ip string) bool {
	ip0 := net.ParseIP(ip)
	for _, v := range m.trustedProxies {
		if v.Contains(ip0) {
			return true
		}
	}
	return false
}
