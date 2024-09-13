package apiutil

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

func TryLock(c *gin.Context, locker *sync.Mutex, timeout time.Duration) bool {
	if !tryLock(locker, timeout) {
		c.JSON(http.StatusOK, Response{
			Code:    4000,
			Message: "Failed to acquire lock",
			Data:    gin.H{},
		})
		return false
	}
	return true
}

func tryLock(locker *sync.Mutex, timeout time.Duration) bool {
	c := make(chan struct{}, 1)
	go func() {
		locker.Lock()
		c <- struct{}{}
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}

func UseErrorHandler(r *gin.Engine) {
	r.Use(ErrorHandler())
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Code:    5000,
					Message: "Internal server error",
					Data:    EmptyResponse{},
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}
