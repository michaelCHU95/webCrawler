package crawler

import (
	"net/url"
)

type Worker struct {
	Visited map[string]struct{}
	Result  urlResults
	root    *url.URL

	// http util methods
	GetUrlResponse      GetUrlResponseFunc
	ParseHTMLToGetLinks ParseHTMLToGetLinksFunc
}

type urlResults []string

func (u *urlResults) append(val string) {
	*u = append(*u, val)
}

func InitWorker(rootURL string) (*Worker, error) {
	w := new(Worker)
	w.Visited = map[string]struct{}{}
	w.GetUrlResponse = GetUrlResponse
	w.ParseHTMLToGetLinks = ParseHTMLToGetLinks

	u, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	w.root = u
	return w, nil
}

func (w *Worker) Start(output chan []string) {
	w.fetchLinks(w.root.String())
	output <- w.Result
}

func (w *Worker) fetchLinks(url string) {
	if _, ok := w.Visited[url]; ok {
		return
	}
	w.Visited[url] = struct{}{}
	w.Result.append(url)

	res, status, err := w.GetUrlResponse(url)
	if err != nil {
		return
	}
	if status >= 400 {
		// TODO: add warning logs
		return
	}

	links, err := w.ParseHTMLToGetLinks(res)
	if err != nil {
		// TODO: add warning logs
		return
	}

	for _, link := range links {
		link = w.convertRelativeUrlToAbsolute(link)
		w.fetchLinks(link)
	}
}

// Convert relative url to absolute link
func (w *Worker) convertRelativeUrlToAbsolute(url string) string {
	urlObj := *w.root
	urlObj.Path = url
	return urlObj.String()
}
