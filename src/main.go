package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"os"
	"time"
)

func main() {
	// シードURL
	seedURL := "https://github.com/"
	// クローラーの深さ
	depth := 15
	// クロールしたURLを格納するスライス
	crawledURLs := make(map[string]bool)

	crawl(seedURL, depth, crawledURLs)

	// URLをファイルに書き出す
	now := time.Now()
	stringNow := now.Format("20060102150405")
	file, err := os.Create("../results/crawled" + stringNow + ".txt")
	if err != nil {
			fmt.Println("Error:", err)
			return
	}
	defer file.Close()

	for url := range crawledURLs {
			fmt.Fprintln(file, url)
	}
}

func crawl(url string, depth int, crawledURLs map[string]bool) {
	// 指定された深さまでクロールする
	if depth <= 0 {
		return
	}

	// すでにクロールしたURLの場合はスキップする
	if _, ok := crawledURLs[url]; ok {
		return
	}

	// URLをクロールする
	fmt.Println("Crawling:", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// HTMLを解析する
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// ページ内のリンクを抽出する
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && strings.HasPrefix(href, "http") {
			// クロールするURLを追加する
			crawl(href, depth-1, crawledURLs)
		}
	})

	// クロールしたURLをマップに追加する
	crawledURLs[url] = true
}
