package exception

import (
	"net/http"

	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	// "gitlab.com/truemeds-dev-team/truemeds-dev-doctor/truemeds-dev-service/doctorportal-auth-service/shared/constants"
)

type ErrorCode string

const (
	USER_ALREADY_EXISTS                 = "USER_ALREADY_EXISTS"
	USER_NOT_FOUND                      = "USER_NOT_FOUND"
	USER_BLOCKED                        = "USER_BLOCKED"
	INTERNAL_SERVER_ERROR               = "INTERNAL_SERVER_ERROR"
	INVALID_CREDENTIALS                 = "Email ID and password do not match"
	INVALID_EMAIL                       = "INVALID_EMAIL"
	INVALID_PAYLOAD                     = "INVALID_PAYLOAD"
	EMAIL_REQUIRED                      = "EMAIL_REQUIRED"
	PASSWORD_REQUIRED                   = "PASSWORD_REQUIRED"
	WEAK_PASSWORD                       = "WEAK_PASSWORD"
	INVALID_OR_EXPIRED_TOKEN            = "INVALID_OR_EXPIRED_TOKEN"
	PASSWORD_CONFIRM_PASSWORD_NOT_MATCH = "PASSWORD_CONFIRM_PASSWORD_NOT_MATCH"
	TOKEN_NOT_FOUND                     = "TOKEN_NOT_FOUND"
	FIRST_NAME_REQUIRED                 = "FIRST_NAME_REQUIRED"
	LAST_NAME_REQUIRED                  = "LAST_NAME_REQUIRED "
	INVALID_USER_ID                     = "INVALID_USER_ID"
	USER_ID_REQUIRED                    = "USER_ID_REQUIRED"
	USER_ALREADY_BLOCKED                = "USER_ALREADY_BLOCKED"
	INVALID_API_KEY                     = "INVALID_API_KEY"
	ACCESS_TOKEN_REQUIRED               = "ACCESS_TOKEN_REQUIRED"
	INVALID_TOKEN_LOGOUT                = "INVALID_TOKEN_LOGOUT"
	INVALID_OR_EXPIRED_REFRESH_TOKEN    = "INVALID_OR_EXPIRED_REFRESH_TOKEN"
	REFRESH_TOKEN_REQUIRED              = "REFRESH_TOKEN_REQUIRED"
	INVALID_OR_EXPIRED_RESET_TOKEN      = "INVALID_OR_EXPIRED_RESET_TOKEN"
	INVALID_USERNAME                    = "INVALID_USERNAME"
	TM_EMAIL_REQUIRED                   = "TM_EMAIL_REQUIRED"
	INVALID_TM_EMAIL_REQUIRED           = "INVALID_TM_EMAIL_REQUIRED"
)

var ErrorCodeErrorMessage = map[ErrorCode]constants.Exception{
	USER_NOT_FOUND: {
		Code:           http.StatusNotFound,
		Message:        "Username does not exist",
		HttpStatusCode: http.StatusNotFound,
	},
	USER_ALREADY_EXISTS: {
		Code:           http.StatusConflict,
		Message:        "User Already exits",
		HttpStatusCode: http.StatusConflict,
	},
	USER_ALREADY_BLOCKED: {
		Code:           http.StatusConflict,
		Message:        "User already blocked",
		HttpStatusCode: http.StatusConflict,
	},
	USER_BLOCKED: {
		Code:           http.StatusForbidden,
		Message:        "Your account is blocked",
		HttpStatusCode: http.StatusForbidden,
	},
	INTERNAL_SERVER_ERROR: {
		Code:           http.StatusInternalServerError,
		Message:        "Something went wrong",
		HttpStatusCode: http.StatusInternalServerError,
	},
	INVALID_CREDENTIALS: {
		Code:           http.StatusUnauthorized,
		Message:        "Invalid Username Or Password",
		HttpStatusCode: http.StatusUnauthorized,
	},
	INVALID_EMAIL: {
		Code:           http.StatusBadRequest,
		Message:        "Invalid Email ID ",
		HttpStatusCode: http.StatusBadRequest,
	},
	INVALID_PAYLOAD: {
		Code:           http.StatusBadRequest,
		Message:        "Invalid payload",
		HttpStatusCode: http.StatusBadRequest,
	},
	INVALID_USER_ID: {
		Code:           http.StatusBadRequest,
		Message:        "Invalid User Id",
		HttpStatusCode: http.StatusBadRequest,
	},
	PASSWORD_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "Password is required",
		HttpStatusCode: http.StatusBadRequest,
	},
	FIRST_NAME_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "First name is required",
		HttpStatusCode: http.StatusBadRequest,
	},
	LAST_NAME_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "Last name is required",
		HttpStatusCode: http.StatusBadRequest,
	},
	EMAIL_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "Username is required",
		HttpStatusCode: http.StatusBadRequest,
	},
	USER_ID_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "User is required",
		HttpStatusCode: http.StatusBadRequest,
	},
	WEAK_PASSWORD: {
		Code:           http.StatusBadRequest,
		Message:        "Password must be 8 - 16 characters long with uppercase, lowercase, number, and special character",
		HttpStatusCode: http.StatusBadRequest,
	},
	INVALID_OR_EXPIRED_TOKEN: {
		Code:           constants.InvalidAccessToken,
		Message:        "Invalid or expired token",
		HttpStatusCode: http.StatusUnauthorized,
	},
	PASSWORD_CONFIRM_PASSWORD_NOT_MATCH: {
		Code:           http.StatusBadRequest,
		Message:        "Password and confirm password do not match",
		HttpStatusCode: http.StatusBadRequest,
	},
	TOKEN_NOT_FOUND: {
		Code:           http.StatusNotFound,
		Message:        "Token Not Found",
		HttpStatusCode: http.StatusNotFound,
	},
	INVALID_API_KEY: {
		Code:           http.StatusUnauthorized,
		Message:        "Invalid or missing api key",
		HttpStatusCode: http.StatusUnauthorized,
	},
	INVALID_TOKEN_LOGOUT: {
		Code:           constants.InvalidTokenLogout,
		Message:        "User logout",
		HttpStatusCode: http.StatusUnauthorized,
	},
	INVALID_OR_EXPIRED_REFRESH_TOKEN: {
		Code:           constants.InvalidTokenLogout,
		Message:        "Invalid or expired token",
		HttpStatusCode: http.StatusUnauthorized,
	},
	REFRESH_TOKEN_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "Refresh token is required",
		HttpStatusCode: http.StatusBadRequest,
	},
	INVALID_OR_EXPIRED_RESET_TOKEN: {
		Code:           constants.ResetTokenInvalidOrExpire,
		Message:        "Invalid or expired token",
		HttpStatusCode: http.StatusUnauthorized,
	},
	INVALID_USERNAME: {
		Code:           http.StatusBadRequest,
		Message:        "Invalid username ",
		HttpStatusCode: http.StatusBadRequest,
	},
	TM_EMAIL_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "truemeds email is required",
		HttpStatusCode: http.StatusBadRequest,
	},
	INVALID_TM_EMAIL_REQUIRED: {
		Code:           http.StatusBadRequest,
		Message:        "Invalid truemeds email",
		HttpStatusCode: http.StatusBadRequest,
	},
}

func GetException(code ErrorCode) *constants.Exception {
	if ex, ok := ErrorCodeErrorMessage[code]; ok {
		return &ex
	}
	return &constants.Exception{
		Code:           500,
		Message:        "Unknown error code",
		HttpStatusCode: 500,
	}
}
