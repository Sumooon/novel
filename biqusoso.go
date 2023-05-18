package main

import (
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
			c.OnHTML(".search-list", func(e *colly.HTMLElement) {
				e.ForEach("li", func(i int, el *colly.HTMLElement) {
					if i > 0 {
						var data TSearchData
						data.Link = strings.Replace(el.ChildAttr(".s2 a", "href"), "http://www.qu-la.com/book/goto/id/", "", -1)
						data.Name = el.ChildText(".s2")
						data.Author = el.ChildText(".s4")
						d.Data = append(d.Data, data)
					}
				})
			})

			url := "https://so.biqusoso.com/s.php?siteid=lvsetxt.com&ie=utf-8"
			if name != "" {
				url = url + "&q=" + name
			}
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
			c.OnHTML(".book-chapter-list", func(e *colly.HTMLElement) {
				e.ForEach("ul", func(i int, el *colly.HTMLElement) {
					if i > 0 {
						el.ForEach("li a", func(li int, a *colly.HTMLElement) {
							href := a.Attr("href")
							link := strings.Replace(strings.Replace(href, ".html", "", -1), "/booktxt/", "", -1)
							d.Data = append(d.Data, TBookData{
								Name: a.Text,
								Link: link,
							})
						})
					}
				})
			})

			c.Visit(url)
		},
		C_fnc: func(c *colly.Collector, d *TChapter, url string) {
			c.OnHTML("#chapter-title h1", func(e *colly.HTMLElement) {
				d.Data.Title = e.Text
			})
			c.OnHTML("#txt", func(e *colly.HTMLElement) {
				ret, _ := e.DOM.Html()
				// 去除网站广告
				content := strings.Split(ret, "『如果章节错误，点此举报』</a>")[1]
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
			c.OnHTML("#pb_prev", func(e *colly.HTMLElement) {
				href := e.Attr("href")
				d.Data.Prev = strings.Replace(strings.Replace(href, ".html", "", -1), "/booktxt/", "", -1)
			})
			c.OnHTML("#pb_mulu", func(e *colly.HTMLElement) {
				d.Data.Dir = strings.Replace(e.Attr("href"), "/booktxt/", "", -1)
			})
			c.OnHTML("#pb_next", func(e *colly.HTMLElement) {
				d.Data.Next = e.Attr("href")
				href := e.Attr("href")
				d.Data.Next = strings.Replace(strings.Replace(href, ".html", "", -1), "/booktxt/", "", -1)
			})

			c.Visit(url)
		},
		Host:    "so.biqusoso.com",
		Origin:  "so.biqusoso.com",
		Referer: "https://so.biqusoso.com",
		BURL: func(key string) string {
			return "http://www.qu-la.com/booktxt/" + key + "116/"
		},
		CURL: func(key string) string {
			return "http://www.qu-la.com/booktxt/" + key + ".html"
		},
	}
}
