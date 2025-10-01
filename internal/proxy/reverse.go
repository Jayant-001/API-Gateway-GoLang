package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)


func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	
	targetUrl, err := url.Parse(target);
	if err != nil {
		return nil, err
	}

	proxyHost := httputil.NewSingleHostReverseProxy(targetUrl)
	proxyHost.Director = func(req *http.Request) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = targetUrl.Host
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
	}
	return proxyHost, nil
}