package middleware

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net"
	"sync"
	"time"
)

type client struct {
	limiters map[string]*rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func init() {
	go cleanupClients()
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func (m middleware) IpRateLimiter(key string, limit rate.Limit, burst int, message string) gin.HandlerFunc {
	if message == "" {
		message = "Rate limit exceeded. Please try again later."
	}
	return func(ctx *gin.Context) {
		ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
		if err != nil {
			response.NewApiResponse(500, "fail to get user's host", err).Send(ctx)
			ctx.Abort()
			return
		}

		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiters: make(map[string]*rate.Limiter)}
		}
		if _, found := clients[ip].limiters[key]; !found {
			clients[ip].limiters[key] = rate.NewLimiter(limit, burst)
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiters[key].Allow() {
			mu.Unlock()
			response.NewApiResponse(429, message, gin.H{}).Send(ctx)
			ctx.Abort()
			return
		}
		mu.Unlock()
		ctx.Next()
	}
}
