package v1

import "github.com/gin-gonic/gin"

type response struct {
	Message string `json:"message"`
}

// newResponse sends a JSON response with the given status code and message.
//
// This function aborts the current HTTP request and writes a JSON response
// using the provided statusCode and message. It sets the response payload to
// include a message field with the specified message.
//
// Parameters:
//   - c: The Gin context for the current HTTP request.
//   - statusCode: The HTTP status code to set in the response.
//   - message: The message to include in the response payload.
func newResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, response{message})
}
