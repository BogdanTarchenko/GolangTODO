package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"runtime/debug"
	"todo/internal/domain/repository"
	"todo/internal/validation"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			log.Printf("Error: %v", err)
			switch {
			case errors.Is(err, repository.ErrTaskNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			case isValidationError(err):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				log.Printf("Stack trace: %s", debug.Stack())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			c.Abort()
		}
	}
}

func isValidationError(err error) bool {
	if errors.As(err, new(validator.ValidationErrors)) {
		return true
	}
	var vErr *validation.ValidationError
	if errors.As(err, &vErr) {
		return true
	}
	return false
}
