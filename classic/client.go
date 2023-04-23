// Referred https://github.com/google/go-github
// BSD-3-Clause https://github.com/google/go-github/blob/master/LICENSE

package classic

import (
	"fmt"
	"net/url"
	"path"

	"github.com/kenchan0130/go-jamf-pro/utils"
)

type service struct {
	client *Client
}

type services struct {
	ComputerExtensionAttributes *ComputerExtensionAttributesService
	ComputerGroups              *ComputerGroupsService
	OSXConfigurationProfiles    *OSXConfigurationProfilesService
	Packages                    *PackagesService
	Policies                    *PoliciesService
}

type Client struct {
	*utils.BaseClient

	common service
	services
}

const apiEndpointPath = "/JSSResource"

func NewClient(serverURL string) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(): %v", err)
	}
	u.Path = path.Join(u.Path, apiEndpointPath)

	client := utils.NewBaseClient(u)

	c := &Client{
		BaseClient: client,
	}

	c.common.client = c

	c.ComputerExtensionAttributes = (*ComputerExtensionAttributesService)(&c.common)
	c.ComputerGroups = (*ComputerGroupsService)(&c.common)
	c.OSXConfigurationProfiles = (*OSXConfigurationProfilesService)(&c.common)
	c.Packages = (*PackagesService)(&c.common)
	c.Policies = (*PoliciesService)(&c.common)

	return c, nil
}