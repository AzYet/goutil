package go_utils

import (
	"bytes"
	"net/http"
	"encoding/json"
	"time"
	"github.com/Sirupsen/logrus"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Fallback   string   `json:"fallback"`
	Color      string   `json:"color"`
	PreText    string   `json:"pretext"`
	AuthorName string   `json:"author_name"`
	AuthorLink string   `json:"author_link"`
	AuthorIcon string   `json:"author_icon"`
	Title      string   `json:"title"`
	TitleLink  string   `json:"title_link"`
	Text       string   `json:"text"`
	ImageUrl   string   `json:"image_url"`
	Fields     []*Field `json:"fields"`
	Footer     string   `json:"footer"`
	FooterIcon string   `json:"footer_icon"`
	TimeStamp  int64   `json:"ts"`
}

type SlackData struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconUrl     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

func (sd *SlackData) Attach(as ...*Attachment) {
	if sd == nil {
		return
	}
	if sd.Attachments == nil {
		sd.Attachments = make([]*Attachment, 0, len(as))
	}
	for _, a := range as {
		sd.Attachments = append(sd.Attachments, a)
	}

}
func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

func NewAttatchment(title string) *Attachment {
	return &Attachment{
		Title:title,
		Color:"#FF0000",
		TimeStamp:time.Now().Unix(),
	}
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
			for i := 0; i < 3; i++ {
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