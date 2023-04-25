// Referred https://github.com/google/go-github
// BSD-3-Clause https://github.com/google/go-github/blob/master/LICENSE

package jamfproapi

import (
	"fmt"
	"net/url"
	"path"

	"github.com/kenchan0130/go-jamf-pro/jamf"
)

type service struct {
	client *Client
}

type services struct {
	APIAuthentication *APIAuthenticationService
	Categories        *CategoriesService
	Icon              *IconService
	Scripts           *ScriptsService
	SSOFailover       *SSOFailoverService
}

type Client struct {
	*jamf.BaseClient

	common service
	services
}

type ListOptions struct {
	Page     *int      `url:"page,omitempty"`
	PageSize *int      `url:"page-size,omitempty"`
	Sort     *[]string `url:"sort,omitempty" del:","`
	Filter   *string   `url:"filter,omitempty"`
}

const apiEndpointPath = "/api"

func NewClient(serverURL string) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(): %v", err)
	}
	u.Path = path.Join(u.Path, apiEndpointPath)

	client := jamf.NewBaseClient(u)

	c := &Client{
		BaseClient: client,
	}

	c.common.client = c

	c.APIAuthentication = (*APIAuthenticationService)(&c.common)
	c.Categories = (*CategoriesService)(&c.common)
	c.Icon = (*IconService)(&c.common)
	c.Scripts = (*ScriptsService)(&c.common)
	c.SSOFailover = (*SSOFailoverService)(&c.common)

	return c, nil
}
