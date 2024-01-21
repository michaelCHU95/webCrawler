package main

import (
	"bytes"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type GetUrlResponseFunc = func(url string) (response []byte, status int, err error)
type ParseHTMLFunc = func(b []byte) (*html.Node, error)

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

// ParseHTML parsing bytes into HTML Nodes
func ParseHTML(b []byte) (node *html.Node, err error) {
	node, err = html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return node, nil
}
