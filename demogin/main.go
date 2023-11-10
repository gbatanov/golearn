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

type ActionServer struct{}

func NewActionServer() *ActionServer {
	ts := ActionServer{}
	return &ts
}

func main() {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("C:\\work\\demogin\\html\\*")

	actionSrv := NewActionServer()

	router.GET("/ping", actionSrv.ping)
	router.GET("/cmd", actionSrv.cmdHandler)

	log.Printf("%v", router.Routes())
	srv := &http.Server{
		Addr:    "192.168.76.95:8089",
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

func (ts *ActionServer) ping(c *gin.Context) {
	//	a, _ := c.GetRawData()
	a := c.Query("id")
	log.Printf("%v", a)
	//JSON serializes the given struct as JSON into the response body.
	// It also sets the Content-Type as "application/json"
	c.JSON(http.StatusOK, gin.H{
		"message": c.Query("id"),
	})
}

// id - это не /cmd/?id=89689, а /cmd/hdjhg при роутере /cmd/:id
// для /cmd/?id=89689 использовать c.Query("id")
func (ts *ActionServer) cmdHandler(c *gin.Context) {
	id := c.Query("id")   //c.Params.ByName("id")
	cmd := c.Query("cmd") //c.Params.ByName("cmd")

	log.Printf("%s %s", id, cmd)

	// HTML ответ на основе шаблона
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "GSB website", "id": id, "cmd": cmd})
}

/*
func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
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
*/
