package urlshorten

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"strings"
)

var validate = validator.New()

type App struct {
	Router      *mux.Router
	Middlewares *Middleware
	Config      *Env
	BaseUrl     string
}

type ShortReq struct {
	Url    string `json:"url" validate:"url"`
	Expire int    `json:"expire" validate:"min=0"`
}

type ShortResp struct {
	ShortLink string `json:"short_link"`
}

//响应json数据
func responseJson(w http.ResponseWriter, code int, message string, data interface{}) {
	m := make(map[string]interface{})
	m["code"] = code
	m["message"] = message
	m["data"] = data
	b, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

//响应404页面
func response404Page(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("<html><body><center>Page Not Found!</center></body></html>"))
}

//初始化
func (a *App) Init(baseUrl string, c *Env) {
	a.Config = c
	a.BaseUrl = baseUrl
	a.Middlewares = &Middleware{}

	//设置路由为mux路由（第三方路由库）
	a.Router = mux.NewRouter()
	a.initRoutes()
}

//初始化路由
func (a *App) initRoutes() {
	m := alice.New(a.Middlewares.LoggingHandler)
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortLink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortLinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirectShortLink)).Methods("GET")
}

//创建短链接
func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req ShortReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseJson(w, 1000, "Invalid Request", nil)
		return
	}
	defer r.Body.Close()

	log.Println("Req:", req)

	//验证POST数据内容
	if err := validate.Struct(req); err != nil {
		responseJson(w, 1000, "Invalid Request", err)
		return
	}

	//生成短链接
	eid, err := a.Config.S.Shorten(req.Url, req.Expire)
	if err != nil {
		responseJson(w, 1001, "System Error", err)
		return
	}

	//返回json数据
	shortUrl := a.BaseUrl + eid
	responseJson(w, 0, "OK", shortUrl)
}

//获取短链接详细信息
func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	link := vals.Get("link")

	log.Println("link:", link)
	//去除url前缀，获取eid
	eid := strings.TrimLeft(link, a.BaseUrl)
	m, err := a.Config.S.ShortenInfo(eid)
	if err != nil {
		responseJson(w, 1002, "Info Not Found", err)
		return
	}

	//返回json数据
	responseJson(w, 0, "OK", m)
}

func (a *App) redirectShortLink(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Path
	if len(link) <= 1 {
		response404Page(w)
		return
	}
	eid := link[1:]

	log.Println("eid:", eid)
	//根据短链接获取原始url
	url, err := a.Config.S.UnShorten(eid)
	if err != nil {
		response404Page(w)
		return
	}

	log.Println("url:", url)
	log.Println("err:", err)

	//必须先设置location，再调用writeHeader才能生效
	w.Header().Set("Location", url)
	w.WriteHeader(302)
}

func (a *App) Run(addr string, baseUrl string) {
	a.Init(baseUrl, GetEnv())
	http.ListenAndServe(addr, a.Router)
}
