package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
)

func testMiddlewareRequest(t *testing.T, r *gin.Engine, expectedHTTPCode int) {
	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == expectedHTTPCode
	})
}

// Sets auth info for use in tests
func setAuthInfoForTest(authInfo handlers.AuthorizerInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set((string)(handlers.CtxKeyAuthInfo), authInfo)
	}
}
