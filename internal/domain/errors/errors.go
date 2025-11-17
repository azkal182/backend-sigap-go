package errors

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenNotFound      = errors.New("token not found")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserInactive      = errors.New("user is inactive")

	// Role errors
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrProtectedRole     = errors.New("cannot modify protected role")

	// Permission errors
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission already exists")
	ErrPermissionDenied        = errors.New("permission denied")

	// Dormitory errors
	ErrDormitoryNotFound      = errors.New("dormitory not found")
	ErrDormitoryAlreadyExists = errors.New("dormitory already exists")
	ErrDormitoryAccessDenied  = errors.New("access denied to this dormitory")

	// General errors
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
)
