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

func crawl_recursive(url_s string, level int, external_crawl bool, seen_urls map[string]struct{}) error {
  var _, ok = seen_urls[url_s]
  if(level >= 0 && !ok) {
    targetURL, err := url.Parse(url_s)
    seen_urls[url_s] = struct{}{}
  	if err != nil {
  		fmt.Println("couldn't parse that URL:", err)
  		os.Exit(1)
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
  		return err
  	}
  	defer resp.Body.Close()
  	contentType := resp.Header.Get("Content-Type")
  	if contentType != "" && !strings.HasPrefix(contentType, "text/html") { // "" to allow for no header being sent
  		return nil
  	}
    fmt.Println("After getting page of url::", targetURL)
    tokens := html.NewTokenizer(resp.Body)
  	for {
  		tokenType := tokens.Next()
  		if tokenType == html.ErrorToken { //an EOF
  			return nil
  		}
  		token := tokens.Token()
  		if tokenType == html.StartTagToken { //opening tag
  			switch token.DataAtom.String() {
  			case "a", "link": //link tags
  				for _, attr := range token.Attr {
  					if attr.Key == "href" {
  						_, ok := seen_urls[attr.Val]
  						if !ok {
  							// seen_urls[attr.Val] = struct{}{} //add this ref to list of those seen on this page
                var child_url = attr.Val
                var domain_c = ""
                if strings.HasPrefix(attr.Val, "/") {
                  fmt.Println("From Relative URL::", child_url)
                  child_url = protocol+ "://"+domain+child_url
                  fmt.Println("Relative to absolute URL::", child_url)
                  domain_c = domain
                } else {
                  temp, err := url.Parse(child_url)
                	if err != nil {
                		fmt.Println("couldn't parse that URL:", child_url, err)
                		os.Exit(1)
                	}
                  domain_c = temp.Host
                }
                if (external_crawl || domain_c == domain) {
    							crawl_recursive(child_url, level-1, external_crawl, seen_urls)
                }
  						}
  					}
  				}
  			}
  		}
  	}
  }
  return nil
}
