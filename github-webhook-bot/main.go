package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func init() {
	// 加载配置文件
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
}

func main() {

	http.HandleFunc("/api/github/webhook", func(w http.ResponseWriter, r *http.Request) {

		payload, err := ioutil.ReadAll(r.Body)

		fmt.Println(string(payload))

		if err != nil || len(payload) == 0 {
			return
		}

		event := r.Header.Get("X-GitHub-Event")
		if event == "" {
			return
		}

		secret := viper.GetString("SECRET")
		if len(secret) > 0 {
			signature := r.Header.Get("X-Hub-Signature")
			if len(signature) == 0 {
				return
			}
			mac := hmac.New(sha1.New, []byte(secret))
			_, _ = mac.Write(payload)
			expectedMAC := hex.EncodeToString(mac.Sum(nil))

			if !hmac.Equal([]byte(signature[5:]), []byte(expectedMAC)) {
				return
			}
		}

		e := struct {
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
		}{}

		var repo string

		err = json.Unmarshal(payload, &e)
		if err != nil {
			repo = "unknown"
		}

		repo = e.Repository.FullName

		notify := struct {
			Msgtype string `json:"msgtype"`
			Markdown    struct {
				Content string `json:"content"`
			} `json:"markdown"`
		}{}

		notify.Msgtype = "markdown"
		notify.Markdown.Content = fmt.Sprintf("## Event: %s \n https://github.com/%s", event, repo)

		fmt.Println(notify)

		payload, err = json.Marshal(notify)
		response, err := http.Post(viper.GetString("WEBHOOK_URL"), "application/json", bytes.NewReader(payload))
		if err != nil {
			fmt.Println(err)
			return
		}

		all, err := ioutil.ReadAll(response.Body)
		fmt.Println(string(all), err)

	})

	_ = http.ListenAndServe(viper.GetString("LISTEN"), nil)
}
