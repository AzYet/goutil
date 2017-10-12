package goutil

import (
	"bytes"
	"net/http"
	"encoding/json"
	"time"
	"github.com/Sirupsen/logrus"
	"fmt"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Action     []map[string]string `json:"action"`
	Fallback   string              `json:"fallback"`
	Color      string              `json:"color"`
	PreText    string              `json:"pretext"`
	AuthorName string              `json:"author_name"`
	AuthorLink string              `json:"author_link"`
	AuthorIcon string              `json:"author_icon"`
	Title      string              `json:"title"`
	TitleLink  string              `json:"title_link"`
	Text       string              `json:"text"`
	ImageUrl   string              `json:"image_url"`
	Fields     []*Field            `json:"fields"`
	Footer     string              `json:"footer"`
	FooterIcon string              `json:"footer_icon"`
	TimeStamp  int64               `json:"ts"`
	MarkdownIn []string            `json:"mrkdwn_in,omitempty"`
}

type SlackData struct {
	Parse       string        `json:"parse,omitempty"`
	Username    string        `json:"username,omitempty"`
	IconUrl     string        `json:"icon_url,omitempty"`
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Channel     string        `json:"channel,omitempty"`
	Text        string        `json:"text,omitempty"`
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

func (sd *SlackData) Send(dsn ...string) {
	FeedSlack(dsn, sd, logrus.New())
}

func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

func NewAttatchment(title string) *Attachment {
	return &Attachment{
		Title:     title,
		Color:     "#FF0000",
		TimeStamp: time.Now().Unix(),
	}
}

func FeedSlack(dsn []string, d *SlackData, Log *logrus.Logger) {
	client := &http.Client{Timeout: time.Second * 45}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(d)
	if err != nil {
		Log.Warnln("json marshal err", err)
		return
	}
	Log.Println(fmt.Sprintf("%+v", d))
	url := "https://hooks.slack.com/services/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"
	if len(dsn) > 0 {
		url = dsn[0]
	}
	r, err := http.NewRequest("POST", url, buf)
	if err != nil {
		Log.Warnln("create http post err", err)
		return
	}
	r.Header.Add("Content-Type", "application/json")
	for i := 0; i < 3; i++ {
		if resp, err := client.Do(r); err != nil {
			Log.Warnln("post err", err)
		} else {
			Log.Printf("Slack api updated, status %v.", resp.StatusCode)
			break
		}
	}
}
