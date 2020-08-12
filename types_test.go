package gopanic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostVotes(t *testing.T) {
	p := Post{
		Votes: map[string]int{
			"negative":  1,
			"positive":  2,
			"important": 3,
			"liked":     4,
			"disliked":  5,
			"lol":       6,
			"toxic":     7,
			"saved":     8,
			"comments":  9,
		},
	}

	assert.Equal(t, 1, p.NegativeVotes())
	assert.Equal(t, 2, p.PositiveVotes())
	assert.Equal(t, 3, p.ImportantVotes())
	assert.Equal(t, 4, p.Likes())
	assert.Equal(t, 5, p.Dislikes())
	assert.Equal(t, 6, p.LolVotes())
	assert.Equal(t, 7, p.ToxicVotes())
	assert.Equal(t, 8, p.Saved())
	assert.Equal(t, 9, p.Comments())
}

func TestPostCurrencyCodes(t *testing.T) {
	p := Post{
		Currencies: []Currency{
			Currency{Code: "ABC"},
			Currency{Code: "DEF"},
			Currency{Code: "GHI"},
		},
	}
	assert.Equal(t, []string{}, Post{}.CurrencyCodes())
	assert.Equal(t, []string{"ABC", "DEF", "GHI"}, p.CurrencyCodes())
}

/*
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
*/

func TestFilterEncode(t *testing.T) {
	assert.Equal(t, "", Filter{}.encode().Encode())
	assert.Equal(t, "public=true", Filter{Public: true}.encode().Encode())
	assert.Equal(t, "filter=xyz", Filter{UIFilter: "xyz"}.encode().Encode())
	assert.Equal(t, "kind=xyz", Filter{Kind: "xyz"}.encode().Encode())

	assert.Equal(t, "currencies=ABC",
		Filter{
			CurrencyCodes: []string{"ABC"},
		}.encode().Encode())

	assert.Equal(t, "currencies=ABC%2CDEF%2CGHI",
		Filter{
			CurrencyCodes: []string{"ABC", "DEF", "GHI"},
		}.encode().Encode())

	assert.Equal(t, "regions=it",
		Filter{
			Regions: []string{"it"},
		}.encode().Encode())

	assert.Equal(t, "regions=it%2Cen%2Cde",
		Filter{
			Regions: []string{"it", "en", "de"},
		}.encode().Encode())

	assert.Equal(t, "currencies=BTC%2CETH&filter=hot&kind=media&public=true&regions=en%2Cde%2Cfr",
		Filter{
			CurrencyCodes: []string{"BTC", "ETH"},
			UIFilter:      "hot",
			Kind:          "media",
			Public:        true,
			Regions:       []string{"en", "de", "fr"},
		}.encode().Encode())
}

func TestPostsResponseHasNext(t *testing.T) {
	assert.True(t, PostsResponse{next: "http://"}.HasNext())
	assert.False(t, PostsResponse{next: ""}.HasNext())
}

func TestPostsResponseHasPrevious(t *testing.T) {
	assert.True(t, PostsResponse{previous: "http://"}.HasPrevious())
	assert.False(t, PostsResponse{previous: ""}.HasPrevious())
}

func TestPostsResponseError(t *testing.T) {
	assert.Nil(t, PostsResponse{}.Error())
}
