package informal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kenchan0130/go-jamf-pro/utils"
)

type DistributionFileUploadService service

type DistributionFileUploadDestination string

const (
	DistributionFileUploadDestinationDefault DistributionFileUploadDestination = "0"
)

type DistributionFileUploadFileType string

const (
	DistributionFileUploadFileTypePackage    DistributionFileUploadFileType = "0"
	DistributionFileUploadFileTypeEbook      DistributionFileUploadFileType = "1"
	DistributionFileUploadFileTypeInHouseApp DistributionFileUploadFileType = "2"
)

const distributionFileUploadPath = "/dbfileupload"

func (s *DistributionFileUploadService) Upload(ctx context.Context, packageID int, packageName string, fileType DistributionFileUploadFileType, destination DistributionFileUploadDestination, src io.Reader) (*utils.Response, error) {
	body := new(bytes.Buffer)

	if _, err := io.Copy(body, src); err != nil {
		return nil, fmt.Errorf("io.Copy(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, utils.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: utils.Uri{
			Entity: distributionFileUploadPath,
		},
		ContentType: "application/octet-stream",
		RequestMiddlewareFunc: func(r *http.Request) {
			r.SetBasicAuth(s.username, s.password)
			r.Header.Set("DESTINATION", string(destination))
			r.Header.Set("OBJECT_ID", fmt.Sprint(packageID))
			r.Header.Set("FILE_TYPE", string(fileType))
			r.Header.Set("FILE_NAME", packageName)
		},
		Body: body,
	})

	if err != nil {
		return nil, fmt.Errorf("client.Post(): %v", err)
	}

	return resp, nil
}
