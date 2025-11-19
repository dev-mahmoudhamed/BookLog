package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(target string) gin.HandlerFunc {
	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.Host = u.Host
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/books")
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		// optionally modify responses here
		return nil
	}

	return func(c *gin.Context) {
		// Pass request through to proxy
		// Copy Authorization explicitly (should be present)
		if auth := c.GetHeader("Authorization"); auth != "" {
			c.Request.Header.Set("Authorization", auth)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
