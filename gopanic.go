package gopanic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// APIURL is the API base URL
	APIURL = "https://cryptopanic.com/api"
	// APIVersion is the API version
	APIVersion = "v1"
	// APIUserAgent is the User-Agent header sent with API requests
	APIUserAgent = "GoPanic +github.com/bfontaine/gopanic"
)

var (
	// ErrNoNext represents the error when a response has no next page
	ErrNoNext error = errors.New("Response has no next page")
	// ErrNoPrevious represents the error when a response has no previous page
	ErrNoPrevious error = errors.New("Response has no previous page")
	// ErrBadToken represents the error when the auth token is invalid
	ErrBadToken error = errors.New("Unrecognized auth token")
	// ErrProOnly represents the error when the query needs you to be PRO
	ErrProOnly error = errors.New("You need to be PRO to get extra metadata")
	// ErrApprovedOnly represents the error when the query needs you to be approved
	ErrApprovedOnly error = errors.New("You need to be Approved to get original sources")
)

// API represents a CryptoPanic API client instance
type API struct {
	authToken string
	client    *http.Client

	// Set this to true if you have a PRO account and want extra metadata on
	// posts such as image and description fields.
	ExtraMetadata bool
	// Set this to true if your API is recognized as an approved Partners API
	// in order to get the original source in posts.
	// See https://cryptopanic.com/developers/api/partners/
	OriginalSource bool
}

// New creates a new API
func New(authToken string) *API {
	return &API{
		authToken: authToken,
		client:    http.DefaultClient,
	}
}

// Posts returns a response containing a sequence of posts
func (api *API) Posts() (*PostsResponse, error) {
	return api.FilteredPosts(Filter{})
}

// News is like Posts but it only returns posts of type 'news' (no 'media')
func (api *API) News() (*PostsResponse, error) {
	return api.FilteredPosts(Filter{
		Kind: "news",
	})
}

// FilteredPosts is like Posts but with filters
func (api *API) FilteredPosts(p Filter) (*PostsResponse, error) {
	url := api.makeURL("/posts/", p.encode())

	return api.postsCall(url)
}

// FilteredNews is like FilteredPosts with Kind set to "news"
func (api *API) FilteredNews(p Filter) (*PostsResponse, error) {
	p.Kind = "news"
	return api.FilteredPosts(p)
}

// Portfolio returns a response containing a user and a portfolio
func (api *API) Portfolio() (*PortfolioResponse, error) {
	url := api.makeURL("/portfolio/", url.Values{})

	resp, err := api.call(url, &PortfolioResponse{})
	if err == nil {
		pResp := resp.(*PortfolioResponse)
		pResp.BaseResponse.api = api
		return pResp, pResp.Error()
	}

	return resp.(*PortfolioResponse), err
}

// makeURL creates an URL given a path and params.
//
// The function might modify the 'param' map.
func (api *API) makeURL(path string, params url.Values) string {
	if params.Get("auth_token") == "" {
		params.Set("auth_token", api.authToken)
	}

	if api.ExtraMetadata && params.Get("metadata") == "" {
		params.Set("metadata", "true")
	}

	if api.OriginalSource && params.Get("approved") == "" {
		params.Set("approved", "true")
	}

	return fmt.Sprintf("%s/%s%s?%s", APIURL, APIVersion, path, params.Encode())
}

func (api *API) postsCall(url string) (*PostsResponse, error) {
	resp, err := api.call(url, &PostsResponse{})
	if err == nil {
		pResp := resp.(*PostsResponse)
		pResp.BaseResponse.api = api

		return pResp, pResp.Error()
	}
	return resp.(*PostsResponse), err
}

func (api *API) call(url string, jsonResp interface{}) (interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", APIUserAgent)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, jsonResp)
	if err != nil {
		return nil, err
	}

	return jsonResp, nil
}
