package docker

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
)

const (
	authService = "registry.docker.io"
	authScope   = "repository:meik99/cloud-computing:pull"

	tokenUrl = "https://auth.docker.io/token?" +
		"service=" + authService +
		"&scope=" + authScope

	manifestUrl           = "https://registry.hub.docker.com/v2/meik99/cloud-computing/manifests/dev"
	manifestContentType   = "application/vnd.docker.distribution.manifest.v2+json"
	manifestAuthorization = "Bearer %s"

	headerAccept        = "Accept"
	headerAuthorization = "Authorization"
)

type authToken struct {
	token string
}

func GetDigest() (string, error) {
	imageManifest, err := getManifest()

	if err != nil {
		return "", errors.WithStack(err)
	}

	return imageManifest.conf.digest, nil
}

func getManifest() (manifest, error) {
	token, err := getToken()

	if err != nil {
		return manifest{}, errors.WithStack(err)
	}

	request, err := createRequest(token)

	if err != nil {
		return manifest{}, errors.WithStack(err)
	}

	response, err := (&http.Client{}).Do(request)

	if err != nil {
		return manifest{}, errors.WithStack(err)
	}

	var imageManifest manifest
	err = parseResponse(response, &imageManifest)
	return imageManifest, errors.WithStack(err)
}

func createRequest(token string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, manifestUrl, nil)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	request.Header.Set(headerAccept, manifestContentType)
	request.Header.Set(headerAuthorization, formatAuthorizationUrl(token))
	return request, nil
}

func formatAuthorizationUrl(token string) string {
	return fmt.Sprintf(manifestAuthorization, token)
}

func getToken() (string, error) {
	response, err := http.Get(tokenUrl)

	if err != nil {
		return "", errors.WithStack(err)
	}

	var result authToken
	err = parseResponse(response, &result)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return result.token, nil
}

func parseResponse(response *http.Response, v interface{}) error {
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return errors.WithStack(err)
	}

	err = json.Unmarshal(data, v)
	return errors.WithStack(err)
}
