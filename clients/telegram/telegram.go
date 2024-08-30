package telegram

import (
	"encoding/json"
	"github.com/tauadam/reading_list-bot/lib/utils"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const TGUpdateMethod = "getUpdates"
const TGSendMessageMethod = "sendMessage"

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}

}

// Updates represents a telegram update
func (c *Client) Updates(offset, limit int) ([]Update, error) {
	queryMap := url.Values{}
	queryMap.Add("offset", strconv.Itoa(offset))
	queryMap.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(TGUpdateMethod, queryMap)
	if err != nil {
		return nil, utils.Wrap("could not get updates", err)
	}

	var resp UpdatesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, utils.Wrap("could not unmarshal response", err)
	}

	return resp.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	queryMap := url.Values{}
	queryMap.Add("chat_id", strconv.Itoa(chatID))
	queryMap.Add("text", text)

	_, err := c.doRequest(TGSendMessageMethod, queryMap)
	if err != nil {
		return utils.Wrap("could not send message", err)
	}
	return nil
}

// doRequest working under the hood of the Client to communicate with the telegram API
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, utils.Wrap("could not create request", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, utils.Wrap("could not send request", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.Wrap("could not read response body", err)
	}

	return body, nil
}
