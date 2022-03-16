package main

import (
	aborter "context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"

	"cube/config"
	"cube/core"
	"cube/lib/database"
	"cube/lib/logger"
	"cube/lib/utils"
	"cube/middleware"
	"cube/service/schedule"

	"github.com/google/shlex"

	"github.com/gin-gonic/gin"
)

var log = logger.Log

func rootHandler(c *core.Core, ctx *gin.Context) {
	var err error
	var text string
	defer func(text *string) {
		resp := gin.H{
			"text": text,
		}

		if err != nil {
			resp["text"] = err.Error()
		}

		ctx.JSON(200, resp)
	}(&text)

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return
	}

	req, err := utils.ParseRequest(body)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(req.Text)

	var cmd string
	if strings.HasPrefix(req.Text, "/fp") {
		cmd = "food poll "
		if len(req.Text) > 3 {
			cmd += req.Text[3:]
		}
	} else {
		cmd = req.Text
	}

	args, err := shlex.Split(cmd)
	if err != nil {
		return
	}
	if len(args) < 1 {
		text = string(core.PrintHelp())
		return
	}

	text = string(c.Handle(req, args))
}

func main() {
	err := config.Init("etc/conf.yaml")
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	if len(config.Conf.Token) < 1 {
		log.Error("bot token is empty")
		return
	}

	db, err := database.New("chat.db")
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	scheduleService := schedule.NewService(db)
	c := core.New(
		db,
		scheduleService,
	)

	r := gin.New()
	r.Use(middleware.LogRequest(), gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.POST("/chat", func(ctx *gin.Context) {
		rootHandler(c, ctx)
	})

	r.GET("/foodmap", c.Food.RenderMap)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	var wg sync.WaitGroup
	abortCtx, abortFunc := aborter.WithCancel(aborter.Background())
	defer abortFunc()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	wg.Add(1)
	go func() {
		scheduleService.Run(abortCtx)
		wg.Done()
	}()

	// wait schedule service starts
	err = c.InitAlarms()
	if err != nil {
		log.Fatal("Core.InitAlarms: ", err)
	}

	wg.Add(1)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
		wg.Done()
	}()

	<-quit
	abortFunc()
	if err := srv.Shutdown(abortCtx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	wg.Wait()
}
