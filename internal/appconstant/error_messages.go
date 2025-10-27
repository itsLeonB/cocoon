package appconstant

const (
	ErrDataSelect = "error retrieving data"
	ErrDataInsert = "error inserting new data"

	ErrAuthUserNotFound       = "user is not found"
	ErrAuthDuplicateUser      = "user with email %s is already registered"
	ErrAuthUnknownCredentials = "unknown credentials, please check your email/password"

	ErrUserNotFound = "user is not found"
	ErrUserDeleted  = "user with ID: %s is deleted"

	ErrFriendshipNotFound = "friendship not found"
	ErrFriendshipDeleted  = "friendship is deleted"

	ErrStructValidation = "error validating struct"

	ErrNilRequest = "request is nil"
)
