package main

import (
	aborter "context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	"cube/config"
	"cube/core"
	"cube/core/helpers/action"
	"cube/lib/context"
	"cube/lib/database"
	"cube/lib/logger"
	"cube/lib/utils"
	"cube/middleware"
	"cube/service/schedule"

	"github.com/google/shlex"

	"github.com/gin-gonic/gin"
)

var log = logger.Log()

func rootHandler(c *core.Core, ctx *gin.Context) {
	var err error
	var h gin.H
	var req context.ChatContext

	defer func(h *gin.H) {
		if err != nil {
			ctx.JSON(200, gin.H{
				"text": err.Error(),
			})
		} else {
			ctx.JSON(200, h)
		}
	}(&h)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return
	}

	bodyStr := string(body)

	if strings.HasPrefix(bodyStr, "payload=") {
		h = handleCallback(c, bodyStr)
		return
	}

	req, err = utils.ParseRequest(bodyStr)
	if err != nil {
		log.Error(err)
		return
	}

	r, err := handleCommand(req, c)
	if err != nil {
		return
	}

	h = r.Normalize()
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

func handleCallback(c *core.Core, body string) (h gin.H) {
	callback, err := utils.ParseActionCallback(body)
	if err != nil {
		log.Error(err)
		return
	}

	actionRecord := database.Action{}
	tx := c.DB.First(&actionRecord, map[string]interface{}{
		"action_id": callback.CallbackID,
	})
	if tx.Error != nil {
		h = context.NewErrorResponse(tx.Error).Normalize()
		return
	}

	if len(callback.Actions[0].Value) < 1 {
		h = context.NewTextResponse(fmt.Sprintf("closed by %v", callback.User.Username)).Normalize()
		return
	}

	req := context.ChatContext{
		Text:     callback.Actions[0].Value,
		UserID:   callback.User.UserID,
		Username: callback.User.Username,
		PostID:   callback.CallbackID,
		Token:    callback.Token,
	}

	r, err := handleCommand(req, c)
	if err != nil {
		return
	}

	callbackID, err := strconv.Atoi(callback.CallbackID)
	if err != nil {
		h = context.NewErrorResponse(err).Normalize()
		return
	}

	resp, err := action.NewActionResponse(c.DB, uint64(callbackID), req.UserID)
	if err != nil {
		h = context.NewErrorResponse(err).Normalize()
		return
	}

	resp.Title = r.Text()
	resp.Message = fmt.Sprintf("%v (%v)", resp.Message, callback.Actions[0].Text)
	h = resp.Normalize()
	return
}

func handleCommand(req context.ChatContext, c *core.Core) (r context.IResponse, err error) {
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
		r = core.PrintHelp()
		return
	}

	r = c.Handle(req, args)
	return
}
