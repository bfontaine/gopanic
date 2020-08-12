# GoPanic

**GoPanic** is an unofficial Go client for the [CryptoPanic API][api].

[Docs](https://pkg.go.dev/github.com/bfontaine/gopanic?tab=doc)

[api]: https://cryptopanic.com/developers/api/

## Example usage

### Posts

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
		// See https://pkg.go.dev/github.com/bfontaine/gopanic?tab=doc#Post
		// for all the fields and methods you can use on Post instances
	}
	// prints something like:
	//   # [] Blablabla
	//   # [BTC] Something Bitcoin Something
	//   # [ETC, DAI] Something ETC+DAI
	//   # [BTC] Another BTC something
	//   ...

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

### Portfolio

```golang
package main

import (
	"fmt"

	"github.com/bfontaine/gopanic"
)

func main() {
	api := gopanic.New("<your token>")

    // Get your portfolio
	resp, err := api.Portfolio()
	if err != nil {
		log.Fatal(err)
	}

    user := resp.User
    portfolio := resp.Portfolio

	fmt.Printf("# User: %s (%s)\n", user.Username, user.Email)
    // print something like:
    //   # foo42 (foo@example.com)

	fmt.Printf("# Portfolio (%s):\n", portfolio.CurrencyCode)
	fmt.Println("Totals:")
	for code, amount := range portfolio.Totals {
		fmt.Printf("- %s: %s\n", code, amount)
	}
    // prints something like:
    //   # Portfolio (USD):
    //   Totals:
    //   - ETH: 123.4567890
    //   - USD: 456.78
    //   - BTC: 0.012345600

	fmt.Println("Entries:")
	for _, entry := range portfolio.Entries {
		fmt.Printf("- %f %s (%.2f%%)\n", entry.Amount, entry.Currency.Code, entry.Percentage)
	}
    // prints something like:
    //   Entries:
    //   - 0.01234 BTC (42.34%)
    //   - 0.56789 ETH (34.56%)
    //   - 50.00000 ALGO (3.23%)
    //   ...
```

## Setup

Sign-up or log-in to CryptoPanic then go to the [API page][api] to retrieve
your API auth token.
