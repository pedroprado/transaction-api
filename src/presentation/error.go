package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"pedroprado.transaction.api/src/core/domain/values"
	"reflect"
	"strings"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorHandler struct {
	err    error
	parser func(interface{}) interface{}
}

func NewErrorHandler(err error) *ErrorHandler {
	return &ErrorHandler{err: err}
}

func (ref *ErrorHandler) WithParser(parser func(interface{}) interface{}) *ErrorHandler {
	ref.parser = parser
	return ref
}

func getValidationErrors(err error) []string {
	var errorCause []string

	validationErrors, ok := err.(validator.ValidationErrors)

	if ok {
		for _, err := range validationErrors {
			switch err.Tag() {
			case "required":
				errorCause = append(errorCause, err.Field()+" is required")
			case "invalid":
				errorCause = append(errorCause, fmt.Sprintf("%v : %v is invalid", err.Field(), err.Param()))
			case "custom":
				errorCause = append(errorCause, fmt.Sprintf("%v", err.Param()))
			}
		}
	}

	if len(errorCause) == 0 {
		errorCause = []string{errors.Cause(err).Error()}
	}
	fmt.Println(errorCause)
	return errorCause
}

func SendBadRequestError(c *gin.Context, err error) {
	validationErrors := getValidationErrors(err)
	errorsString := strings.Join(validationErrors, ", ")
	errorResponse := ErrorResponse{
		Error: errorsString,
	}
	c.JSON(http.StatusBadRequest, errorResponse)
}

func SendNotFoundError(c *gin.Context, err error) {
	errorResponse := ErrorResponse{
		Error: errors.Cause(err).Error(),
	}
	c.JSON(http.StatusNotFound, errorResponse)
}

func (ref *ErrorHandler) Handle(c *gin.Context) {
	var (
		err    = ref.err
		parser = func(original interface{}) interface{} {
			return original
		}
	)

	if ref.parser != nil {
		parser = ref.parser
	}

	output := ErrorResponse{
		Error: errors.Cause(err).Error(),
	}

	logrus.WithField("exception", err).
		WithField("type", reflect.TypeOf(errors.Cause(err))).
		WithField("detail", errors.Cause(err).Error()).
		Error("Request failed")

	switch e := errors.Cause(err).(type) {
	case validator.ValidationErrors:
		SendBadRequestError(c, errors.Cause(err))
		return
	case *values.ErrorDuplicated:
		c.JSON(http.StatusConflict, parser(e.Original))
		return
	case *values.ErrorNotFound:
		c.JSON(http.StatusNotFound, output)
		return
	case *values.ErrorPreconditionNotFound:
		c.JSON(http.StatusPreconditionFailed, output)
		return
	case *values.ErrorPrecondition:
		c.JSON(http.StatusPreconditionFailed, output)
	case *values.ErrorValidation:
		err := ErrorResponse{
			Error: output.Error,
		}
		c.JSON(http.StatusBadRequest, err)
		return
	default:
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, output)
		return
	}
}
