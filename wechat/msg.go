package wechat

import (
	"encoding/xml"
	"time"
)

type Msg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName,omitempty"`
	FromUserName string   `xml:"FromUserName,omitempty"`
	CreateTime   int64    `xml:"CreateTime,omitempty"`
	MsgType      string   `xml:"MsgType,omitempty"`
	Event        string   `xml:"Event,omitempty"`
	Content      string   `xml:"Content,omitempty"`
	Recognition  string   `xml:"Recognition,omitempty"`
	MsgId        int64    `xml:"MsgId,omitempty"`
	EventKey     string   `xml:"EventKey,omitempty"`
	PicUrl       string   `xml:"PicUrl,omitempty"`
	MediaId      string   `xml:"MediaId,omitempty"`
	Encrypt      string   `xml:"Encrypt,omitempty"`
}

func NewMsg(data []byte) *Msg {
	var msg Msg
	if err := xml.Unmarshal(data, &msg); err != nil {
		return nil
	}
	return &msg
}

func (msg *Msg) GenerateEchoData(s string) []byte {
	data := Msg{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      s,
	}
	bs, _ := xml.Marshal(&data)
	return bs
}
