package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	//"net/http"
	"os"
	"strconv"
	"time"
	//"strings"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

func parseTV(pagina int, ch chan string, url string, count *Count) {
	uri := url + strconv.Itoa(pagina)

	doc, err := goquery.NewDocument(uri)

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	fmt.Printf("imprimiendo %#v \n", uri)
	//value := 0
	doc.Find("#categoryProductContainer").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := s.Find("a.productName span").Text()
		shortDescription := s.Find("div.product10ShortDescription p").Text()

		perm := os.FileMode(0777)

		f, err := os.OpenFile("/home/amartinez/GO/files/"+strconv.Itoa(pagina)+".txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, perm)
		if err != nil {
			fmt.Printf("%#v", err)
		}
		defer f.Close()

		if _, err = f.WriteString(title + shortDescription + "\n"); err != nil {
			panic(err)
		}

		fmt.Printf("%s, %s\n", title, shortDescription)

	})

	ch <- "termina"

}

type Count struct {
	count int
}

// Extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool, numHref int) {
	count := Count{0}

	i := 0
	go func() {
		for i < 24 {
			i++
			fmt.Printf("aaa %#v", i)

			parseTV(i, ch, url, &count)
			//fmt.Println(count.count)
		}
	}()

	for count.count < 24 {
		select {
		case <-ch:
			count.count++
			fmt.Printf("canal %#v\n", count.count)

		}
	}
	chFinished <- true

	//fmt.Printf("Review %s\n", doc)

	// b := resp.Body

	// defer b.Close() // close Body when the function returns

	// z := html.NewTokenizer(b)

	// for {
	// 	tt := z.Next()

	// 	switch {
	// 	case tt == html.ErrorToken:
	// 		// End of the document, we're done
	// 		return
	// 	case tt == html.StartTagToken:
	// 		t := z.Token()
	// 		// Check if the token is an <a> tag
	// 		isAnchor := t.Data == "li.div"
	// 		if !isAnchor {
	// 			continue
	// 		}

	// 		fmt.Printf("%v", t)

	// 		// Extract the href value, if there is one
	// 		ok, url := getHref(t)
	// 		if !ok {
	// 			continue
	// 		}

	// 		// Make sure the url begines in http**
	// 		hasProto := strings.Index(url, "http") == 0
	// 		if hasProto {
	// 			numHref++
	// 			ch <- url
	// 			//go crawl(url, ch, chFinished, numHref)
	// 		}
	// 	}
	// }
}

func main() {
	start := time.Now()

	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]
	numHref := 0
	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished, numHref)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case <-chFinished:
			c++
		}
	}

	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s", elapsed)
	close(chUrls)
}
