// Referred https://github.com/google/go-github
// BSD-3-Clause https://github.com/google/go-github/blob/master/LICENSE

package informal

import (
	"fmt"
	"net/url"

	"github.com/kenchan0130/go-jamf-pro/jamf"
)

type service struct {
	client *Client

	username string
	password string
}

type services struct {
	DistributionFileUpload *DistributionFileUploadService
	Session                *SessionService
}

type Client struct {
	*jamf.BaseClient

	common service
	services
}

func NewClient(serverURL string, username string, password string) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(): %v", err)
	}

	client := jamf.NewBaseClient(u)

	c := &Client{
		BaseClient: client,
	}

	c.common.client = c
	c.common.username = username
	c.common.password = password
	c.DistributionFileUpload = (*DistributionFileUploadService)(&c.common)
	c.Session = (*SessionService)(&c.common)

	return c, nil
}
