# GoPanic

**GoPanic** is an unofficial Go client for the [CryptoPanic API][api].

[api]: https://cryptopanic.com/developers/api/

## Example usage

```golang
package main

import (
	"fmt"
	"strings"

	"github.com/bfontaine/gopanic"
)

func main() {
	api := gopanic.New("<your token>")

    // Get all posts
    resp, err := api.Posts()
    if err != nil {
        log.Fatal(err)
    }

	for _, post := range resp.Posts {
		fmt.Printf("# [%s] %s\n",
			strings.Join(post.CurrencyCodes(), ", "),
			post.Title)
	}
    // prints something like:
    // # [] Blablabla
    // # [BTC] Something Bitcoin Something
    // # [ETC, DAI] Something ETC+DAI
    // # [BTC] Another BTC something
    // ...

    // Same, but exclude media posts
    resp, err = api.News()

    // Same, but get only news about BTC
	resp, err = api.FilteredNews(gopanic.Filter{
        CurrencyCodes: []string{"BTC"},
    })

    // You can also combine filters:
	resp, err = api.FilteredNews(gopanic.Filter{
        // News about BTC and/or ETH
        CurrencyCodes: []string{"BTC", "ETH"},

        // using the UI filter 'rising'
        UIFilter: "rising",

        // in the "regions" (= languages) 'en' and 'de'
        Regions: []string{"en", "de"},
    })

    // You can set additional options on the API:
    // - Set this *only* if you have a PRO account and want the extra metadata
    api.ExtraMetadata = true
    // - Set this *only* if you have an approved API access and want the source URLs
    api.Approved = true
    // Setting these without meeting the requirements will result in errors on
    // all your calls: ErrProOnly and ErrApprovedOnly.
}
```

## Setup

Sign-up or log-in to CryptoPanic then go to the [API page][api] to retrieve
your API auth token.
