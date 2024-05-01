package constants

type key int

const (
	SessionIsAuthorized key = iota
	SessionSkipAuthorization
	SessionID
	SessionIPAddress
	SessionUser
	SessionUserRole
	SessionUserID
	SessionUserUUID
	SessionUserTimezone
	SessionUserName
	SessionUserLexicalName
	SessionUserFirstName
	SessionUserLastName
	SessionUserOrganizationID
	SessionUserOrganizationName
	SessionUserOTPValidated
	SessionRequestID
)
