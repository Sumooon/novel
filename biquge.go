package main

import (
	"fmt"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

func biquge() *THandler {
	return &THandler{
		S_fnc: func(c *colly.Collector, d *TSearch, name string) {
			c.OnHTML(".grid tr", func(e *colly.HTMLElement) {
				link := strings.Replace(e.ChildAttr("td a", "href"), "https://www.ibiquges.com/", "", -1)
				var data TSearchData
				data.Link = link
				e.ForEach("td", func(i int, el *colly.HTMLElement) {
					text := strings.Trim(strings.Replace(el.Text, "\n", "", -1), " ")
					switch i {
					case 0:
						data.Name = text
					case 1:
						data.Chapter = text
					case 2:
						data.Author = text
					case 3:
						data.Time = text
						// default:
						//     ...
					}
				})
				if link != "" {
					d.Data = append(d.Data, data)
				}
			})

			err := c.Post("https://www.ibiquges.com/modules/article/waps.php", map[string]string{"searchkey": name})
			if err != nil {
				fmt.Printf("异常 %s", err)
				d.Code = 10001
				d.Msg = "搜索接口异常"
			}
		},
		B_fnc: func(c *colly.Collector, d *TBook, url string) {
			c.OnHTML("body", func(e *colly.HTMLElement) {
				// <dd>
				// <a href="/79/79077/43685553.html">第1651章 玄聪</a>
				// </dd>
				e.ForEach("dd a", func(i int, el *colly.HTMLElement) {
					link := strings.Replace(el.Attr("href"), ".html", "", -1)
					title := el.Text
					d.Data = append(d.Data, TBookData{
						Name: title,
						Link: link,
					})
				})
			})

			c.Visit(url)
		},
		C_fnc: func(c *colly.Collector, d *TChapter, url string) {
			c.OnHTML("div .bookname h1", func(e *colly.HTMLElement) {
				d.Data.Title = e.Text
			})
			c.OnHTML("#content", func(e *colly.HTMLElement) {
				// 去除网站广告
				content := strings.Split(e.Text, "亲,点击进去")[0]
				d.Data.Content = content
			})
			c.OnHTML(".bottem2", func(e *colly.HTMLElement) {
				e.ForEach("a", func(i int, el *colly.HTMLElement) {
					link := strings.Replace(el.Attr("href"), ".html", "", -1)
					url := strings.Replace(link, "https://www.ibiquges.com", "", -1)
					if !strings.Contains(url, "javascript:;") {
						switch i {
						case 1:
							{
								d.Data.Prev = url
							}
						case 2:
							{
								d.Data.Dir = url
							}
						case 3:
							{
								d.Data.Next = url
							}
						}
					}
				})
			})

			c.Visit(url)
		},
		Host:    "www.ibiquges.com",
		Origin:  "www.ibiquges.com",
		Referer: "https://www.ibiquges.com/modules/article/waps.php",
		BURL: func(key string) string {
			return "https://www.ibiquges.com/" + key
		},
		CURL: func(key string) string {
			return "https://www.ibiquges.com/" + key + ".html"
		},
	}
}
