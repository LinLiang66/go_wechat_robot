package redisService

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"go_wechat_robot/model"
	"io"
	"log"
	"net/http"
	"strconv"
)

type RedisUtil struct {
	client redis.Conn
}

// RedisClient 全局变量, 外部使用utils.RedisClient来访问
var RedisClient RedisUtil

// InitRedisUtil 初始化redis
func InitRedisUtil(address string, port int, pwd string) (*RedisUtil, error) {
	//连接redis
	client, err := redis.Dial("tcp", address+":"+strconv.Itoa(port))
	if err != nil {
		panic("failed to redis:" + err.Error())
	}
	//验证redis redis的配置文件redis.conf中一定要设置quirepass=password, 不然连不上
	_, err = client.Do("auth", pwd)
	if err != nil {
		panic("failed to auth redis:" + err.Error())
	}
	//初始化全局redis结构体
	RedisClient = RedisUtil{client: client}
	return &RedisClient, nil
}

// SetStr 设置数据到redis中（string）
func (rs *RedisUtil) SetStr(key string, value string) error {
	_, err := rs.client.Do("Set", key, value)
	return err
}

// SetStrNotExist 设置数据到redis中（string）
func (rs *RedisUtil) SetStrNotExist(key string, value string, expireSecond int) bool {
	val, err := rs.client.Do("SET", key, value, "EX", expireSecond, "NX")
	if err != nil || val == nil {
		return false
	}
	return true
}

// SetStrWithExpire 设置数据到redis中（string）
func (rs *RedisUtil) SetStrWithExpire(key string, value string, expireSecond int) error {
	_, err := rs.client.Do("Set", key, value, "ex", expireSecond)
	return err
}

// GetStr 获取redis中数据（string）
func (rs *RedisUtil) GetStr(key string) (string, error) {
	val, err := rs.client.Do("Get", key)
	if err != nil {
		return "", err
	}
	res := val.([]byte)
	return string(res), nil
}

// HSet 设置数据到redis中（hash）
func (rs *RedisUtil) HSet(key string, field string, value string) error {
	_, err := rs.client.Do("HSet", key, field, value)
	return err
}

// HGet 获取redis中数据（hash）
func (rs *RedisUtil) HGet(key string, field string) (string, error) {
	val, err := rs.client.Do("HGet", key, field)
	if err != nil {
		return "", err
	}
	return string(val.([]byte)), nil
}

// DelByKey 删除
func (rs *RedisUtil) DelByKey(key string) error {
	_, err := rs.client.Do("DEL", key)
	return err
}

// SetExpire 设置key过期时间
func (rs *RedisUtil) SetExpire(key string, expireSecond int) error {
	_, err := rs.client.Do("EXPIRE", key, expireSecond)
	return err
}

// KEYEXISTS 判断KEY在redis中是否存在
func (rs *RedisUtil) KEYEXISTS(KEY string) bool {
	exists, _ := redis.Bool(rs.client.Do("EXISTS", KEY))
	return exists
}

// KEYEXISTSGetStr 判断KEY在redis中是否存在,存在则获取内容
func (rs *RedisUtil) KEYEXISTSGetStr(KEY string) (bool, string) {
	exists, _ := redis.Bool(rs.client.Do("EXISTS", KEY))
	if exists {
		val, err := rs.client.Do("Get", KEY)
		if err != nil {
			return false, ""
		}
		return true, string(val.([]byte))
	}

	return exists, ""
}

// GetWechatCache 获取公众号缓存信息
func GetWechatCache(appId string) (model.WechatAppCache, bool) {
	if RedisClient.KEYEXISTS("robot:wechat_app_key:" + appId) {
		str, err := RedisClient.GetStr("robot:wechat_app_key:" + appId)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		var appCache model.WechatAppCache
		err = json.Unmarshal([]byte(str), &appCache)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		return appCache, true
	}
	return model.WechatAppCache{}, false
}

func GetWechatToken(appid string, AppCache model.WechatAppCache) (model.WechatAppCache, bool) {
	var appCache model.WechatAppCache
	//判断 redis中是否有accessToken,有则redis中取，没有就微信官方取，因为设置了redis中的accessToken有效期，过期自动失效
	if RedisClient.KEYEXISTS("robot:wehcat_app:access_token:" + appid) {
		str, err := RedisClient.GetStr("robot:wehcat_app:access_token:" + appid)
		if err != nil {
			log.Printf("获取access_token 报错: %v\n", err)
			return appCache, false
		}
		err = json.Unmarshal([]byte(str), &appCache)
		if err != nil {
			log.Printf("failed to start server: %v\n", err)
		}
		return appCache, true
	} else {
		var AccessToken string
		req, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + AppCache.AppID + "&secret=" + AppCache.AppSecret)
		if err != nil {
			log.Printf("读取请求结果报错: %v\n", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
			}
		}(req.Body)
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("读取请求结果报错: %v", err)
		}
		var res model.RequestEntity
		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Printf("请求结果转换实体报错: %v", err)
		}
		AccessToken = res.AccessToken
		appCache.AppID = appid
		appCache.AccessToken = AccessToken
		bytes, _ := json.Marshal(appCache)
		RedisClient.SetStrWithExpire("robot:wehcat_app:access_token:"+appid, string(bytes), 7000)
		return appCache, true
	}
}

func GetMessageContext(userid string) []model.IFlytekMessage {
	var MessageContext []model.IFlytekMessage
	if RedisClient.KEYEXISTS("robot:message_context:" + userid) {
		str, err := RedisClient.GetStr("robot:message_context:" + userid)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		err = json.Unmarshal([]byte(str), &MessageContext)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		return MessageContext
	}
	return MessageContext
}

func SetMessageContext(userid string, MessageContext []model.IFlytekMessage) {
	bytes, _ := json.Marshal(MessageContext)
	RedisClient.SetStr("robot:message_context:"+userid, string(bytes))
}

func GetUserCache(userid string) (model.AppCache, bool) {
	if RedisClient.KEYEXISTS("robot:robot_user_model:" + userid) {
		str, err := RedisClient.GetStr("robot:robot_user_model:" + userid)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		var appCache model.AppCache
		err = json.Unmarshal([]byte(str), &appCache)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		return appCache, true
	}
	appCache := model.AppCache{
		UserID:           userid,
		RobotAppID:       "讯飞控制台获取到的APPID",
		RobotApiSecret:   "讯飞控制台获取到的APISecret",
		RobotApiKey:      "讯飞控制台获取到的APIKey",
		RobotDomain:      "spark3.5-chat",                       //默认最新版
		RobotSparkUrl:    "ws://spark-api.xf-yun.com/v3.5/chat", //默认最新版
		RobotTemperature: 1}                                     //1 发散模式为严谨
	bytes, _ := json.Marshal(appCache)
	RedisClient.SetStr("robot:robot_user_model:"+userid, string(bytes))
	return appCache, true
}
