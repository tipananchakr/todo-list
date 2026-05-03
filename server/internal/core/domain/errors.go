package domain

import "errors"

var ErrInvalidTodoID = errors.New("invalid todo id")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidToken = errors.New("invalid token")
var ErrUserNotFound = errors.New("user not found")
