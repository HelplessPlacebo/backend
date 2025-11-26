package shared

import "net/http"

var hopByHopHeaders = map[string]struct{}{
	"Connection":          {},
	"Proxy-Connection":    {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"Te":                  {},
	"Trailers":            {},
	"Transfer-Encoding":   {},
	"Upgrade":             {},
}

func CopyIncomingHeaders(in *http.Request, out *http.Request) {
	for k, vv := range in.Header {
		if _, skip := hopByHopHeaders[k]; skip {
			continue
		}
		for _, v := range vv {
			out.Header.Add(k, v)
		}
	}
}

func CopyIncomingCookies(in *http.Request, out *http.Request) {
	for _, c := range in.Cookies() {
		out.AddCookie(c)
	}
}

func ProxyCookiesFromResponse(in *http.Response, w http.ResponseWriter) {
	cookies := in.Header.Values("Set-Cookie")
	for _, c := range cookies {
		w.Header().Add("Set-Cookie", c)
	}
}
