package shared

import (
	"fmt"
	"os"
	"sort"

	goenvtostruct "github.com/abdullah2993/go-env-to-struct"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/exp/constraints"
)

type config struct {
	AccessToken       string `dotenv:"ACCESS_TOKEN"`
	AccessTokenSecret string `dotenv:"ACCESS_TOKEN_SECRET"`
	ConsumerKey       string `dotenv:"CONSUMER_KEY"`
	ConsumerSecret    string `dotenv:"CONSUMER_SECRET"`
}

var Config = goenvtostruct.GetConfigFromEnv[config]()

func FailOnError(err error, message string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, message, args...)
		fmt.Fprintf(os.Stderr, "\nError: %v", err)
		os.Exit(1)
	}
}

func Min[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

type Pair[K comparable, V constraints.Integer] struct {
	Key   K
	Value V
}

type PairList[K comparable, V constraints.Integer] []Pair[K, V]

func (p PairList[K, V]) Len() int           { return len(p) }
func (p PairList[K, V]) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList[K, V]) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type FreqMap[K comparable, V constraints.Integer] map[K]V

func (r FreqMap[K, V]) Add(key K) {
	if count, ok := r[key]; ok {
		r[key] = count + 1
	} else {
		r[key] = 1
	}
}

func (r FreqMap[K, V]) Count(key K) (V, bool) {
	v, ok := r[key]
	return v, ok
}

func (r FreqMap[K, V]) Pairs() PairList[K, V] {
	pl := make(PairList[K, V], len(r))
	i := 0
	for k, v := range r {
		pl[i] = Pair[K, V]{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (r FreqMap[K, V]) DecSorted() PairList[K, V] {
	pl := r.Pairs()
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (r FreqMap[K, V]) AscSorted() PairList[K, V] {
	pl := r.Pairs()
	sort.Sort(pl)
	return pl
}

func (r FreqMap[K, V]) Top(count int) PairList[K, V] {
	pl := r.DecSorted()
	return pl[:Min(len(pl), count)]
}

func (r FreqMap[K, V]) Len() int {
	return len(r)
}

func NewFreqMap[K comparable, V constraints.Integer]() FreqMap[K, V] {
	return make(map[K]V)
}
