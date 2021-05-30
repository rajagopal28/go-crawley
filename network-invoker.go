package main

import (
  "fmt"
	"net/url"
	"net/http"
	"golang.org/x/net/html"
	"os"
  "strings"
  "sync"
)

type VisitTracker struct {
  visited_map   map[string]bool
	lock sync.Mutex
}

func (t VisitTracker)checkvisited(url string)bool{
	t.lock.Lock()
	defer t.lock.Unlock()
	_,ok:=t.visited_map[url]
	if !ok {
		t.visited_map[url]=true
		return false
	}
	return true
}


func crawlPage(url_s string, max_level int) {
  fmt.Printf("crawling %s\n", url_s)
  // crawl_recursive(url_s, max_level, false, seen_urls)
  tracker := VisitTracker{visited_map: make(map[string]bool)}
  doc := crawl_master(url_s, max_level, false, tracker)
  fmt.Println(tracker)

  fmt.Println("\nDocument is :\n")
  fmt.Println(doc)
}

func crawl_master(url_s string, level int, external_crawl bool, tracker VisitTracker) FetchedDocument {
   parent_doc := make(chan FetchedDocument)
   var d = FetchedDocument{}
   go crawl_concurrent(url_s, level, external_crawl, tracker, d, parent_doc)
   // defer close(parent_doc)
   val, ok := <- parent_doc
   if ok {
     return val
   }
   return FetchedDocument{err: nil}
}

func domain(url_s string) string {
  targetURL, err := url.Parse(url_s)
  if err != nil {
    fmt.Println("couldn't parse that URL:", err)
    return ""
  }
  return targetURL.Host
}

func crawl_concurrent(url_s string, level int, external_crawl bool, tracker VisitTracker, fetcher DocumentFetcher, document_p chan FetchedDocument) {
  if level > 0 && !tracker.checkvisited(url_s) {
    document := fetcher.fetch_document(url_s)
  	if document.err != nil {
  		fmt.Println(document.err)
  		return
  	}
  	// fmt.Println("found: ", document)
    urls := document.links
  	for _, u := range urls {
      domain_p := domain(url_s)
      child_url := prune_url(u)
      domain_c := domain(child_url)
      if external_crawl || domain_c == domain_p {
        child_doc := make(chan FetchedDocument)
        fmt.Println("found child_url: ", u)
        fmt.Println("Spawing child crawl")
    		go crawl_concurrent(child_url, level-1, external_crawl, tracker, fetcher, child_doc)
        fmt.Println("Defering child close: ")
        // defer close(child_doc)
        document.link_docs = append(document.link_docs, <-child_doc)
        fmt.Println("Setting up child documents")
      }
  	}
    document_p <- document
  }
  defer close(document_p)
}

func prune_url(url_s string) string {
  index := strings.Index(url_s, "#")
  if index != -1 {
    url_s = url_s[index:]
  }
  index = strings.Index(url_s, "?")
  if index != -1 {
    url_s = url_s[index:]
  }
  return url_s
}
