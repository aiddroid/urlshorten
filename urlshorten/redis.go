package urlshorten

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
	"log"
	"time"
)

const (
	URL_ID_KEY          = "next.url.id"
	SHORT_LINK_KEY      = "shortlink:%s"
	URL_HASH_KEY        = "urlhash:%s"
	SHORT_LINK_INFO_KEY = "shortlinkinfo:%s"
)

type RedisCli struct {
	//持有第三方redis库实例
	Cli *redis.Client
}

func NewRedisCli(addr string, password string, db int) *RedisCli {
	//创建RedisCli
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}

	return &RedisCli{Cli: c}
}

//字符串转sha1
func ToSha1(s string) interface{} {
	hash := sha1.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//url转换为短链接
func (r *RedisCli) Shorten(url string, expire int) (string, error) {
	h := ToSha1(url)

	// 检查url是否已经转换过，已经存在
	urlHashKey := fmt.Sprintf(URL_HASH_KEY, h)
	d, err := r.Cli.Get(urlHashKey).Result()
	log.Println("URL hash:", h)
	if err != nil && err != redis.Nil {
		return "", err
	} else if d != "" {
		if d == "{}" {
			d = ""
		}

		log.Println("Cache HIT!")

		return d, nil
	}

	//id自增
	id, err := r.Cli.Incr(URL_ID_KEY).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	//计算短地址
	eid := base62.EncodeInt64(id)

	log.Println("Starting redis TX pipeline")

	//开始redis事务
	pipeline := r.Cli.TxPipeline()
	exp := time.Duration(expire) * time.Second

	//设置hashkey
	pipeline.Set(urlHashKey, eid, exp)

	//设置短地址到URL的映射
	shortLinkKey := fmt.Sprintf(SHORT_LINK_KEY, eid)
	pipeline.Set(shortLinkKey, url, exp)

	//设置短地址详情信息
	shortLinkInfoKey := fmt.Sprintf(SHORT_LINK_INFO_KEY, eid)
	m := make(map[string]interface{})
	m["url"] = url
	m["created_at"] = time.Now().Unix()
	m["expire"] = expire
	r.Cli.HMSet(shortLinkInfoKey, m)
	r.Cli.Expire(shortLinkInfoKey, exp)

	//提交事务
	if _, err := pipeline.Exec(); err != nil {
		log.Println("Redis TX exec error:", err)
		return "", err
	}

	log.Println("eid: ", eid)

	return eid, err
}

//短链接转换为原始url
func (r *RedisCli) UnShorten(eid string) (string, error) {
	shortLinkKey := fmt.Sprintf(SHORT_LINK_KEY, eid)
	url, err := r.Cli.Get(shortLinkKey).Result()
	if url != "" {
		return url, nil
	}
	return "", err
}

//获取短链接的详细信息
func (r *RedisCli) ShortenInfo(eid string) (interface{}, error) {
	shortLinkInfoKey := fmt.Sprintf(SHORT_LINK_INFO_KEY, eid)
	m, err := r.Cli.HGetAll(shortLinkInfoKey).Result()
	if err != nil {
		return nil, err
	} else if len(m) == 0 {
		return nil, errors.New("info not found")
	}

	return m, nil
}
