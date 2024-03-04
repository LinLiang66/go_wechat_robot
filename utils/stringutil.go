package utils

import (
	"encoding/xml"
	"go_wechat_robot/model"
)

// H is a shortcut for map[string]any
type H map[string]any

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "xml",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
func getLength(Messages []model.IFlytekMessage) int {
	length := 0
	for _, Message := range Messages {
		temp := Message.Content
		leng := len(temp)
		length += leng
	}
	return length
}

func Checklen(Messages []model.IFlytekMessage) []model.IFlytekMessage {
	for getLength(Messages) > 8000 {
		Messages = Messages[1:]
	}
	return Messages
}

func GetDomainType(DomainType string) string {
	switch DomainType {
	case "spark1.5-chat":
		return "general"
	case "spark2.1-chat":
		return "generalv2"
	case "spark3.1-chat":
		return "generalv3"
	case "spark3.5-chat":
		return "generalv3.5"
	default:
		break
	}
	return ""
}
