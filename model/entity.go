package model

type WechatAppCache struct {
	AppID       string `json:"appid"`
	AppSecret   string `json:"app_secret"`
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	AesKey      string `json:"aes_key"`
}

// RequestEntity 公众号 通用实体
type RequestEntity struct {
	Errcode     int    `json:"errcode,omitempty"`
	Errmsg      string `json:"errmsg,omitempty"`
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

// IFlytekMessage 讯飞星火大模型
type IFlytekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AppCache struct {
	UserName          string  `json:"user_name"`
	UserID            string  `json:"user_id"`
	AppID             string  `json:"appid"`
	AppSecret         string  `json:"app_secret"`
	AppRoleType       float64 `json:"app_role_type"`
	VerificationToken string  `json:"verification_token"`
	EncryptKey        string  `json:"encrypt_key"`
	RobotAppID        string  `json:"robot_appid"`
	RobotApiSecret    string  `json:"robot_api_secret"`
	RobotApiKey       string  `json:"robot_api_key"`
	RobotDomain       string  `json:"robot_domain"`
	RobotSparkUrl     string  `json:"robot_spark_url"`
	RobotTemperature  float64 `json:"robot_temperature"`
}
