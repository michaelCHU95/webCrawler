package main

import (
	"os"

	"github.com/michaelCHU95/webCrawler/crawler"
)

func main() {
	sites := []string{"https://monzo.com/"}
	file, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	crawler := crawler.NewCrawler(sites, file)
	crawler.Run()
}
