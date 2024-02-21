package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var ipLimitMap = make(map[string][]time.Time)
var mapMutex = &sync.Mutex{}

func rateLimiter(c *gin.Context) {
	ip := c.ClientIP()
	now := time.Now()

	mapMutex.Lock()
	timestamps, exists := ipLimitMap[ip]
	if !exists {
		ipLimitMap[ip] = []time.Time{now}
		mapMutex.Unlock()
		c.Next()
		return
	}

	newTimestamps := []time.Time{now}
	for _, t := range timestamps {
		if now.Sub(t) < time.Minute {
			newTimestamps = append(newTimestamps, t)
		}
	}

	ipLimitMap[ip] = newTimestamps
	mapMutex.Unlock()

	if len(newTimestamps) > 3 {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	c.Next()
}
