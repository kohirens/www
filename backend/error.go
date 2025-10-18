package backend

import "fmt"

type ProviderNotFound struct {
	name string
}

func (e *ProviderNotFound) Error() string {
	return fmt.Sprintf(stderr.ProviderNotFound, e.name)
}

// ReferralError Thrown when processing the request has failed in some way and
// there is somewhere else to send the client.
// NOTE: This is good for HTTP codes such as SeeOther, 302,301, etc. Consider
// this a treatable error, that you choose to log or not.
type ReferralError struct {
	Code        int
	ContentType string
	Location    string
	Log         bool
	msg         string
}

func (e *ReferralError) Error() string {
	return fmt.Sprintf(stderr.SeeOther, e.Location)
}

func NewReferralError(contentType, msg, loc string, code int, log bool) *ReferralError {
	return &ReferralError{
		Code:        code,
		ContentType: contentType,
		msg:         msg,
		Location:    loc,
		Log:         log,
	}
}

type ServiceNotFoundError struct {
	name string
}

func (e *ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service %v not found", e.name)
}

// UnauthorizedError Thrown when processing a request has failed authorization.
// If "Location" is set, you can opt to send the client there.
// NOTE: This is good for HTTP codes such as 401. Consider
// this a treatable error, that you choose to log or not.
type UnauthorizedError struct {
	Body        []byte
	Code        int
	ContentType string
	Location    string
	Log         bool
	msg         string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf(stderr.SeeOther, e.Location)
}

func NewUnauthorizedError(contentType, msg, loc string, code int, log bool, body []byte) *UnauthorizedError {
	return &UnauthorizedError{
		Body:        body,
		ContentType: contentType,
		Code:        code,
		Location:    loc,
		Log:         log,
		msg:         msg,
	}
}
