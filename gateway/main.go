package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Read optional config from env
	userSvc := getEnv("USER_SERVICE_URL", "http://user-service:8080")
	bookSvc := getEnv("BOOK_SERVICE_URL", "http://book-service:8081")
	addr := getEnv("GATEWAY_ADDR", ":8000")
	enableAuth := getEnv("GATEWAY_AUTH_ENABLED", "true") == "true"

	// Gin router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(BasicLoggerMiddleware())

	// Health & metrics
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if userSvc == "" || bookSvc == "" {
		log.Fatalf("service URL not set: user=%q book=%q", userSvc, bookSvc)
	}

	r.POST("/register", ProxyHandler(userSvc))
	r.POST("/login", ProxyHandler(userSvc))

	// JWT middleware (applies to proxied endpoints)
	var jwtMiddleware gin.HandlerFunc
	if enableAuth {
		jwtMiddleware = JWTMiddleware() // validates token or calls introspect based on env
	} else {
		jwtMiddleware = func(c *gin.Context) { c.Next() }
	}

	// Proxy routes
	r.Any("/users/*path", jwtMiddleware, ProxyHandler(userSvc))
	r.Any("/books/*path", jwtMiddleware, ProxyHandler(bookSvc))

	log.Printf("Gateway listening on %s (proxying users -> %s, books -> %s). Auth enabled=%v", addr, userSvc, bookSvc, enableAuth)
	if err := r.Run(addr); err != nil {
		log.Fatalf("gateway failed: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
