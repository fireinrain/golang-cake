package main

import (
	"context"
	"fmt"
	"github.com/shadowscatcher/shodan"
	"github.com/shadowscatcher/shodan/search"
	"log"
	"net/http"
	"os"
)

func main() {
	//port:443 http.status:403 country:KR ssl:cloudflare
	nginxSearch := search.Params{
		Page: 1,
		Query: search.Query{
			Port: 443,
			HTTP: search.HTTP{
				Status: 403,
			},
			ASN:     "AS31898",
			SSL:     "cloudflare",
			Country: "KR",
		},
	}

	client, _ := shodan.GetClient(os.Getenv("SHODAN_API_KEY"), http.DefaultClient, true)
	ctx := context.Background()
	result, err := client.Search(ctx, nginxSearch)
	if err != nil {
		log.Fatal(err)
	}

	for _, match := range result.Matches {
		// a lot of returned data can be used in another searches
		// it's easy because you will get response with almost all possible fields, just don't forget to check them
		fmt.Println(match.IP)
	}

	// later on you can change every part of search query or parameters:
	//nginxSearch.Page++                            // for example, increase page
	//nginxSearch.Query.Port = 443                  // or add new search term
	//result, err = client.Search(ctx, nginxSearch) // and reuse modified parameters object
	//if err != nil {
	//	log.Fatal(err)
	//}
}
