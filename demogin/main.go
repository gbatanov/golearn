package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const VERSION = "0.0.2"
const SERVER = "192.168.76.95:8089"

type ActionHandler struct{}

func NewActionHandler() *ActionHandler {
	ts := ActionHandler{}
	return &ts
}

func main() {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob(".\\html\\*")

	actionSrv := NewActionHandler()

	router.GET("/ping", actionSrv.ping)
	router.GET("/cmd", actionSrv.cmdHandler)

	log.Printf("%v", router.Routes())
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
	log.Printf("%v", c.Query("id"))
	//JSON serializes the given struct as JSON into the response body.
	// It also sets the Content-Type as "application/json"
	c.JSON(http.StatusOK, gin.H{
		"message": c.Query("id"),
	})
}

// id - это не /cmd/?id=89689, а /cmd/hdjhg при роутере /cmd/:id
// для /cmd/?id=89689 использовать c.Query("id")
func (ts *ActionHandler) cmdHandler(c *gin.Context) {
	id := c.Query("id")   //c.Params.ByName("id")
	cmd := c.Query("cmd") //c.Params.ByName("cmd")

	log.Printf("%s %s", id, cmd)

	// HTML ответ на основе шаблона
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "GSB website", "id": id, "cmd": cmd})
}
