package storage

import (
	"errors"

	errtext "github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/domain/errors"
)

var ErrIvalidAuthData = errors.New(errtext.InvalidAuthDataError)
var ErrUserNotFound = errors.New(errtext.UserNotFoundError)
var ErrBookNotFound = errors.New(errtext.BookNotFoundError)
var ErrBookListEmrty = errors.New(errtext.BookListEmptyError)
