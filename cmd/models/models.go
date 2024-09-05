package models

import "errors"

var (
	ErrUserBrokeConn      error = errors.New("user broked Connection")
	ErrNameExistsInServer error = errors.New("name exists in this Chat")
	ErrNameIsEmpety       error = errors.New("name is empety")
	ErrNameHasIllegalSims error = errors.New("name has illegal sims")
	ErrServerIsFull       error = errors.New("server is full")
)
