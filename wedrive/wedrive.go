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

package wedrive

import (
	"github.com/fastwego/wxwork/corporation"
	"github.com/fastwego/wxwork/corporation/apis/efficiency/wedrive"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Corp *corporation.Corporation
var WedriveApp *corporation.App

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	Corp = corporation.New(corporation.Config{Corpid: viper.GetString("CROPID")})
	WedriveApp = Corp.NewApp(corporation.AppConfig{
		Secret: viper.GetString("WeDriveSecret"),
	})
}

func Demo(c *gin.Context) {

	payload := []byte(`{
    "userid": "USERID",
    "spaceid": "SPACEID"
}`)
	resp, err := wedrive.SpaceInfo(WedriveApp, payload)

	c.Writer.Write(resp)
	if err != nil {
		c.Writer.WriteString(err.Error())
	}

}
