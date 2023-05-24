package main

import (
	"fmt"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

func biqusoso() *THandler {
	return &THandler{
		S_fnc: func(c *colly.Collector, d *TSearch, name string) {
			/*
				<div class="search-list">
					<h2>
						搜索"噩梦惊袭"相关小说
					</h2>
					<ul>
						<li><span class="s1"><b>序号</b></span>
							<span class="s2"><b>作品名称</b></span>
							<span class="s4"><b>作者</b></span>
						</li>
						<li>
							<span class="s1">1</span>
							<span class="s2"><a href="http://www.qu-la.com/book/goto/id/94897903" target="_blank">噩梦惊袭</a></span>
							<span class="s4">温柔劝睡师</span>
						</li>
					</ul>
				</div>
			*/
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

			url := "https://www.166xs.org/search.php"
			if name != "" {
				url = url + "&keyword=" + name
			}
			c.OnError(func(_ *colly.Response, err error) {
				fmt.Println("search error ", err)
				d.Data = []TSearchData{}
			})
			c.Visit(url)
		},
		B_fnc: func(c *colly.Collector, d *TBook, url string) {
			/*
				<div class="book-chapter-list">
				<h3>最新章节预览</h3>
				<ul class="cf">
					<li><a href="/booktxt/82989775116/7979411116.html">460. 育英综合大学 温简言 “冤枉啊”……</a></li>
					.....
				</ul>
				<h3>正文</h3>
				<ul class="cf">
					<li><a href="/booktxt/82989775116/1920260116.html">第1章 德才中学</a></li>
				</ul>
			*/
			c.OnHTML("#listt", func(e *colly.HTMLElement) {
				e.ForEach("dt", func(i int, el *colly.HTMLElement) {
					if i > 0 {
						el.ForEach("dd a", func(li int, a *colly.HTMLElement) {
							link := a.Attr("href")
							d.Data = append(d.Data, TBookData{
								Name: a.Text,
								Link: link,
							})
						})
					}
				})
			})
			c.OnError(func(_ *colly.Response, err error) {
				fmt.Println("book error ", err)
				d.Data = []TBookData{}
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
			c.OnError(func(_ *colly.Response, err error) {
				fmt.Println("chapter error ", err)
				d.Data = TChapterData{}
			})

			c.Visit(url)
		},
		Host:    "www.166xs.org",
		Origin:  "www.166xs.org",
		Referer: "https://www.166xs.org",
		BURL: func(key string) string {
			return "https://www.166xs.org" + key
		},
		CURL: func(key string) string {
			return "https://www.166xs.org" + key
		},
	}
}
