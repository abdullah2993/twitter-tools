package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/abdullah2993/twitter-tools/pkgs/shared"
)

func usage() {
	fmt.Fprint(os.Stderr, "USAGE: search <SEARCH_TERM> <OUTPUT_FILE>")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage

	pretty := flag.Bool("pretty", true, "Format JSON")

	flag.Parse()
	nargs := flag.NArg()

	if nargs != 2 {
		flag.Usage()
	}

	api := anaconda.NewTwitterApiWithCredentials(shared.Config.AccessToken, shared.Config.AccessTokenSecret, shared.Config.ConsumerKey, shared.Config.ConsumerSecret)
	v := url.Values{}
	v.Add("count", "100")
	res, err := api.GetSearch(flag.Arg(0), v)
	shared.FailOnError(err, "unable to search via api")

	f, err := os.OpenFile(flag.Arg(1), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	shared.FailOnError(err, "unable to open file: %s", flag.Arg(1))
	defer f.Close()
	enc := json.NewEncoder(f)
	if *pretty {
		enc.SetIndent("", "  ")
	}
	count := 0
	for {
		if err != nil {
			shared.FailOnError(err, "unable to search via api")
		}
		for _, tw := range res.Statuses {
			err := enc.Encode(tw)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to save tweet\nError: %v\n", err)
				continue
			}
			count++
		}
		fmt.Println("Tweets gathered", count)
		res, err = res.GetNext(api)
	}
}
