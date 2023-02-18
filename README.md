# twitter-tools

Small utilities for twitter

### search

Search for tweets with the given keyword/hashtag use to get all data related to a keyword/hashtag

```
USAGE: search <SEARCH_TERM> <OUTPUT_FILE>
Flags:
  -pretty
        Format JSON (default true)
```

will create a new output file if it doesn't exist otherwise will append

### track

Track tweets with the given keyword/hashtag in realtime usually used with search so you don't miss new tweets while the search is getting older tweets

```
USAGE: track <SEARCH_TERM> <OUTPUT_FILE>
Flags:
  -pretty
        Format JSON (default true)
```

will create a new output file if it doesn't exist otherwise will append

### merge

Deduplicate and merge multiple tweet files into one file, used to merge output of search and track utilities for analysis.

```
USAGE: analyze <OUTPUT_FILE> <TWEETS_FILE_1> <TWEETS_FILE_2> <TWEETS_FILE_...>
Flags:
  -fail
        Fail on error (default true)
  -pretty
        Format JSON (default true)
```

will create a new output file if it doesn't exist otherwise will append

### analyze

Do some preliminary analysis on a given tweets file

```
USAGE: analyze <TWEETS_FILE> <REPORT_FILE>
Flags:
  -top int
        Number of results to print(less than 1 to print all) (default 30)
```

will create a new report file if it doesn't exist otherwise it will override any existing file

### Usage

To analyze a given hashtag run `search` for the hashtag to get all the older tweets, in parallel run `track` to get the new tweets once done, merge the output files using `merge` and run `analyze` on the output file

search and track use twitter api keys so make sure to [set the environment variables](./.env.sample) or provide a `.env` in the same folder as the executable
