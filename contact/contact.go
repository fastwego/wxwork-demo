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

package contact

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/fastwego/wxwork/corporation"
	"github.com/fastwego/wxwork/corporation/apis/contact/department"
	"github.com/fastwego/wxwork/corporation/apis/contact/user"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Corp *corporation.Corporation
var ContactApp *corporation.App

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	Corp = corporation.New(corporation.Config{Corpid: viper.GetString("CROPID")})
	ContactApp = Corp.NewApp(corporation.AppConfig{
		AgentId:        viper.GetString("AGENTID"),
		Secret:         viper.GetString("SECRET"),
		Token:          viper.GetString("TOKEN"),
		EncodingAESKey: viper.GetString("EncodingAESKey"),
	})
}

func Demo(c *gin.Context) {

	params := url.Values{}
	resp, err := department.List(ContactApp, params)

	c.Writer.Write(resp)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}

	dept := struct {
		Department []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Parentid int    `json:"parentid"`
			Order    int    `json:"order"`
		} `json:"department"`
	}{}

	json.Unmarshal(resp, &dept)

	params = url.Values{}
	params.Add("department_id", strconv.Itoa(dept.Department[0].ID))
	resp, err = user.SimpleList(ContactApp, params)
	c.Writer.Write(resp)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}

}
