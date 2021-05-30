package main

import "fmt"

func main() {
  fmt.Println("crawling page")
	doc := crawlPage("https://monzo.com" , 2, false)
  recursive_print_tree(doc, 0)
}

func crawl(f DocumentFetcher) {
  fmt.Println("Cast based crawl")
  var ret = f.fetch_document("https://monzo.com")
  fmt.Println(ret)
}

func recursive_print_tree(doc FetchedDocument, level int) {
  if doc.err == nil {
    links := doc.links
    links_docs := doc.link_docs
    statics := doc.statics
    fmt.Println("|");
    for i := 0; i < level; i++ {
  		fmt.Printf("----")
  	}
    for _, v := range statics {
      fmt.Println("-->STATICS:", v)
    }
    for i, v := range links {
        fmt.Println("-->LINK:", v)
      if len(links_docs) > i {
        defer recursive_print_tree(links_docs[i], level+1)
      }
    }
  }
}
