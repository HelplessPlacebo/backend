package authclient

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/HelplessPlacebo/backend/pkg/shared"
)

type SessionResponse struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

type AuthClient interface {
	RequestSession(ctx context.Context, r *http.Request) (bool, *SessionResponse, int, error)

	RequestRefresh(ctx context.Context, r *http.Request) (bool, *SessionResponse, []*http.Cookie, int, error)
}

type client struct {
	sessionEndpoint string
	refreshEndpoint string
	client          *http.Client
}

func New(sessionEndpoint string, refreshEndpoint string, timeout time.Duration) AuthClient {
	return &client{
		sessionEndpoint: sessionEndpoint,
		refreshEndpoint: refreshEndpoint,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *client) RequestSession(ctx context.Context, r *http.Request) (bool, *SessionResponse, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.sessionEndpoint, nil)

	if err != nil {
		return false, nil, 0, err
	}

	shared.CopyIncomingCookies(r, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, nil, 0, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode

	if status != http.StatusOK {
		return false, nil, status, nil
	}

	var s SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return false, nil, status, err
	}
	return true, &s, status, nil
}

func (c *client) RequestRefresh(ctx context.Context, r *http.Request) (bool, *SessionResponse, []*http.Cookie, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.refreshEndpoint, nil)

	if err != nil {
		return false, nil, nil, 0, err
	}

	shared.CopyIncomingCookies(r, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, nil, nil, 0, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status != http.StatusOK {
		return false, nil, nil, status, nil
	}

	var s SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return false, nil, nil, status, err
	}

	setCookies := resp.Cookies()
	return true, &s, setCookies, status, nil
}
