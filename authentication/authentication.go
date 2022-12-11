// Package authentication provides helpful authentication-related
// functions.
package authentication

import (
	"errors"
	"sync"
)

const (
	AuthHeaderKey    = "Authorization"
	AuthHeaderPrefix = "Bearer "
	OrgHeaderKey     = "OpenAI-Organization"
)

var (
	apiKey   = ""
	keyMutex = sync.RWMutex{}
)

// SetAPIKey sets the API user's API key
// for later retrieval with APIKey.
func SetAPIKey(key string) error {
	if len(key) == 0 {
		return errors.New("provided key was the empty string")
	}
	keyMutex.Lock()
	defer keyMutex.Unlock()
	apiKey = key
	return nil
}

// APIKey returns the API key set by
// calling SetAPIKey. It is required for
// this to be set.
func APIKey() string {
	keyMutex.RLock()
	defer keyMutex.RUnlock()
	return apiKey
}

var (
	defaultOrgID = ""
	orgMutex     = sync.RWMutex{}
)

func SetDefaultOrganizationID(orgID string) error {
	if len(orgID) == 0 {
		return errors.New("provided ID was the empty string")
	}
	orgMutex.Lock()
	defer orgMutex.Unlock()
	defaultOrgID = orgID
	return nil
}

// DefaultOrganizationID returns the default organization ID
// set by SetDefaultOrganizationID. It is not required for this
// to be set.
func DefaultOrganizationID() string {
	orgMutex.RLock()
	defer orgMutex.RUnlock()
	return defaultOrgID
}
