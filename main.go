package main

import (
	"context"
	"github.com/fastwego/wechat4work-demo/calendar"
	"github.com/fastwego/wechat4work-demo/contact"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fastwego/wechat4work/corporation"

	"github.com/fastwego/wechat4work/corporation/type/type_message"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

var Corp *corporation.Corporation
var ContactApp *corporation.App

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	Corp = corporation.New(corporation.CorporationConfig{Corpid: viper.GetString("CROPID")})
	ContactApp = Corp.NewApp(corporation.AppConfig{
		AgentId:        viper.GetString("AGENTID"),
		Secret:         viper.GetString("SECRET"),
		Token:          viper.GetString("TOKEN"),
		EncodingAESKey: viper.GetString("EncodingAESKey"),
	})
}

func HandleMessage(c *gin.Context) {

	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))

	message, err := ContactApp.Server.ParseXML(body)
	if err != nil {
		log.Println(err)
	}

	var output interface{}
	switch message.(type) {
	case type_message.MessageText: // 文本 消息
		msg := message.(type_message.MessageText)

		// 回复文本消息
		output = type_message.ReplyMessageText{
			ReplyMessage: type_message.ReplyMessage{
				ToUserName:   type_message.CDATA(msg.FromUserName),
				FromUserName: type_message.CDATA(msg.ToUserName),
				CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
				MsgType:      type_message.ReplyMsgTypeText,
			},
			Content: type_message.CDATA(msg.Content),
		}
	}

	ContactApp.Server.Response(c.Writer, c.Request, output)
}

func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/api/weixin/contact", func(c *gin.Context) {
		ContactApp.Server.EchoStr(c.Writer, c.Request)
	})
	router.POST("/api/weixin/contact", HandleMessage)

	router.GET("/api/weixin/demo", contact.Demo)
	router.GET("/api/weixin/calendar", calendar.Demo)

	svr := &http.Server{
		Addr:    viper.GetString("LISTEN"),
		Handler: router,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
