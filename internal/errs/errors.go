package errs

import "errors"

var (
    // ===== User =====
    UserNameRequired = errors.New("userName is required")
    UserNotFound = errors.New("user not found")
)