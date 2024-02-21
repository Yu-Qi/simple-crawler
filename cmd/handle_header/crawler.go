package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func crawler() {
	c := colly.NewCollector()

	// get header from postman default setting
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
		r.Headers.Set("Referer", "https://www.google.com/")
		r.Headers.Set("Host", "https://www.google.com/")
		r.Headers.Set("Connect", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("visit success:", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("request failed:", r.Request.URL, err.Error())
	})

	c.Visit(fmt.Sprintf("http://localhost:%d/", port))
}
