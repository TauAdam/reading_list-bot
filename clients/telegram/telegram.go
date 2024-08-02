package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) Client {
	return Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}

}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	queryMap := url.Values{}
	queryMap.Add("offset", strconv.Itoa(offset))
	queryMap.Add("limit", strconv.Itoa(limit))

}

func (c *Client) SendMessage() {

}
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
}
