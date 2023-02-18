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
	fmt.Fprint(os.Stderr, "USAGE: analyze <OUTPUT_FILE> <TWEETS_FILE_1> <TWEETS_FILE_2> <TWEETS_FILE_...>")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	pretty := flag.Bool("pretty", true, "Format JSON")
	fail := flag.Bool("fail", true, "Fail on error")
	flag.Parse()
	nargs := flag.NArg()
	if nargs < 2 {
		flag.Usage()
	}

	fw, err := os.OpenFile(flag.Arg(0), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	shared.FailOnError(err, "unable to open file: %s", flag.Arg(0))
	defer fw.Close()

	enc := json.NewEncoder(fw)
	if *pretty {
		enc.SetIndent("", "  ")
	}
	index := make(map[string]interface{})
	dc := 0
	total := 0

	for _, fp := range flag.Args()[1:] {
		f, err := os.OpenFile(fp, os.O_RDONLY, 0600)
		if err != nil {
			if *fail {
				shared.FailOnError(err, "unable to open file: %s", fp)
			} else {
				fmt.Fprintf(os.Stderr, "unable to open file: %s\nError: %v", fp, err)
				continue
			}
		}
		dec := json.NewDecoder(f)

		for {
			t := new(anaconda.Tweet)
			err := dec.Decode(t)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else {
					if *fail {
						shared.FailOnError(err, "unable to read tweet from file: %s", fp)
					} else {
						fmt.Fprintf(os.Stderr, "unable to read tweet from file: %s\nError: %v", fp, err)
						break
					}
				}
			}

			if _, ok := index[t.IdStr]; ok {
				dc++
				continue
			}

			index[t.IdStr] = true
			enc.Encode(t)
			total++
		}
	}
	fmt.Fprintln(os.Stdout, "Total:", total+dc)
	fmt.Fprintln(os.Stdout, "Duplicates:", dc)
	fmt.Fprintln(os.Stdout, "Actual:", total)
}
