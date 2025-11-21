package errors

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid username or password")
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

	// Student errors
	ErrStudentNotFound      = errors.New("student not found")
	ErrStudentAlreadyExists = errors.New("student already exists")

	// Fan errors
	ErrFanNotFound = errors.New("fan not found")

	// Teacher errors
	ErrTeacherNotFound      = errors.New("teacher not found")
	ErrTeacherAlreadyExists = errors.New("teacher already exists")
	ErrTeacherUserAssigned  = errors.New("user already linked to another teacher")

	// Schedule slot errors
	ErrScheduleSlotNotFound = errors.New("schedule slot not found")
	ErrScheduleSlotConflict = errors.New("schedule slot conflict")
	ErrScheduleSlotInactive = errors.New("schedule slot inactive")

	// Subject errors
	ErrSubjectNotFound = errors.New("subject not found")

	// Class schedule errors
	ErrClassScheduleNotFound = errors.New("class schedule not found")
	ErrClassScheduleConflict = errors.New("class schedule conflict")

	// SKS errors
	ErrSKSDefinitionNotFound     = errors.New("sks definition not found")
	ErrSKSDefinitionAlreadyExist = errors.New("sks definition already exists")
	ErrSKSExamScheduleNotFound   = errors.New("sks exam schedule not found")
	ErrStudentSKSResultNotFound  = errors.New("student sks result not found")

	// Attendance errors
	ErrAttendanceSessionNotFound = errors.New("attendance session not found")
	ErrAttendanceAlreadyLocked   = errors.New("attendance session already locked")
	ErrAttendanceInvalidStatus   = errors.New("invalid attendance status")

	// Leave/health errors
	ErrLeavePermitNotFound   = errors.New("leave permit not found")
	ErrLeavePermitConflict   = errors.New("leave permit conflict")
	ErrLeavePermitStatus     = errors.New("invalid leave permit status transition")
	ErrHealthStatusNotFound  = errors.New("health status not found")
	ErrHealthStatusActive    = errors.New("health status already active")
	ErrHealthStatusForbidden = errors.New("operation not allowed for current health status")

	// Class errors
	ErrClassNotFound          = errors.New("class not found")
	ErrStudentAlreadyEnrolled = errors.New("student already enrolled in class")
	ErrClassStaffExists       = errors.New("staff already assigned to class")

	// General errors
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
)
