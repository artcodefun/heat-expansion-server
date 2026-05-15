package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// bindRequest centralizes validation and error responses for Gin binding steps.
// It always returns false after writing a 400 response so callers can just return.
func bindRequest[U, Q, B any](c *gin.Context, req *dtos.Request[U, Q, B]) bool {
	if needsBind[U]() {
		if err := bindUri(c, &req.Uri); err != nil {
			writeBindError(c, "uri", err, &req.Uri)
			return false
		}
	}
	if needsBind[Q]() {
		if err := bindQuery(c, &req.Query); err != nil {
			writeBindError(c, "query", err, &req.Query)
			return false
		}
	}
	if needsBind[B]() {
		if err := bindJSON(c, &req.Body); err != nil {
			writeBindError(c, "json", err, &req.Body)
			return false
		}
	}
	return true
}

func writeBindError(c *gin.Context, source string, err error, dest interface{}) {
	message := bindErrorMessage(err, dest)
	slog.WarnContext(c.Request.Context(), "request rejected; invalid input", "method", c.Request.Method, "path", c.Request.URL.Path, "source", source, "error", message)
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func bindErrorMessage(err error, dest interface{}) string {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) && len(verrs) > 0 {
		fe := verrs[0]
		fieldName := friendlyFieldName(fe, dest)
		if param := fe.Param(); param != "" {
			return fmt.Sprintf("%s failed validation %s=%s", fieldName, fe.Tag(), param)
		}
		return fmt.Sprintf("%s failed validation %s", fieldName, fe.Tag())
	}
	return err.Error()
}

func friendlyFieldName(fe validator.FieldError, dest interface{}) string {
	if dest == nil {
		return fe.Field()
	}
	t := reflect.TypeOf(dest)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fe.Field()
	}
	fieldName := fe.StructField()
	if field, ok := t.FieldByName(fieldName); ok {
		for _, tagKey := range []string{"json", "form", "uri"} {
			if tag := field.Tag.Get(tagKey); tag != "" {
				name := strings.Split(tag, ",")[0]
				if name != "" && name != "-" {
					return name
				}
			}
		}
	}
	return fe.Field()
}

func needsBind[T any]() bool {
	var zero T
	return reflect.TypeOf(zero) != reflect.TypeOf(dtos.None{})
}

func bindUri(c *gin.Context, dest interface{}) error {
	return c.ShouldBindUri(dest)
}

func bindQuery(c *gin.Context, dest interface{}) error {
	return c.ShouldBindQuery(dest)
}

func bindJSON(c *gin.Context, dest interface{}) error {
	return c.ShouldBindJSON(dest)
}
