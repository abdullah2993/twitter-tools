package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/abdullah2993/twitter-tools/pkgs/shared"
)

func usage() {
	fmt.Fprint(os.Stderr, "USAGE: analyze <TWEETS_FILE> <REPORT_FILE>")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	topCount := flag.Int("top", 30, "Number of results to print(less than 1 to print all)")

	flag.Parse()
	nargs := flag.NArg()
	if nargs != 2 {
		flag.Usage()
	}

	rf, err := os.Create(flag.Arg(1))
	shared.FailOnError(err, "unable to open/create report file")
	defer rf.Close()

	f, err := os.OpenFile(flag.Arg(0), os.O_RDONLY, 0600)
	shared.FailOnError(err, "unable to open file: %s", flag.Arg(0))
	defer f.Close()

	totalFreq := shared.NewFreqMap[string, int]()
	originalFreq := shared.NewFreqMap[string, int]()
	retweetFreq := shared.NewFreqMap[string, int]()

	userDetails := make(map[string]anaconda.User)

	total := 0
	retweets := 0
	originals := 0

	dec := json.NewDecoder(f)
	for {
		t := new(anaconda.Tweet)
		err := dec.Decode(t)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				shared.FailOnError(err, "unable to read tweet")
			}
		}
		userDetails[t.User.ScreenName] = t.User
		totalFreq.Add(t.User.ScreenName)
		total++
		if t.RetweetedStatus != nil {
			retweetFreq.Add(t.User.ScreenName)
			retweets++
		} else {
			originals++
			originalFreq.Add(t.User.ScreenName)
		}

	}

	rw := io.MultiWriter(os.Stdout, rf)

	fmt.Fprintln(rw, "Total Tweets: ", total)
	fmt.Fprintln(rw, "Total Original Tweets: ", originals)
	fmt.Fprintln(rw, "Total Retweets: ", retweets)
	fmt.Fprintln(rw, "Total Engagers: ", totalFreq.Len())
	fmt.Fprintln(rw, "Total Original Engagers: ", originalFreq.Len())
	fmt.Fprintln(rw, "Total Retweet Engagers: ", retweetFreq.Len())

	count := *topCount
	if count == 0 {
		count = totalFreq.Len()
	}
	fmt.Fprintln(rw, "------------")
	fmt.Fprintf(rw, "Top %d Engagers\n", count)
	for _, p := range totalFreq.Top(count) {
		fmt.Fprintln(rw, p.Value, p.Key)
	}

	count = *topCount
	if count == 0 {
		count = originalFreq.Len()
	}
	fmt.Fprintln(rw, "------------")
	fmt.Fprintf(rw, "Top %d Authors\n", count)
	for _, p := range originalFreq.Top(count) {
		fmt.Fprintln(rw, p.Value, p.Key)
	}

	count = *topCount
	if count == 0 {
		count = retweetFreq.Len()
	}
	fmt.Fprintln(rw, "------------")
	fmt.Fprintf(rw, "Top %d Retweets", count)
	for _, p := range retweetFreq.Top(count) {
		fmt.Fprintln(rw, p.Value, p.Key)
	}
}
