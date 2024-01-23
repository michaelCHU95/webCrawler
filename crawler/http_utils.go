package crawler

import (
	"bytes"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type GetUrlResponseFunc = func(url string) (response []byte, status int, err error)
type ParseHTMLToGetLinksFunc = func(b []byte) (*html.Node, error)

// GetUrlResponse call GET request to get HTTP response in bytes by given URL
func GetUrlResponse(url string) (response []byte, status int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return response, resp.StatusCode, err
	}

	// Check if there is an error HTTP code
	if resp.StatusCode >= 400 {
		return response, resp.StatusCode, err
	}

	defer resp.Body.Close()

	response, _ = io.ReadAll(resp.Body)
	return response, resp.StatusCode, nil
}

// getLinks traverses through HTML nodes to search for <a>
func getLinks(node *html.Node) (links []string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		getLinks(c)
	}
	return links
}

// ParseHTML parsing bytes into HTML Nodes and get links
func ParseHTMLToGetLinks(b []byte) (links []string, err error) {
	node, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	links = getLinks(node)
	return links, nil
}
