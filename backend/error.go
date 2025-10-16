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
// NOTE: This is good for HTTP codes such ass SeeOther, 302,301, etc. Consider
// this a treatable error, that you choose to log or not.
type ReferralError struct {
	msg      string
	Location string
	Code     int
	Log      bool
}

func (e *ReferralError) Error() string {
	return fmt.Sprintf(stderr.SeeOther, e.Location)
}

func NewReferralError(msg, loc string, code int, log bool) *ReferralError {
	return &ReferralError{
		msg:      msg,
		Location: loc,
		Code:     code,
		Log:      log,
	}
}

type ServiceNotFoundError struct {
	name string
}

func (e *ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service %v not found", e.name)
}
