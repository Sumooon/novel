package main

import (
    "io"
    "fmt"
    "net/http"
    "strings"
    "encoding/json"
    colly "github.com/gocolly/colly/v2"
)

func main() {
    http.HandleFunc("/bqj", HttpHandler)
    // 启动http服务
    http.ListenAndServe(":9090", nil)
}

func getClient() *colly.Collector {
    // NewCollector(options ...func(*Collector)) *Collector
    // 声明初始化NewCollector对象时可以指定Agent，连接递归深度，URL过滤以及domain限制等
    c := colly.NewCollector(
        //colly.AllowedDomains("www.ibiquges.com"),
        colly.UserAgent("Opera/9.80 (Windows NT 6.1; U; zh-cn) Presto/2.9.168 Version/11.50"))
	
    // 发出请求时附的回调
    c.OnRequest(func(r *colly.Request) {
        // Request头部设定
        r.Headers.Set("Host", "www.ibiquges.com")
        r.Headers.Set("Connection", "keep-alive")
        r.Headers.Set("Accept", "*/*")
        r.Headers.Set("Origin", "https://www.ibiquges.com")
        r.Headers.Set("Referer", "https://www.ibiquges.com/modules/article/waps.php")
        r.Headers.Set("Accept-Encoding", "gzip, deflate")
        r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")

        fmt.Println("Visiting", r.URL)
    })

    // extract status code
    c.OnResponse(func(r *colly.Response) {
        fmt.Println("response received", r.StatusCode)
        // 设置context
        // fmt.Println(r.Ctx.Get("url"))
    })

    // 对visit的线程数做限制，visit可以同时运行多个
    c.Limit(&colly.LimitRule{
        Parallelism: 2,
        //Delay:      5 * time.Second,
    })
    return c
}

type TEmpty struct {
    Code    int     `json:"code"`
    Msg     string  `json:"msg"`
    Data    string  `json:"data"`
}
func HttpHandler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  if len(params) == 0 {
    res := new(TEmpty)
    res.Code = 0
    res.Msg = ""
    ret_json,_ := json.Marshal(res)
    io.WriteString(w, string(ret_json))
  } else {
    for key, value := range params {
    switch key {
        case "name":
            ret_json,_ := json.Marshal(searchBook(value[0]))
	        io.WriteString(w, string(ret_json))
        case "book":
            ret_json,_ := json.Marshal(getBook("https://www.ibiquges.com/" + value[0]))
	        io.WriteString(w, string(ret_json))
        case "chapter":
            ret_json,_ := json.Marshal(getChapter("https://www.ibiquges.com/" + value[0] + ".html"))
	        io.WriteString(w, string(ret_json))
        default:
            res := new(TEmpty)
            res.Code = 0
            res.Msg = ""
            ret_json,_ := json.Marshal(res)
            io.WriteString(w, string(ret_json))
    }
  }
  }
}

/*
    搜索书本/作者名称
*/
type TSearchData struct {
    Name    string `json:"name"`
    Author  string `json:"author"`
    Chapter  string `json:"chapter"`
    Link    string `json:"link"`
    Time    string `json:"time"`
}
type TSearch struct {
    Code    int             `json:"code"`
    Msg     string          `json:"msg"`
    Data    []TSearchData   `json:"data"`
}
func searchBook(name string) *TSearch {
    res := new(TSearch)
    res.Code = 0
    res.Msg = ""
    if name == "" {
        return res
    }
    c := getClient()

    c.OnHTML(".grid tr", func(e *colly.HTMLElement) {
        link := strings.Replace(e.ChildAttr("td a", "href"), "https://www.ibiquges.com/", "", -1)
        var data TSearchData
        data.Link = link
        e.ForEach("td", func(i int, el *colly.HTMLElement) {
            text := strings.Trim(strings.Replace(el.Text, "\n", "",  -1), " ")
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
            res.Data = append(res.Data, data)
        }
    })

	err := c.Post("https://www.ibiquges.com/modules/article/waps.php", map[string]string{"searchkey": name})
	if err != nil {
		fmt.Printf("异常 %s", err)
        res.Code = 10001
        res.Msg = "搜索接口异常"
	}
    return res
}

/*
    获取书目录
*/
type TBookData struct {
    Name    string `json:"name"`
    Link    string `json:"link"`
}
type TBook struct {
    Code    int             `json:"code"`
    Msg     string          `json:"msg"`
    Data    []TBookData   `json:"data"`
}

func getBook(url string) *TBook {
    res := new(TBook)
    res.Code = 0
    res.Msg = ""
    if url == "" {
        return res
    }
    c := getClient()

    c.OnHTML("body", func(e *colly.HTMLElement) {
        // <dd>
        // <a href="/79/79077/43685553.html">第1651章 玄聪</a>
        // </dd>
        e.ForEach("dd a", func(i int, el *colly.HTMLElement) {
            link := strings.Replace(el.Attr("href"), ".html", "", -1)
            title := el.Text
            res.Data = append(res.Data, TBookData{
                Name: title,
                Link: link,
            })
        })
    })

    c.Visit(url)
    return res
}

/*
    获取章节内容
*/
type TChapterData struct {
    Title       string  `json:"title"`
    Content     string  `json:"content"`
}
type TChapter struct {
    Code    int             `json:"code"`
    Msg     string          `json:"msg"`
    Data    TChapterData   `json:"data"`
}
func getChapter(url string) *TChapter {
    res := new(TChapter)
    res.Code = 0
    res.Msg = ""
    if url == "" {
        return res
    }
    c := getClient()

    c.OnHTML("div .bookname h1", func(e *colly.HTMLElement) {
        res.Data.Title = e.Text
    })
    c.OnHTML("#content", func(e *colly.HTMLElement) {
        // 去除网站广告
        content := strings.Split(e.Text, "亲,点击进去")[0]
        res.Data.Content = content
    })

    c.Visit(url)
    return res
}