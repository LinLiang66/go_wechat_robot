package model

// MsgEntity 消息结构体
type MsgEntity struct {
	Image         Image         `json:"image,omitempty"`
	Mpnews        Mpnews        `json:"mpnews,omitempty"`
	Mpnewsarticle Mpnewsarticle `json:"mpnewsarticle,omitempty"`
	Msgmenu       Msgmenu       `json:"msgmenu,omitempty"`
	Msgtype       string        `json:"msgtype,omitempty"`
	Music         Music         `json:"music,omitempty"`
	News          News          `json:"news,omitempty"`
	Text          Text          `json:"text,omitempty"`
	Touser        string        `json:"touser,omitempty"`
	Video         Video         `json:"video,omitempty"`
	Voice         Voice         `json:"voice,omitempty"`
	Wxcard        Wxcard        `json:"wxcard,omitempty"`
}

type Image struct {
	MediaID string `json:"media_id,omitempty"`
}

type Mpnews struct {
	MediaID string `json:"media_id,omitempty"`
}

type Mpnewsarticle struct {
	ArticleID string `json:"article_id,omitempty"`
}

type Msgmenu struct {
	HeadContent string `json:"head_content,omitempty"`
	List        []List `json:"list,omitempty"`
	TailContent string `json:"tail_content,omitempty"`
}

type List struct {
	Content string `json:"content,omitempty"`
	ID      string `json:"id,omitempty"`
}

type Music struct {
	Description  string `json:"description,omitempty"`
	Hqmusicurl   string `json:"hqmusicurl,omitempty"`
	Musicurl     string `json:"musicurl,omitempty"`
	ThumbMediaID string `json:"thumb_media_id,omitempty"`
	Title        string `json:"title,omitempty"`
}

type News struct {
	Articles []Article `json:"articles,omitempty"`
}

type Article struct {
	Description string `json:"description,omitempty"`
	Picurl      string `json:"picurl,omitempty"`
	Title       string `json:"title,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Text struct {
	Content string `json:"content,omitempty"`
}

type Video struct {
	Description  string `json:"description,omitempty"`
	MediaID      string `json:"media_id,omitempty"`
	ThumbMediaID string `json:"thumb_media_id,omitempty"`
	Title        string `json:"title,omitempty"`
}

type Voice struct {
	MediaID string `json:"media_id,omitempty"`
}

type Wxcard struct {
	CardID string `json:"card_id,omitempty"`
}
