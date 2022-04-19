package errutil

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	gopi "github.com/teejays/gopi"
)

var ErrNotFound = fmt.Errorf("not found")
var ErrNothingToUpdate = fmt.Errorf("nothing to update")
var ErrNotAuthorized = fmt.Errorf("you are not authorized to perform this action")

func Combine(errs ...error) error {
	return MultiErr{errs: removeNilErrs(errs)}
}

func removeNilErrs(errs []error) []error {
	var clean []error
	for _, err := range errs {
		if err != nil {
			clean = append(clean, err)
		}
	}
	return clean
}

func IsErrNoRows(err error) bool {
	return err == sql.ErrNoRows
}

var ErrBadCredentials = GokuError{
	internalError:      fmt.Errorf("bad credentials"),
	externalError:      fmt.Errorf("Provided credentials did not match those on the system"),
	externalHTTPStatus: http.StatusUnauthorized,
}

var ErrBadToken = GokuError{
	internalError:      fmt.Errorf("bad token"),
	externalError:      fmt.Errorf("Provided authentication token cannot be verified"),
	externalHTTPStatus: http.StatusUnauthorized,
}

type GokuError struct {
	internalError      error
	externalError      error // if left empty, the internal message will be used
	externalHTTPStatus int
}

func (err GokuError) Error() string {
	return err.internalError.Error()
}

// GetHTTPStatus returns the status, if set, and defaults to InternalServerError
func (err GokuError) GetExternalError() error {
	if err.externalError != nil {
		return err.externalError
	}
	return err.internalError
}

// GetHTTPStatus returns the status, if set, and defaults to InternalServerError
func (err GokuError) GetHTTPStatus() int {
	if err.externalHTTPStatus > 0 {
		return err.externalHTTPStatus
	}
	return http.StatusInternalServerError
}

func AsGokuError(err error) (GokuError, bool) {
	if errors.Is(err, ErrBadCredentials) {
		return ErrBadCredentials, true
	}
	return GokuError{}, false
}

// HandleHTTPResponseError handles the logic that: if the error is a GokuError, get the right external stuff otherwise defaults to teh default message
func HandleHTTPResponseError(w http.ResponseWriter, err error) {
	if gErr, ok := AsGokuError(err); ok {
		gopi.WriteError(w, gErr.GetHTTPStatus(), err, true, gErr.GetExternalError())
		return
	}
	gopi.WriteError(w, http.StatusBadRequest, err, false, nil)
	return
}
