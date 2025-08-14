package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type IPResponse struct {
	IP string `json:"ip"`
}

func getClientIP(c *gin.Context) string {
	headers := []string{
		"CF-Connecting-IP",
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Client-IP",
		"X-Forwarded",
		"X-Cluster-Client-IP",
		"Forwarded-For",
		"Forwarded",
	}

	for _, header := range headers {
		ip := c.GetHeader(header)
		if ip != "" {
			if header == "X-Forwarded-For" {
				ips := strings.Split(ip, ",")
				ip = strings.TrimSpace(ips[0])
			}
			if isValidIP(ip) {
				return ip
			}
		}
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

func getIPv4(c *gin.Context) string {
	ip := getClientIP(c)
	parsedIP := net.ParseIP(ip)
	if parsedIP != nil && parsedIP.To4() != nil {
		return ip
	}
	return ""
}

func getIPv6(c *gin.Context) string {
	ip := getClientIP(c)
	parsedIP := net.ParseIP(ip)
	if parsedIP != nil && parsedIP.To4() == nil {
		return ip
	}
	return ""
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func handleRoot(c *gin.Context) {
	ip := getClientIP(c)
	c.JSON(http.StatusOK, IPResponse{IP: ip})
}

func handleIPv4(c *gin.Context) {
	ip := getIPv4(c)
	if ip == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No IPv4 address found"})
		return
	}
	c.JSON(http.StatusOK, IPResponse{IP: ip})
}

func handleIPv6(c *gin.Context) {
	ip := getIPv6(c)
	if ip == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No IPv6 address found"})
		return
	}
	c.JSON(http.StatusOK, IPResponse{IP: ip})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"timestamp": time.Now().Unix(),
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", handleRoot)
	r.GET("/ipv4", handleIPv4)
	r.GET("/ipv6", handleIPv6)
	r.GET("/health", handleHealth)

	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}