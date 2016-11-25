package go_utils

import (
	"bytes"
	"net/http"
	"encoding/json"
	"time"
	"github.com/Sirupsen/logrus"
)

type SlackData struct {
	Text        string `json:"text"`
	UserName    string `json:"username"`
	IconUrl     string `json:"icon_url"`
	Attachments []map[string]string `json:"attachments"`
}

func FeedSlack(d *SlackData, Logger *logrus.Logger) {
	client := &http.Client{Timeout:time.Second * 45}
	bz, err := json.Marshal(*d)
	if err != nil {
		Logger.Println("json marshal err", err)
	} else {
		Logger.Println(string(bz))
		if r, err := http.NewRequest("POST", "https://hooks.slack.com/services/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", bytes.NewBuffer(bz)); err != nil {
			Logger.Println("create http post err", err)
		} else {
			r.Header.Add("Content-Type", "application/json")
			for i := 0; i < 1; i++ {
				if resp, err := client.Do(r); err != nil {
					Logger.Println("post err", err)
				} else {
					Logger.Printf("Slack api updated, status %v.", resp.StatusCode)
					break
				}
			}
		}
	}
}