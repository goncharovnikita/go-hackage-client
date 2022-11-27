package gohackageclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type packageGetter struct {
	baseURL    string
	httpClient *http.Client
}

func NewPackageGetter() *packageGetter {
	return &packageGetter{
		baseURL:    "https://hackage.haskell.org",
		httpClient: http.DefaultClient,
	}
}

func (c *packageGetter) GetPackageVersions(
	ctx context.Context,
	name string,
) (PackageVersions, error) {
	r := make(PackageVersions)

	url := fmt.Sprintf("%s/package/%s.json", c.baseURL, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("non-ok response: %s", res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *packageGetter) GetPackageCabalReader(
	ctx context.Context,
	name string,
	version string,
) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/package/%s-%s/%s.cabal", c.baseURL, name, version, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("non-ok response: %s", res.Status)
	}

	return res.Body, nil
}
