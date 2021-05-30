package main

import (
  "fmt"
	"net/url"
	"net/http"
	"golang.org/x/net/html"
  "strings"
)

type FetchedDocument struct {
  links []string
  statics []string
  err error
  link_docs []FetchedDocument
}

type DocumentFetcher interface {
	// Fetch returns the links and URLs
	// a slice of URLs found on that page.
	fetch_document(url_s string) (d FetchedDocument)
}


func (d FetchedDocument) fetch_document(url_s string) FetchedDocument {
    targetURL, err := url.Parse(url_s)
    if err != nil {
  		fmt.Println("couldn't parse that URL:", err)
  		return FetchedDocument{err:err}
  	}
    var domain = targetURL.Host
    var protocol = "http" // default value
    // get the protocol and domain for Relative URL
    if idx := strings.IndexByte(url_s, ':'); idx >= 0 {
      protocol = url_s[:idx]
    }
    resp, err := http.Get(url_s)
    if err != nil {
  		fmt.Println("failed to get URL %s: %v", url_s, err)
  		return FetchedDocument{err:err}
  	}
  	defer resp.Body.Close()
  	contentType := resp.Header.Get("Content-Type")
  	if contentType != "" && !strings.HasPrefix(contentType, "text/html") { // "" to allow for no header being sent
  		return FetchedDocument{err:nil}
  	}
    fmt.Println("After getting page of url::", targetURL)
    tokens := html.NewTokenizer(resp.Body)
    var links []string
    var statics []string
  	for {
  		tokenType := tokens.Next()
  		if tokenType == html.ErrorToken { //an EOF
        fmt.Println("End token")
  			break
  		}
  		token := tokens.Token()
  		if tokenType == html.StartTagToken { //opening tag
        var tag_s = token.DataAtom.String()
  			switch tag_s {
  			case "a", "link": //link tags
  				for _, attr := range token.Attr {
  					if attr.Key == "href" {
                var child_url = parse_link_token(attr.Val, protocol, domain)
                if tag_s == "a" {
                  links = append(links, child_url)
                } else {
                  statics = append(statics, child_url)
                }
  					}
  				}
  			}
  		}
  	}
    return FetchedDocument{links: links, statics: statics, err: nil}
}


func parse_link_token(url_s string, protocol string, domain string) string {
  var child_url = url_s
  if strings.HasPrefix(url_s, "/") {
    // fmt.Println("From Relative URL::", child_url)
    child_url = protocol+ "://"+domain+child_url
    // fmt.Println("Relative to absolute URL::", child_url)
  }
  return child_url
}
