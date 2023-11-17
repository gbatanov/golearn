package api

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"git.makves.ru/test/tlgserv/tlg"
)

const BOT_SERVER_ADDRESS = "192.168.76.95:8055"
const BOT_NAME = "Makves_test_bot"
const MY_ID = int64(836487770)

type TlgBlock struct {
	tlg32      *tlg.Tlg32
	tlgMsgChan chan tlg.Message
}

type HttpApi struct {
	MsgIn    string
	MsgOut   string
	Address  string
	srv      *http.Server
	Quit     chan string
	tlgBlock TlgBlock
	wg       sync.WaitGroup
}

func ApiCreate() (*HttpApi, error) {
	q := make(chan string)
	api := HttpApi{
		MsgIn:  "",
		MsgOut: "",
		Quit:   q,
		srv:    &http.Server{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/", api.apiHandler)
	api.srv.Addr = BOT_SERVER_ADDRESS
	api.srv.Handler = mux

	api.tlgBlock.tlgMsgChan = make(chan tlg.Message, 16)
	api.tlgBlock.tlg32 = tlg.Tlg32Create(BOT_NAME, "prod", MY_ID, api.tlgBlock.tlgMsgChan, &api.wg)

	return &api, nil
}

func (api *HttpApi) Start() error {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		api.Stop()
	}()

	// Старт телеграм бота
	err := api.tlgBlock.tlg32.Run()

	if err != nil {
		log.Println(err.Error())
		return err
	}
	// Старт HTTP-сервера
	go func() {
		if err := api.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}() // listen and serve

	<-api.Quit
	api.tlgBlock.tlg32.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = api.srv.Shutdown(ctx)
	if err == nil {
		<-ctx.Done()
	}
	return nil
}

func (api *HttpApi) Stop() {
	api.Quit <- "stop"
}

// минимальное требование - параметр msg
// если параметр id отсутствует, берется дефолтный id разработчика
func (api *HttpApi) apiHandler(w http.ResponseWriter, req *http.Request) {
	uri := req.RequestURI
	u, err := url.Parse(uri)
	if err == nil {
		m, _ := url.ParseQuery(u.RawQuery)
		_, ok := m["msg"]
		if ok {
			api.MsgIn = m["msg"][0]
			if len(api.MsgIn) > 0 {
				_, ok = m["id"]
				chatId := MY_ID
				if ok {
					id, err := strconv.ParseInt(m["id"][0], 10, 64)
					if err == nil {
						chatId = id
					}
				}
				api.tlgBlock.tlg32.MsgChan <- tlg.Message{ChatId: chatId, Msg: api.MsgIn}
				api.MsgOut = "Ok"
				api.sendAnswer(w, 200)
				return
			}
		}
	}
	api.MsgOut = "Error"
	api.sendAnswer(w, 404)
}

// Отправка заголовка клиенту
func (api *HttpApi) sendHeaders(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	w.WriteHeader(code)
}

func (api *HttpApi) sendAnswer(w http.ResponseWriter, code int) {
	api.sendHeaders(w, code)
	w.Write([]byte(api.MsgOut))
	api.MsgOut = ""
}
