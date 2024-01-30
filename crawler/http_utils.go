package crawler

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type GetUrlResponseFunc = func(url string) (response string, status int, err error)
type ParseHTMLToGetLinksFunc = func(b string) (links []string, err error)

// GetUrlResponse call GET request to get HTTP response in bytes by given URL
func GetUrlResponse(url string) (response string, status int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}

	// Check if there is an error HTTP code
	if resp.StatusCode >= 400 {
		return "", resp.StatusCode, err
	}

	defer resp.Body.Close()

	content, _ := io.ReadAll(resp.Body)
	return string(content), resp.StatusCode, nil
}

// getLinks traverses through HTML nodes to search for <a>
func getLinks(node *html.Node, links []string) []string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = getLinks(c, links)
	}

	return links
}

// ParseHTML parsing bytes into HTML Nodes and get links
func ParseHTMLToGetLinks(b string) (links []string, err error) {
	node, err := html.Parse(strings.NewReader(b))
	if err != nil {
		return nil, err
	}

	links = getLinks(node, []string{})
	return links, nil
}
