package main

import "fmt"

func main() {
  fmt.Println("crawling page")
	crawlPage("https://monzo.com" , 2)
}

func crawl(f DocumentFetcher) {
  fmt.Println("Cast based crawl")
  var ret = f.fetch_document("https://monzo.com")
  fmt.Println(ret)
}
