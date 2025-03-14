package realdebrid

import (
	"net/http"
)

type RealDebridType string

const (
	RealDebridTypeFree    RealDebridType = "free"
	RealDebridTypePremium RealDebridType = "premium"
)

type RealDebridUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Fidelity points
	Points int `json:"points"`
	// User language
	Locale string `json:"locale"`
	// URL
	Avatar string `json:"avatar"`
	// "premium" or "free"
	Type RealDebridType `json:"type"`
	// seconds left as a Premium user
	Premium int `json:"premium"`
	// jsonDate
	Expiration string `json:"expiration"`
}

// Returns some informations on the current user.
func (client *RealDebridClient) GetUser() (*RealDebridUser, error) {
	req, err := client.newRequest(http.MethodGet, "/user", nil, "", nil)

	if err != nil {
		return nil, err
	}

	var user RealDebridUser

	err = client.do(req, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
