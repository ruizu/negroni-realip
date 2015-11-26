package realip

import (
	"strings"
	"net/http"
)

type Middleware struct {
	header         string
	trustedProxies []string
}

func New(ip []string) *Middleware {
	return &Middleware{
		header: "X-Forwarded-For",
		trustedProxies: ip,
	}
}

func NewCustomHeader(h string, ip []string) *Middleware {
	return &Middleware{
		header: h,
		trustedProxies: ip,
	}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h := r.Header.Get("X-Forwarded-For")
	if h == "" {
		next(w, r)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	if checkIP(m.trustedProxies, ip) {
		r.RemoteAddr = strings.TrimSpace(strings.Split(h, ",")[0])
	}

	next(w, r)
}

func checkIP(p []string, ip string) bool {
	for _, v := range p {
		if v == ip {
			return true
		}
	}
	return false
}
