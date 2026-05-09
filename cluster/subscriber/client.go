package subscriber

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gohttp"
)

type Client struct {
	client  gohttp.Client
	baseURL string
}

func NewClient(client gohttp.Client, baseURL string) *Client {
	return &Client{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *Client) GetLastContractsByObjectIDs(ctx goctx.Context, objectIDs []int) ([]Contract, error) {
	var response []Contract
	requestURL := fmt.Sprintf("%s/contracts/objects/last", c.baseURL)

	values := url.Values{}
	for _, objectID := range objectIDs {
		values.Add("id", strconv.Itoa(objectID))
	}
	if len(values) > 0 {
		requestURL = fmt.Sprintf("%s?%s", requestURL, values.Encode())
	}

	status, err := c.client.DoJson(ctx, http.MethodGet, requestURL, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("c.client.DoJson: %w", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	return response, nil
}
