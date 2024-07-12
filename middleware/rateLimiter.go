package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/JGLTechnologies/gin-rate-limit"
)

func RateLimiter() gin.HandlerFunc {
    store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
        Rate:  time.Second,
        Limit: 5, // the amount of requests that can be made every Rate
    })

    mw := ratelimit.RateLimiter(store, &ratelimit.Options{
        ErrorHandler: errorHandler,
        KeyFunc:      keyFunc,
    })

    return mw
}

func keyFunc(c *gin.Context) string {
    return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
    c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}
