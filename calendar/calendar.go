package calendar

import (
	"github.com/fastwego/wechat4work/corporation"
	"github.com/fastwego/wechat4work/corporation/apis/calendar"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var Corp *corporation.Corporation
var CalendarApp *corporation.App

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	Corp = corporation.New(corporation.CorporationConfig{Corpid: viper.GetString("CROPID")})
	CalendarApp = Corp.NewApp(corporation.AppConfig{
		Secret:         viper.GetString("CalendarSECRET"),
	})
}


func Demo(c *gin.Context) {

	payload := []byte(``)
	resp, err := calendar.CalendarGet(CalendarApp, payload)

	c.Writer.Write(resp)
	c.Writer.WriteString(err.Error())
}