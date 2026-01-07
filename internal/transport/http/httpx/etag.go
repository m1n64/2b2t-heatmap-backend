package httpx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ETagHandler struct {
	OnHit  func()
	OnMiss func()
}

func HandleETag(c *gin.Context, etag string, h ETagHandler) bool {
	if c.GetHeader("If-None-Match") == etag {
		if h.OnHit != nil {
			h.OnHit()
		}
		c.Status(http.StatusNotModified)
		return true
	}

	c.Header("ETag", etag)
	if h.OnMiss != nil {
		h.OnMiss()
	}
	return false
}
