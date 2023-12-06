package main

import (
	"context"
	"fmt"
	"hello_server/ws"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const Version = "0.1.10"

type MyResponse struct {
	head string
	body string
}

func main() {

	fmt.Println("\nHTTP Server start")
	sigs := make(chan os.Signal, 1)
	// признак прерывания по Ctrl+C
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// Создаем переменную для сервера
	var srv http.Server

	// Завершение по сигналам прерывания
	go func() {
		sig := <-sigs
		fmt.Println("\nSignal:", sig)
		srv.Shutdown(context.Background())

	}()
	// Создаем маршрутизатор
	mux := http.NewServeMux()
	// Наполняем его обрабатываемыми маршрутами
	register_routing(mux)
	srv.Addr = "192.168.88.240:8180"
	srv.Handler = mux
	srv.ListenAndServe()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	fmt.Println("\nHTTP Server finished")
}

func register_routing(mux *http.ServeMux) {
	// Фундаментальная концепция серверов net/http - это обработчики.
	// Обработчик - это объект, реализующий интерфейс http.Handler.

	// Мы регистрируем наши обработчики на сервере,
	// используя удобную функцию  HandleFunc.
	// Она устанавливает маршрут по умолчанию в пакете net/http и принимает функцию в качестве аргумента.

	// Обработчики на функциях
	mux.HandleFunc("/api/", apiHandler)
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/headers", show_headers)
	mux.HandleFunc("/ws/", ws.Ws_handler)
	mux.HandleFunc("/", other_handler)
}

// Пример обработчика с формированием корректной ссылки из заданных параметров
func apiHandler(w http.ResponseWriter, r *http.Request) {

	host := r.Host
	var protocol string
	_, proto_redirect := r.Header["X-Forwarded-Proto"]
	if proto_redirect && r.Header["X-Forwarded-Proto"][0] == "https" {
		protocol = "https://"
	} else {
		protocol = "http://"
	}

	var my_resp MyResponse

	my_resp.head = "<title>API Page</title>"
	my_resp.body = "<div>Это АПИ</div>"
	// Базовый url
	baseUrl, _ := url.Parse(protocol + host)

	// добавляем путь (автоматически будет escaped и добавлен стартовый слеш)
	baseUrl.Path += "hello"

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("tyt", "Hello, Георгий!")
	params.Add("a", "12")
	params.Add("a", "14")
	params.Add("tail", "")
	params.Add("gg", "2,67,iuu")

	// добавляем параметры к Uri с кодированием
	baseUrl.RawQuery = params.Encode()

	my_resp.body += fmt.Sprintf("<a href=\"%s\">Переход</a>", baseUrl.String())

	send_answer(w, my_resp, 200)

}

// Тут я переопределяю поведение при ошибке 404
func NotFound(w http.ResponseWriter, r *http.Request) {

	host := r.Host
	var my_resp MyResponse
	var protocol string
	_, proto_redirect := r.Header["X-Forwarded-Proto"]
	if proto_redirect && r.Header["X-Forwarded-Proto"][0] == "https" {
		protocol = "https://"
	} else {
		protocol = "http://"
	}

	baseUrl, _ := url.Parse(protocol + host)
	my_resp.body = "<div>Wrong URL</div>"
	my_resp.body += fmt.Sprintf("<a href=\"%s\">Home page</a>", baseUrl.String())

	my_resp.head = "<title>Page not found</title>"
	send_answer(w, my_resp, 404)
}

// Обработчики
// принимают в качестве аргументов http.ResponseWriter и http.Request.
// ResponseWriter используется для наполнения HTTP-ответа.

func hello(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	var my_resp MyResponse
	var protocol string
	_, proto_redirect := r.Header["X-Forwarded-Proto"]
	if proto_redirect && r.Header["X-Forwarded-Proto"][0] == "https" {
		protocol = "https://"
	} else {
		protocol = "http://"
	}

	baseUrl, _ := url.Parse(protocol + host)
	a := ""
	args, _ := get_params(r) // args это map[string][]string
	// gg=2,67,iuu аргумент воспринимается как одна строка, создается срез из одного элемента gg:[2,67,iuu]
	// a=12&a=14 - создается срез a:[12 14]
	// tal - создается пустой срез tal:[]
	if len(args) > 0 {

		for k, v := range args {
			if len(v) > 1 {
				a += k + "=["
				for _, v1 := range v {
					a += (v1 + ",")
				}
				a += "]"
			} else if len(v) == 1 && len(v[0]) > 0 {
				a += k + "=" + v[0]
			} else {
				a += k
			}
			a += " | "
		}
	}
	b := fmt.Sprint(args)
	my_resp.body = "<div><h4>Hello, Георгий!!!</h4><p>"
	my_resp.body += b + " <br>" + a
	my_resp.body += "</p></div>"
	my_resp.head = "<title>Hello</title>"
	my_resp.body += fmt.Sprintf("<a href=\"%s\">Home page</a>", baseUrl.String())

	send_answer(w, my_resp, 200)

}
func other_handler(w http.ResponseWriter, r *http.Request) {
	var my_resp MyResponse
	host := r.Host

	var protocol string
	_, proto_redirect := r.Header["X-Forwarded-Proto"]
	if proto_redirect && r.Header["X-Forwarded-Proto"][0] == "https" {
		protocol = "https://"
	} else {
		protocol = "http://"
	}

	baseUrl, _ := url.Parse(protocol + host)
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path != "/" {
		NotFound(w, r)
		return
	}
	my_resp.body = "<div>Main Page normal</div>"
	my_resp.body += fmt.Sprintf("<a href=\"%s\">Home page</a>", baseUrl.String())
	my_resp.body += fmt.Sprintf("<br><a href=\"%s\">НКБ НИР</a>", baseUrl.String()+"/horizon-cmm11/")

	my_resp.head = "<title>Main page</title>"
	send_answer(w, my_resp, 200)

}
func show_headers(w http.ResponseWriter, req *http.Request) {

	// Этот обработчик делает что-то более сложное,
	// читая все заголовки HTTP-запроса и вставляя их в тело ответа.
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

// Локальные функции
// Отправка ответа клиенту
func send_answer(w http.ResponseWriter, my_resp MyResponse, code int) {

	var result string = "<html>"

	result += "<head>"
	result += my_resp.head
	result += "</head><body>"
	result += my_resp.body
	result += "</body></html>"
	send_headers(w, code)
	w.Write([]byte(result))
}

// Извлечение параметров запроса из Uri
func get_params(req *http.Request) (url.Values, error) {
	uri := req.RequestURI
	u, err := url.Parse(uri)
	if err != nil {
		return url.Values{}, err
	}
	m, _ := url.ParseQuery(u.RawQuery)
	return m, nil
}

// Отправка заголовка клиенту
func send_headers(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(code)
}
