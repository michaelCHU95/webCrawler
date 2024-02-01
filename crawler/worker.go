package crawler

type Worker struct {
	Visited map[string]struct{}
	Result  urlResults

	// http util methods
	GetUrlResponse      GetUrlResponseFunc
	ParseHTMLToGetLinks ParseHTMLToGetLinksFunc
}

type urlResults []string

func (u *urlResults) append(val string) {
	*u = append(*u, val)
}

func InitWorker() *Worker {
	w := new(Worker)
	w.Visited = map[string]struct{}{}
	w.GetUrlResponse = GetUrlResponse
	w.ParseHTMLToGetLinks = ParseHTMLToGetLinks
	return w
}

func (w *Worker) Start(root string, output chan []string) {
	w.fetchLinks(root)
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
		return
	}

	links, err := w.ParseHTMLToGetLinks(res)
	if err != nil {
		return
	}

	for _, link := range links {
		w.fetchLinks(link)
	}
}
