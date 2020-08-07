package contact

import (
	"github.com/fastwego/wechat4work/corporation"
	"github.com/fastwego/wechat4work/corporation/apis/contact"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/url"
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


func Demo(c *gin.Context) {

	params := url.Values{}
	params.Add("department_id", "10086")
	resp, err := contact.UserList(ContactApp, params)

	c.Writer.Write(resp)
	c.Writer.WriteString(err.Error())
}