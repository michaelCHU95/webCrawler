package crawler

type Worker struct {
	Visited map[string]struct{}

	// http util methods
	GetUrlResponse      GetUrlResponseFunc
	ParseHTMLToGetLinks ParseHTMLToGetLinksFunc
}

func InitWorker() *Worker {
	w := new(Worker)
	w.Visited = map[string]struct{}{}
	w.GetUrlResponse = GetUrlResponse
	w.ParseHTMLToGetLinks = ParseHTMLToGetLinks
	return w
}

func (w *Worker) Start(root string, output chan []string) {
	result := []string{}
	w.fetchLinks(root, result)
	output <- result
}

func (w *Worker) fetchLinks(url string, result []string) {
	if _, ok := w.Visited[url]; ok {
		return
	}
	w.Visited[url] = struct{}{}

	res, status, err := w.GetUrlResponse(url)
	if err != nil {
		return
	}
	if status >= 400 {
		return
	}

	result = append(result, url)

	links, err := w.ParseHTMLToGetLinks(res)
	if err != nil {
		return
	}

	for _, link := range links {
		w.fetchLinks(link, result)
	}
}
