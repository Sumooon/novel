package main

import (
	"fmt"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

func biqusoso() *THandler {
	return &THandler{
		S_fnc: func(c *colly.Collector, d *TSearch, name string) {
			c.OnHTML(".novelslist2", func(e *colly.HTMLElement) {
				e.ForEach("li", func(i int, el *colly.HTMLElement) {
					if i > 0 {
						var data TSearchData
						data.Link = el.ChildAttr(".s2 a", "href")
						data.Name = el.ChildText(".s2")
						data.Author = el.ChildText(".s4")
						data.Chapter = el.ChildText(".s3")
						data.Time = el.ChildText(".s6")
						d.Data = append(d.Data, data)
					}
				})
			})
			c.OnError(func(_ *colly.Response, err error) {
				fmt.Println("search err:", err)
			})

			url := "https://www.166xs.org/search.php"
			if name != "" {
				url = url + "?keyword=" + encodeURI(name)
			}
			c.Visit(url)
		},
		B_fnc: func(c *colly.Collector, d *TBook, url string) {
			c.OnHTML("#list dl", func(e *colly.HTMLElement) {
				// ret, _ := e.DOM.Html()
				// str := strings.Split(ret, "正文</dt>")[1]
				// fmt.Println("res ", str)
				e.ForEach("dd", func(i int, el *colly.HTMLElement) {
					if i > 14 {
						link := el.ChildAttr("a", "href")
						d.Data = append(d.Data, TBookData{
							Name: el.Text,
							Link: link,
						})
					}
				})
			})

			c.Visit(url)
		},
		C_fnc: func(c *colly.Collector, d *TChapter, url string) {
			c.OnHTML(".bookname h1", func(e *colly.HTMLElement) {
				d.Data.Title = e.Text
			})
			c.OnHTML("#content", func(e *colly.HTMLElement) {
				ret, _ := e.DOM.Html()
				// 去除网站广告
				content := strings.Split(ret, "<div align=\"center\">")[0]
				d.Data.Content = content
			})
			/*
				<div class="chapter-control">
					<a class="url_pre" id="pb_prev" href="/booktxt/82989775116/7979410116.html">上一章</a>
					<span>|</span>
					<a id="pb_mulu" href="/booktxt/82989775116/">返回目录</a>
					<span>|</span>
					<a class="url_next" id="pb_next" href="/booktxt/82989775116/">下一章</a>
				</div>
			*/
			c.OnHTML(".pre", func(e *colly.HTMLElement) {
				d.Data.Prev = e.Attr("href")
			})
			c.OnHTML(".back", func(e *colly.HTMLElement) {
				d.Data.Dir = e.Attr("href")
			})
			c.OnHTML(".next", func(e *colly.HTMLElement) {
				d.Data.Next = e.Attr("href")
			})

			c.Visit(url)
		},
		Header: map[string]string{
			"Host":    "www.166xs.org",
			"Origin":  "www.166xs.org",
			"Referer": "https://www.166xs.org",
		},
		BURL: func(key string) string {
			return "https://www.166xs.org" + key
		},
		CURL: func(key string) string {
			return "https://www.166xs.org" + key
		},
	}
}
