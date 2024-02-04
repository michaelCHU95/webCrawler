package crawler

import (
	"fmt"
	"io"
	"sync"
)

type Crawler struct {
	Sites []string
	Out   io.Writer
}

// NewCrawler is the constructor of Crawler
func NewCrawler(urls []string, writer io.Writer) *Crawler {
	c := new(Crawler)
	c.Sites = append(c.Sites, urls...)
	c.Out = writer
	return c
}

// Run is the main function to run the crawler
func (impl *Crawler) Run() {
	if len(impl.Sites) == 0 {
		return
	}

	wg := new(sync.WaitGroup)
	results := make(chan []string, len(impl.Sites))

	// Crawling process
	for _, s := range impl.Sites {
		wg.Add(1)
		go func(rootURL string) {
			worker := InitWorker()
			worker.Start(rootURL, results)
			wg.Done()
		}(s)
	}
	wg.Wait()
	close(results)

	// Write results
	for r := range results {
		for _, url := range r {
			fmt.Fprintln(impl.Out, url)
		}
	}
}
