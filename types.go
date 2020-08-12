package gopanic

import (
	"fmt"
	"net/url"
	"strings"
)

// Post represents a post
type Post struct {
	// Kind can be either "news" or "media"
	Kind       string
	Domain     string
	Title      string
	Slug       string
	ID         int
	URL        string
	Votes      map[string]int
	Currencies []Currency
	Source     Source

	PublishedAt string `json:"published_at"`
	CreatedAt   string `json:"created_at"`

	// only for PRO users
	Image       string
	Description string
}

// Comments returns the comments count on the post
func (p Post) Comments() int { return p.Votes["comments"] }

// Saved returns the 'saved' count on the post
func (p Post) Saved() int { return p.Votes["saved"] }

// Likes returns the 'like' votes count on the post
func (p Post) Likes() int { return p.Votes["liked"] }

// Dislikes returns the 'dislike' votes count on the post
func (p Post) Dislikes() int { return p.Votes["disliked"] }

// PositiveVotes returns the 'positive' votes count on the post
func (p Post) PositiveVotes() int { return p.Votes["positive"] }

// NegativeVotes returns the 'negative' votes count on the post
func (p Post) NegativeVotes() int { return p.Votes["negative"] }

// ImportantVotes returns the 'important' votes count on the post
func (p Post) ImportantVotes() int { return p.Votes["important"] }

// LolVotes returns the 'lol' votes count on the post
func (p Post) LolVotes() int { return p.Votes["lol"] }

// ToxicVotes returns the 'toxic' votes count on the post
func (p Post) ToxicVotes() int { return p.Votes["toxic"] }

// CurrencyCodes returns an array of currency codes the post is about
func (p Post) CurrencyCodes() []string {
	codes := make([]string, len(p.Currencies))
	for i, currency := range p.Currencies {
		codes[i] = currency.Code
	}
	return codes
}

// Currency represents a currency
type Currency struct {
	Code  string
	Title string
	Slug  string
	URL   string
}

// Source represents a source
type Source struct {
	Title  string
	Region string
	Domain string
	Path   string
}

// Filter represents a way to provide filters to the Posts() method
type Filter struct {
	Public        bool
	UIFilter      string
	CurrencyCodes []string
	Regions       []string
	Kind          string
}

func (p Filter) encode() url.Values {
	v := url.Values{}

	if p.Public {
		v.Set("public", "true")
	}

	if p.UIFilter != "" {
		v.Set("filter", p.UIFilter)
	}

	if len(p.CurrencyCodes) > 0 {
		v.Set("currencies", strings.Join(p.CurrencyCodes, ","))
	}

	if len(p.Regions) > 0 {
		v.Set("regions", strings.Join(p.Regions, ","))
	}

	if p.Kind != "" {
		v.Set("kind", p.Kind)
	}

	return v
}

// Response represents a response from the API
type Response struct {
	TotalCount int    `json:"count"`
	Posts      []Post `json:"results"`

	next     string
	previous string

	// errors
	Status string
	Info   string

	api *API
}

// HasNext returns true if the response has a next page
func (r Response) HasNext() bool { return r.next != "" }

// HasPrevious returns true if the response has a previous page
func (r Response) HasPrevious() bool { return r.previous != "" }

// Next returns the next page of results, if any
func (r Response) Next() (*Response, error) {
	if !r.HasNext() {
		return nil, ErrNoNext
	}
	return r.api.call(r.next)
}

// Previous returns the previous page of results, if any
func (r Response) Previous() (*Response, error) {
	if !r.HasPrevious() {
		return nil, ErrNoPrevious
	}
	return r.api.call(r.previous)
}

// Error returns the error (if any) for this response
func (r Response) Error() error {
	// Statuses aren't documented;here are the ones I encountered:
	switch r.Status {
	case "":
		return nil
	case "Incomplete":
		switch r.Info {
		case "Token not found":
			return ErrBadToken
		default:
			return fmt.Errorf("Incomplete query: %s", r.Info)
		}
	case "invalid":
		switch r.Info {
		case "Metadata param requires PRO account":
			return ErrProOnly
		case "Access Denied. This API key is not approved for Partners API.":
			return ErrApprovedOnly
		default:
			return fmt.Errorf("Invalid query: %s", r.Info)
		}
	default:
		return fmt.Errorf("Error: %s", r.Info)
	}
}
