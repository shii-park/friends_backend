package errs

import "errors"

var (
    // ===== User =====
    UserNameRequired = errors.New("ユーザーネームは必須です")
    UserNotFound = errors.New("ユーザーが存在しません")
)