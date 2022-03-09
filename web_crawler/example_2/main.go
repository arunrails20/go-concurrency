package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/html"
)

var fetched map[string]bool

type result struct {
	url   string
	urls  []string
	err   error
	depth int
}

func main() {
	// Use all cores in the machine
	runtime.GOMAXPROCS(runtime.NumCPU())
	fetched = make(map[string]bool)
	now := time.Now()
	Crawl("http://google.com", 2)
	fmt.Println("time taken:", time.Since(now))
}

func Crawl(url string, depth int) {
	// creating channel to communicates between the goroutines
	// we use result struct to communicate.
	ch := make(chan *result)

	fetch := func(url string, depth int) {
		urls, err := findLinks(url)
		// pass the results to the channel
		ch <- &result{url, urls, err, depth}
	}

	go fetch(url, depth)
	// mark the current url as true
	fetched[url] = true

	for fetching := 1; fetching > 0; fetching-- {
		res := <-ch

		if res.err != nil {
			continue
		}

		fmt.Printf("found: %s\n", res.url)

		if res.depth > 0 {
			for _, u := range res.urls {
				if !fetched[u] {
					fetching++
					go fetch(u, res.depth-1)
					fetched[u] = true
				}
			}
		}
	}
	close(ch)
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
