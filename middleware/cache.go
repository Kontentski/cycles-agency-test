package middleware

import (
    "log"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/patrickmn/go-cache"
)

var c *cache.Cache

// InitializeCache initializes the cache
func InitializeCache(defaultExpiration, cleanupInterval time.Duration) {
    c = cache.New(defaultExpiration, cleanupInterval)
}

// CacheResponse middleware to cache responses
func CacheResponse() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // Check if the cache is initialized
        if c == nil {
            log.Println("Cache is not initialized")
            ctx.Next()
            return
        }

        // Check if the response is cached
        if cachedData, found := c.Get(ctx.Request.RequestURI); found {
            log.Printf("Cache hit for %s", ctx.Request.RequestURI)
            ctx.Data(200, "application/json", cachedData.([]byte))
            ctx.Abort()
            return
        }

        log.Printf("Cache miss for %s", ctx.Request.RequestURI)
        // Capture the response
        writer := &responseWriter{ctx.Writer, &[]byte{}}
        ctx.Writer = writer
        ctx.Next()

        // Store the response in the cache if the status code is 200
        if ctx.Writer.Status() == 200 {
            c.Set(ctx.Request.RequestURI, *writer.body, cache.DefaultExpiration)
        }
    }
}

type responseWriter struct {
    gin.ResponseWriter
    body *[]byte
}

func (w *responseWriter) Write(b []byte) (int, error) {
    *w.body = append(*w.body, b...)
    return w.ResponseWriter.Write(b)
}
