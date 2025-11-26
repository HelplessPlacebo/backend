package proxy

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/HelplessPlacebo/backend/pkg/shared"
)

type Client interface {
	ForwardJSON(ctx context.Context, targetPath string, body interface{}) (*http.Response, error)
}

type client struct {
	base    *url.URL
	httpCli *shared.HTTPClient
}

func NewClient(baseURL string) Client {
	u, _ := url.Parse(baseURL)
	return &client{
		base:    u,
		httpCli: shared.NewHTTPClient(shared.Duration("HTTP_CLIENT_TIMEOUT", 5)),
	}
}

func (c *client) ForwardJSON(ctx context.Context, targetPath string, body interface{}, incomingCookies []*http.Cookie) (*http.Response, error) {
	ref := &url.URL{Path: path.Join(c.base.Path, targetPath)}
	full := c.base.ResolveReference(ref).String()

	req, err := c.httpCli.NewJSONRequest(ctx, http.MethodPost, full, body)
	if err != nil {
		return nil, err
	}

	// 🔥 Прокидываем куки
	for _, ck := range incomingCookies {
		req.AddCookie(ck)
	}

	return c.httpCli.Do(req)
}
