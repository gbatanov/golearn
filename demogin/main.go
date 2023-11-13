package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const VERSION = "0.0.7"
const SERVER = "192.168.76.95:8089"

// Пример встраивания целой папки в переменную
//
//go:embed  html/*
var f embed.FS

// Пример встраивания текстового файла в строку
//
//go:embed README.md
var readme string

//go:embed html/img/fav32.png
var icon []byte

type ActionHandler struct{}

func NewActionHandler() *ActionHandler {
	ts := ActionHandler{}
	return &ts
}

func main() {
	//	log.Println(readme)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	gin.SetMode(gin.ReleaseMode)
	//	router := gin.Default()
	router := gin.New()

	// Используем встроенные шаблоны (включаются в тело программы)
	templ := template.Must(template.New("").ParseFS(f, "html/tpl/*.tmpl"))
	router.SetHTMLTemplate(templ)
	// Использования шаблонов из ФС
	//	  router.LoadHTMLGlob(".\\html\\tpl\\*")
	//
	//router.StaticFS("/css", http.FS(f)) // Не заработало /html /html/css /css
	router.StaticFS("/css", http.Dir("html/css")) // так работает, но это не embed
	//router.Static("/css", ".\\html\\css") // так работает

	// картинка
	router.StaticFile("/img/fav32.png", ".\\html\\img\\fav32.png")
	// иконка из файла
	//	router.StaticFile("/favicon.ico", ".\\html\\img\\fav32.png")
	// иконка из внедренной картинки
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Data(
			http.StatusOK,
			"image/x-icon",
			icon,
		)
	})

	actionSrv := NewActionHandler()

	router.GET("/ping", actionSrv.ping)
	router.GET("/cmd", actionSrv.cmdHandler)
	router.GET("/", actionSrv.mainPage)

	//	log.Printf("%v", router.Routes())
	srv := &http.Server{
		Addr:    SERVER,
		Handler: router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}() // listen and serve
	<-quit
	log.Println("Shutdown Server ...")
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
	log.Println("Server exiting")
}

func (ts *ActionHandler) ping(c *gin.Context) {
	log.Printf("%s", "ping redirect to main page")
	c.Redirect(http.StatusPermanentRedirect, "/")
}

// id - это не /cmd/?id=89689, а /cmd/hdjhg при роутере /cmd/:id
// для /cmd/?id=89689 использовать c.Query("id")
func (ts *ActionHandler) cmdHandler(c *gin.Context) {
	id := c.Query("id")   //c.Params.ByName("id")
	cmd := c.Query("cmd") //c.Params.ByName("cmd")

	log.Printf("%s %s", id, cmd)

	// HTML ответ на основе шаблона
	c.HTML(http.StatusOK, "command.tmpl", gin.H{"title": "GSB website", "id": id, "cmd": cmd})
}
func (ts *ActionHandler) mainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "GSB website"})
}
