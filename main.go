package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	colly "github.com/gocolly/colly/v2"
)

type TSearchFunc func(*colly.Collector, *TSearch, string)
type TBookFunc func(*colly.Collector, *TBook, string)
type TChapterFunc func(*colly.Collector, *TChapter, string)
type TBookUrl func(string) string
type TChapterUrl func(string) string

type THandler struct {
	S_fnc   TSearchFunc
	B_fnc   TBookFunc
	C_fnc   TChapterFunc
	BURL    TBookUrl
	CURL    TChapterUrl
	Host    string
	Origin  string
	Referer string
}

func main() {
	biquge := biquge()
	biqusoso := biqusoso()
	// 97xiaoshuo()
	http.HandleFunc("/bqj", func(w http.ResponseWriter, r *http.Request) {
		HttpHandler(w, r, biquge)
	})
	http.HandleFunc("/bqss", func(w http.ResponseWriter, r *http.Request) {
		HttpHandler(w, r, biqusoso)
	})
	// 启动http服务
	http.ListenAndServe(":9090", nil)
}

func getClient(host string, origin string, referer string) *colly.Collector {
	// NewCollector(options ...func(*Collector)) *Collector
	// 声明初始化NewCollector对象时可以指定Agent，连接递归深度，URL过滤以及domain限制等
	c := colly.NewCollector(
		//colly.AllowedDomains("www.ibiquges.com"),
		colly.UserAgent("Opera/9.80 (Windows NT 6.1; U; zh-cn) Presto/2.9.168 Version/11.50"))

	// 发出请求时附的回调
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Host", host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", origin)
		r.Headers.Set("Referer", referer)
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
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []any  `json:"data"`
}

func HttpHandler(w http.ResponseWriter, r *http.Request, d *THandler) {
	params := r.URL.Query()
	if len(params) == 0 {
		res := new(TEmpty)
		res.Code = 0
		res.Msg = ""
		res.Data = []any{}
		ret_json, _ := json.Marshal(res)
		io.WriteString(w, string(ret_json))
	} else {
		for key, value := range params {
			switch key {
			case "name":
				ret_json, _ := json.Marshal(searchBook(value[0], d))
				io.WriteString(w, string(ret_json))
			case "book":
				ret_json, _ := json.Marshal(getBook(d.BURL(value[0]), d))
				io.WriteString(w, string(ret_json))
			case "chapter":
				ret_json, _ := json.Marshal(getChapter(d.CURL(value[0]), d))
				io.WriteString(w, string(ret_json))
			default:
				res := new(TEmpty)
				res.Code = 0
				res.Msg = ""
				res.Data = []any{}
				ret_json, _ := json.Marshal(res)
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
	Chapter string `json:"chapter"`
	Link    string `json:"link"`
	Time    string `json:"time"`
}
type TSearch struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []TSearchData `json:"data"`
}

func searchBook(name string, d *THandler) *TSearch {
	res := new(TSearch)
	res.Code = 0
	res.Msg = ""
	if name == "" {
		return res
	}
	c := getClient(d.Host, d.Origin, d.Referer)

	d.S_fnc(c, res, name)
	return res
}

/*
获取书目录
*/
type TBookData struct {
	Name string `json:"name"`
	Link string `json:"link"`
}
type TBook struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data []TBookData `json:"data"`
}

func getBook(url string, d *THandler) *TBook {
	res := new(TBook)
	res.Code = 0
	res.Msg = ""
	if url == "" {
		return res
	}
	c := getClient(d.Host, d.Origin, d.Referer)
	d.B_fnc(c, res, url)

	return res
}

/*
获取章节内容
*/
type TChapterData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Prev    string `json:"prev"`
	Dir     string `json:"dir"`
	Next    string `json:"next"`
}
type TChapter struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data TChapterData `json:"data"`
}

func getChapter(url string, d *THandler) *TChapter {
	res := new(TChapter)
	res.Code = 0
	res.Msg = ""
	if url == "" {
		return res
	}
	c := getClient(d.Host, d.Origin, d.Referer)
	d.C_fnc(c, res, url)

	return res
}
