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

package material

import (
	"encoding/json"
	"github.com/fastwego/wxwork/corporation"
	"github.com/fastwego/wxwork/corporation/apis/material"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var Corp *corporation.Corporation
var App *corporation.App

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	Corp = corporation.New(corporation.Config{Corpid: viper.GetString("CROPID")})
	App = Corp.NewApp(corporation.AppConfig{
		AgentId:        viper.GetString("AGENTID"),
		Secret:         viper.GetString("SECRET"),
		Token:          viper.GetString("TOKEN"),
		EncodingAESKey: viper.GetString("EncodingAESKey"),
	})
}

func Demo(c *gin.Context) {

	// 上传图片
	resp, err := material.UploadImg(App,"material/bilibili.png")
	c.Writer.Write(resp)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}

	// 上传文件
	params := url.Values{}
	params.Add("type","file")
	resp, err = material.Upload(App,"material/material.go",params)

	c.Writer.Write(resp)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}

	mediaResp := struct {
		Type      string `json:"type"`
		MediaID   string `json:"media_id"`
		CreatedAt string `json:"created_at"`
	}{}

	json.Unmarshal(resp, &mediaResp)

	// 下载素材
	params = url.Values{}
	params.Add("media_id", mediaResp.MediaID)
	header := http.Header{}
	header.Add("Range","bytes=0-127")
	response, err := material.Get(App, params, header)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(string(all))
	log.Println(response.Header)
}
