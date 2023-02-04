package values

import "fmt"

type ErrorDuplicated struct {
	error
	Original interface{}
}

// To use when a resource is duplicated
func NewErrorDuplicated(original interface{}) error {
	return &ErrorDuplicated{Original: original}
}

func (ref *ErrorDuplicated) Error() string {
	return "Duplicated"
}

type ErrorNotFound struct {
	error
	s string
}

// To use when a entity is not found
func NewErrorNotFound(entityName string) error {
	return &ErrorNotFound{s: entityName}
}

func (ref *ErrorNotFound) Error() string {
	return fmt.Sprintf("%s not found", ref.s)
}

type ErrorPreconditionNotFound struct {
	error
	s string
}

// To use when a precondition entity is not found
func NewErrorPreconditionNotFound(entityName string) error {
	return &ErrorPreconditionNotFound{s: entityName}
}

func (ref *ErrorPreconditionNotFound) Error() string {
	return fmt.Sprintf("%s not found", ref.s)
}

type ErrorPrecondition struct {
	error
	message string
}

func NewErrorPrecondition(message string) error {
	return &ErrorPrecondition{message: message}
}

func (ref *ErrorPrecondition) Error() string {
	return ref.message
}

type ErrorValidation struct {
	error
	s string
}

// To use when something in the request is not valid
func NewErrorValidation(message string) error {
	return &ErrorValidation{s: message}
}

func (ref *ErrorValidation) Error() string {
	return ref.s
}

var (
	ErrorDestinationAccountNotFound  = NewErrorValidation("destination account not found")
	ErrorOriginAccountNotFound       = NewErrorValidation("origin account not found")
	ErrorIntermediaryAccountNotFound = NewErrorValidation("intermediary account not found")
)
