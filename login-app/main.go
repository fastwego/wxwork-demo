// Copyright 2021 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastwego/wxwork/corporation/apis/message"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/fastwego/wxwork/corporation"
	"github.com/fastwego/wxwork/corporation/apis/contact/user"
	"github.com/fastwego/wxwork/corporation/apis/oauth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Corp *corporation.Corporation
var App *corporation.App

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	Corp = corporation.New(corporation.Config{
		Corpid: viper.GetString("CorpId"),
	})

	App = Corp.NewApp(corporation.AppConfig{
		AgentId: viper.GetString("AgentId"),
		Secret:  viper.GetString("Secret"),
	})

	fmt.Println(Corp, App)
}

func main() {

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("gosession", store))

	router.GET("/", Index)
	router.GET("/login", Login)

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

type User struct {
	Userid  string `json:"userid"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Message string `json:"message"`
}

func Index(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("user")

	loginUser, ok := user.(User)
	if !ok {
		loginUser = User{}
	}

	join := c.Query("join")
	if len(join) > 0 {
		// 发送 报名信息
		type Textcard struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Btntxt      string `json:"btntxt"`
		}

		textCard := struct {
			Touser                 string   `json:"touser"`
			Toparty                string   `json:"toparty"`
			Totag                  string   `json:"totag"`
			Msgtype                string   `json:"msgtype"`
			Agentid                string   `json:"agentid"`
			Textcard               Textcard `json:"textcard"`
			EnableIDTrans          int      `json:"enable_id_trans"`
			EnableDuplicateCheck   int      `json:"enable_duplicate_check"`
			DuplicateCheckInterval int      `json:"duplicate_check_interval"`
		}{
			Touser:  loginUser.Userid,
			Msgtype: "textcard",
			Agentid: App.Config.AgentId,
		}

		textCard.Textcard = Textcard{Title: "报名成功", Description: "请来我办公室 <br> <div class='highlight'>记得带上吃饭的家伙 ~</div>", URL: viper.GetString("ServerUrl"), Btntxt: "好的"}

		payload, err := json.Marshal(textCard)
		fmt.Println(string(payload), err)
		if err != nil {
			return
		}

		resp, err := message.Send(App, payload)
		fmt.Println(string(resp), err)

		loginUser.Message = "报名成功~"
	}

	t1, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	t1.Execute(c.Writer, loginUser)
}

func Login(c *gin.Context) {

	code := c.Query("code")

	// 跳转登录
	if len(code) == 0 {
		var redirectUri = viper.GetString("ServerUrl") + "/login"
		link := oauth.GetAuthorizeUrl(viper.GetString("CorpId"), redirectUri, "STATE")
		c.Redirect(302, link)
		return
	}

	// 获取用户身份
	accessToken, err := App.AccessToken.GetAccessTokenHandler(App)
	userInfo, err := oauth.GetUserInfo(accessToken, code)
	fmt.Println(userInfo, err)
	if err != nil {
		return
	}

	// 获取员工详细信息
	params := url.Values{}
	params.Add("userid", userInfo.UserID)
	resp, err := user.Get(App, params)
	fmt.Println(string(resp), err)
	if err != nil {
		return
	}

	user := User{}

	err = json.Unmarshal(resp, &user)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 记录 Session
	gob.Register(User{})
	session := sessions.Default(c)
	session.Set("user", user)
	fmt.Println(user)
	err = session.Save()

	if err != nil {
		fmt.Println(err)
		return
	}

	// 返回首页
	c.Header("Content-Type", "text/html")
	_, _ = c.Writer.WriteString(`<html"><head><meta http-equiv="refresh" content="0;URL='/'" /></head><body></body></html>`)
}
