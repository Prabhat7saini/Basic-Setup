package constants

type Exception struct {
	Code           int    `json:"code"`
	Message        string `json:"message"`
	HttpStatusCode int    `json:"httpStatusCode"`
}

// IServiceOutput<T> equivalent
	type ServiceOutput[T any] struct {
		Message        string     `json:"message,omitempty"`
		OutputData     T          `json:"outputData,omitempty"`
		Exception      *Exception `json:"exception,omitempty"`
		HttpStatusCode int        `json:"httpStatusCode"`
		RespStatusCode int        `json:"respStatusCode"`
	}

// Final API response structure
type ApiResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

const (
	Access_Token  string = "access_token"
	Refresh_Token string = "refresh_token"
)

type HttpServiceResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// RedisKey
const (
	ForgotPasswordRedisKey    = "auth:forgotPassword:%s"
	LoginAccessTokenRedisKey  = "auth:user:accessToken:%s"
	LoginRefreshTokenRedisKey = "auth:user:refreshToken:%s"
	VerifyAccessTokenRedisKey = "auth:user:ids:%s"
)

// custom error codes
const (
	InvalidAccessToken        int = 601 // take refreshToken
	InvalidTokenLogout        int = 602 //logout
	ResetTokenInvalidOrExpire int = 603 //resetToken
)

// API endpoints
const (
	SendEmailUrl = "/email/v1/send"
)

const TruemedsEmailExtension = "@truemeds.in"

// type tracing string

const (
	XTracingID          = "x-tracing-id"
	TracingIDKey    string     = "tracing_id"
	RequestStartTimeKey string = "request_start_time"
)
