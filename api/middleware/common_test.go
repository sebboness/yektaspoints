package middleware

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func testHTTPResponse(r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder)) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	// Run assertion function
	f(w)
}
