package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"

	userCtx = "id"
)

// userIdentity is a middleware that extracts the user ID from the Authorization header
// and stores it in the request context.
//
// The middleware expects the Authorization header to be in the format "Bearer <token>".
// If the header is empty or invalid, or if the token is invalid, the middleware returns
// a 401 error with a corresponding error message.
//
// The user ID is stored in the request context under the key "userId".
func (h *Handler) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}

// parseAuthHeader extracts and validates the JWT token from the Authorization header.
//
// This function retrieves the Authorization header from the provided Gin context,
// verifies that it is in the format "Bearer <token>", and returns the token if valid.
// If the header is missing, improperly formatted, or the token is empty, an error is returned.
//
// Parameters:
//   - c: The Gin context for the current HTTP request.
//
// Returns:
//   - string: The extracted token if the header is valid.
//   - error: An error if the header is empty, invalid, or the token cannot be retrieved.
func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

// getUserId retrieves the user ID from the Gin context.
//
// The function retrieves the value associated with the "userId" context key,
// verifies that it is a string, and attempts to parse it as a UUID.
// If the value is not found, is of an invalid type, or cannot be parsed as a UUID,
// an error is returned.
//
// Parameters:
//   - c: The Gin context for the current HTTP request.
//
// Returns:
//   - uuid.UUID: The parsed UUID if the value is found and valid.
//   - error: An error if the value is not found, is of an invalid type, or cannot be parsed as a UUID.
func getUserId(c *gin.Context) (uuid.UUID, error) {
	return getIdByContext(c, userCtx)
}

// getIdByContext retrieves the UUID value from the provided Gin context.
//
// The function retrieves the value associated with the provided context key,
// verifies that it is a string, and attempts to parse it as a UUID.
// If the value is not found, is of an invalid type, or cannot be parsed as a UUID,
// an error is returned.
//
// Parameters:
//   - c: The Gin context for the current HTTP request.
//   - context: The key of the value to be retrieved from the context.
//
// Returns:
//   - uuid.UUID: The parsed UUID if the value is found and valid.
//   - error: An error if the value is not found, is of an invalid type, or cannot be parsed as a UUID.
func getIdByContext(c *gin.Context, context string) (uuid.UUID, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return uuid.Nil, errors.New("id not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return uuid.Nil, errors.New("id is of invalid type")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
