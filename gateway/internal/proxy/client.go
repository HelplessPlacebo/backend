package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"time"

	"github.com/HelplessPlacebo/backend/pkg/shared"
)

type Client interface {
	ForwardJSON(ctx context.Context, req *http.Request, targetPath string, body interface{}) (*http.Response, error)
}

type client struct {
	base    *url.URL
	httpCli *http.Client
}

func NewClient(baseURL string) Client {
	u, _ := url.Parse(baseURL)

	jar, _ := cookiejar.New(nil)

	return &client{
		base: u,
		httpCli: &http.Client{
			Timeout: time.Second * 10,
			Jar:     jar,
		},
	}
}

func (c *client) ForwardJSON(ctx context.Context, originalReq *http.Request, targetPath string, body interface{}) (*http.Response, error) {
	ref := &url.URL{Path: path.Join(c.base.Path, targetPath)}
	fullURL := c.base.ResolveReference(ref)

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	upReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL.String(), &buf)
	if err != nil {
		return nil, err
	}

	upReq.Header.Set("Content-Type", "application/json")

	shared.CopyIncomingHeaders(originalReq, upReq)
	shared.CopyIncomingCookies(originalReq, upReq)

	return c.httpCli.Do(upReq)
}
