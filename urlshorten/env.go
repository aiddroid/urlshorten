package urlshorten

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	S Storage
}

func GetEnv() *Env {
	//从系统环境变量中读取配置
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("APP_REDIS_PASSWORD")
	dbs := os.Getenv("APP_REDIS_DB")
	db, _ := strconv.Atoi(dbs)

	log.Println("Connecting to redis:", addr)
	r := NewRedisCli(addr, password, db)

	return &Env{S: r}
}
