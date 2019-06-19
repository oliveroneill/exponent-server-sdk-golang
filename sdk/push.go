package expo

import (
	"errors"
	"net/http"
	"strings"
)

// ExponentPushToken is a valid Expo push token
type ExponentPushToken string

// ErrMalformedToken is returned if a token does not start with 'ExponentPushToken'
var ErrMalformedToken = errors.New("Token should start with ExponentPushToken")

//NewExponentPushToken returns a token and may return an error if the input token is invalid
func NewExponentPushToken(token string) (ExponentPushToken, error) {
	if !strings.HasPrefix(token, "ExponentPushToken") {
		return "", ErrMalformedToken
	}
	return ExponentPushToken(token), nil
}

const (
	// DefaultPriority is the standard priority used in PushMessage
	DefaultPriority = "default"
	// NormalPriority is a priority used in PushMessage
	NormalPriority = "normal"
	// HighPriority is a priority used in PushMessage
	HighPriority = "high"
)

// PushMessage is an object that describes a push notification request.
// Fields:
//	To: an ExponentPushToken
//	Data: A dict of extra data to pass inside of the push notification.
//	      The total notification payload must be at most 4096 bytes.
//	Title: The title to display in the notification. On iOS, this is
//	       displayed only on Apple Watch.
//	Body: The message to display in the notification.
//	Sound: A sound to play when the recipient receives this
//	       notification. Specify "default" to play the device's default
//	       notification sound, or omit this field to play no sound.
//	TTLSeconds: The number of seconds for which the message may be kept around
//	     for redelivery if it hasn't been delivered yet. Defaults to 0.
//	Expiration: UNIX timestamp for when this message expires. It has
//	        the same effect as ttl, and is just an absolute timestamp
//	        instead of a relative one.
//	Priority: Delivery priority of the message. Use the *Priority constants
//          specified above.
//	Badge: An integer representing the unread notification count. This
//	       currently only affects iOS. Specify 0 to clear the badge count.
//	ChannelID: ID of the Notification Channel through which to display this
//         notification on Android devices.
type PushMessage struct {
	To         ExponentPushToken `json:"to"`
	Body       string            `json:"body"`
	Data       map[string]string `json:"data,omitempty"`
	Sound      string            `json:"sound,omitempty"`
	Title      string            `json:"title,omitempty"`
	TTLSeconds int               `json:"ttl,omitempty"`
	Expiration int64             `json:"expiration,omitempty"`
	Priority   string            `json:"priority,omitempty"`
	Badge      int               `json:"badge,omitempty"`
	ChannelID  string            `json:"channelId,omitempty"`
}

// Response is the HTTP response returned from an Expo publish HTTP request
type Response struct {
	Data   []PushResponse      `json:"data"`
	Errors []map[string]string `json:"errors"`
}

// SuccessStatus is the status returned from Expo on a success
const SuccessStatus = "ok"

// ErrorDeviceNotRegistered indicates the token is invalid
const ErrorDeviceNotRegistered = "DeviceNotRegistered"

// ErrorMessageTooBig indicates the message went over payload size of 4096 bytes
const ErrorMessageTooBig = "MessageTooBig"

// ErrorMessageRateExceeded indicates messages have been sent too frequently
const ErrorMessageRateExceeded = "MessageRateExceeded"

// PushResponse is a wrapper class for a push notification response.
// A successful single push notification:
//     {'status': 'ok'}
// An invalid push token
//     {'status': 'error',
//      'message': '"adsf" is not a registered push notification recipient'}
type PushResponse struct {
	PushMessage PushMessage
	Status      string            `json:"status"`
	Message     string            `json:"message"`
	Details     map[string]string `json:"details"`
}

func (r *PushResponse) isSuccess() bool {
	return r.Status == SuccessStatus
}

// ValidateResponse returns an error if the response indicates that one occurred.
// Clients should handle these errors, since these require custom handling
// to properly resolve.
func (r *PushResponse) ValidateResponse() error {
	if r.isSuccess() {
		return nil
	}
	err := &PushResponseError{
		Response: r,
	}
	// Handle specific errors if we have information
	if r.Details != nil {
		e := r.Details["error"]
		if e == ErrorDeviceNotRegistered {
			return &DeviceNotRegisteredError{
				PushResponseError: *err,
			}
		} else if e == ErrorMessageTooBig {
			return &MessageTooBigError{
				PushResponseError: *err,
			}
		} else if e == ErrorMessageRateExceeded {
			return &MessageRateExceededError{
				PushResponseError: *err,
			}
		}
	}
	return err
}

// PushResponseError is a base class for all push reponse errors
type PushResponseError struct {
	Response *PushResponse
}

func (e *PushResponseError) Error() string {
	if e.Response != nil {
		return e.Response.Message
	}
	return "Unknown push response error"
}

// DeviceNotRegisteredError is raised when the push token is invalid
// To handle this error, you should stop sending messages to this token.
type DeviceNotRegisteredError struct {
	PushResponseError
}

// MessageTooBigError is raised when the notification was too large.
// On Android and iOS, the total payload must be at most 4096 bytes.
type MessageTooBigError struct {
	PushResponseError
}

// MessageRateExceededError is raised when you are sending messages too frequently to a device
// You should implement exponential backoff and slowly retry sending messages.
type MessageRateExceededError struct {
	PushResponseError
}

// PushServerError is raised when the push token server is not behaving as expected
// For example, invalid push notification arguments result in a different
// style of error. Instead of a "data" array containing errors per
// notification, an "error" array is returned.
// {"errors": [
//   {"code": "API_ERROR",
//    "message": "child \"to\" fails because [\"to\" must be a string]. \"value\" must be an array."
//   }
// ]}
type PushServerError struct {
	Message      string
	Response     *http.Response
	ResponseData *Response
	Errors       []map[string]string
}

// NewPushServerError creates a new PushServerError object
func NewPushServerError(message string, response *http.Response,
	responseData *Response,
	errors []map[string]string) *PushServerError {
	return &PushServerError{
		Message:      message,
		Response:     response,
		ResponseData: responseData,
		Errors:       errors,
	}
}

func (e *PushServerError) Error() string {
	return e.Message
}
