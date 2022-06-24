package main

import (
	"fmt"
	"bytes"
	"time"
    "io/ioutil"
    "net/http"
    "regexp"

    "github.com/PuerkitoBio/goquery"
    "github.com/saintfish/chardet"
    "golang.org/x/net/html/charset"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
	IsFetched(url string) bool
}

func fetch(url string) [10]string {
    var validUrl = regexp.MustCompile(`^/wiki/`)
    // Getリクエスト
    res, _ := http.Get(url)
    defer res.Body.Close()

    // 読み取り
    buf, _ := ioutil.ReadAll(res.Body)

    // 文字コード判定
    det := chardet.NewTextDetector()
    detRslt, _ := det.DetectBest(buf)
    // fmt.Println(detRslt.Charset)
    // => EUC-JP

    // 文字コード変換
    bReader := bytes.NewReader(buf)
    reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)

    // HTMLパース
    doc, _ := goquery.NewDocumentFromReader(reader)
    // titleを抜き出し
	title := doc.Find("title").Text()
	fmt.Println(title)

	var urls [10]string
	var i int = 0
    doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		if i > 9 {
			return
		}
        url, _ := s.Attr("href")
        matched := validUrl.MatchString(url)
        if matched {
			urls[i] = "https://ja.wikipedia.org"+ url
			i++
        }
  })
  	return urls
}



// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int) {
	if depth <= 0 {
		return
	}
	urls := fetch(url)

	for _, u := range urls {
		func(u string) {
			Crawl(u, depth-1)
		}(u)
	}
}

func main() {
	start := time.Now()
	Crawl("https://ja.wikipedia.org/wiki/SCHOOL_OF_LOCK!", 4)
	end := time.Now()
	fmt.Printf("%f秒\n", (end.Sub(start)).Seconds())
}
