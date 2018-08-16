package goutil

import (
	"bytes"
	"net/http"
	"encoding/json"
	"time"
	"github.com/sirupsen/logrus"
	"fmt"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Actions    []map[string]string `json:"actions,omitempty"`
	Fallback   string              `json:"fallback,omitempty"`
	Color      string              `json:"color,omitempty"`
	PreText    string              `json:"pretext,omitempty"`
	AuthorName string              `json:"author_name,omitempty"`
	AuthorLink string              `json:"author_link,omitempty"`
	AuthorIcon string              `json:"author_icon,omitempty"`
	Title      string              `json:"title,omitempty"`
	TitleLink  string              `json:"title_link,omitempty"`
	Text       string              `json:"text,omitempty"`
	ImageUrl   string              `json:"image_url,omitempty"`
	Fields     []*Field            `json:"fields,omitempty"`
	Footer     string              `json:"footer,omitempty"`
	FooterIcon string              `json:"footer_icon,omitempty"`
	TimeStamp  int64               `json:"ts,omitempty"`
	MarkdownIn []string            `json:"mrkdwn_in,omitempty"`
	CallbackID string              `json:"callback_id,omitempty"`
}

type SlackData struct {
	Token       string        `json:"token,omitempty"`
	Parse       string        `json:"parse,omitempty"`
	Username    string        `json:"username,omitempty"`
	AsUser      string        `json:"as_user,omitempty"`
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
	sd.Attachments = append(sd.Attachments, as...)
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
	url := ""
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
