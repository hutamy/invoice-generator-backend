package errors

import e "errors"

var (
	ErrUserAlreadyExists = e.New("email already exists")
	ErrLoginFailed       = e.New("invalid email and password")
	ErrBadRequest        = e.New("please check your input")
)
