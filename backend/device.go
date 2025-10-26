package backend

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
)

type Device struct {
	ID           string              `json:"id"`
	OIDCProvider string              `json:"oidc_provider"`
	SessionID    string              `json:"session_id"`
	UserAgent    useragent.UserAgent `json:"user_agent"`
}

func NewDeviceId(b []byte) string {
	id := uuid.NewSHA1(uuid.NameSpaceOID, b)
	return id.String()
}

func NewDevice(uaMeta []byte, sessionID, oidcProvider string) (*Device, error) {
	var ua useragent.UserAgent
	if e := json.Unmarshal(uaMeta, &ua); e != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, e.Error())
	}

	return &Device{
		ID:           NewDeviceId(uaMeta),
		OIDCProvider: oidcProvider,
		SessionID:    sessionID,
		UserAgent:    ua,
	}, nil
}
