package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

var fetched map[string]bool

func main() {
	fetched = make(map[string]bool)
	now := time.Now()
	Crawl("http://google.com", 2)
	fmt.Println("time taken:", time.Since(now))
}

// Crawl uses findlinks to recursively Crawl
func Crawl(url string, depth int) {
	if depth < 0 {
		return
	}

	urls, err := findLinks(url)

	if err != nil {
		return
	}
	fmt.Printf("found: %s\n", url)
	// mark the url as true to stop visit again and again
	fetched[url] = true
	for _, u := range urls {
		if !fetched[u] {
			Crawl(u, depth-1)
		}
	}
	return
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
