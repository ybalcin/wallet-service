package errr

import (
	"errors"
	"net/http"
)

const AnErrorOccurredText = "an error occurred"

// ThrowBadRequestError throws Error with code http.StatusBadRequest
func ThrowBadRequestError(e error) *Error {
	return New(e, http.StatusBadRequest)
}

// ThrowForbiddenError throws Error with code http.StatusForbidden
func ThrowForbiddenError(e error) *Error {
	return New(e, http.StatusForbidden)
}

// ThrowInternalServerError throws Error with code http.StatusInternalServerError
func ThrowInternalServerError(e error) *Error {
	if e == nil {
		e = errors.New(AnErrorOccurredText)
	}

	return New(e, http.StatusInternalServerError)
}

// ThrowNotFoundError throws Error with http.StatusNotFound
func ThrowNotFoundError(e error) *Error {
	return New(e, http.StatusNotFound)
}

// ThrowUnauthorizedError throw Error with http.StatusUnauthorized
func ThrowUnauthorizedError(e error) *Error {
	return New(e, http.StatusUnauthorized)
}
