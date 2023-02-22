package app

import (
	"time"

	"github.com/visonlv/go-vkit/miniox"
	"github.com/visonlv/go-vkit/mongox"
	"github.com/visonlv/go-vkit/mysqlx"
	"github.com/visonlv/go-vkit/natsx"
	"github.com/visonlv/go-vkit/redisx"
	"github.com/visonlv/vkit-example/config"
)

var (
	Cfg   *config.Config
	Mysql *mysqlx.MysqlClient
	Minio *miniox.MinioClient
	Redis *redisx.RedisClient
	Mongo *mongox.MongoClient
	Nats  *natsx.NatsClient
)

type clientProxy struct {
}

func Init(cfgPath string) {
	initConfig(cfgPath)
	initMysql()
	// initMinIO()
	// initRedis()
	// initMongoDb()
	// initNats()
}

// 初始化配置， 第一步运行
func initConfig(fpath string) {
	if cfg, err := config.LoadConfig(fpath); err == nil {
		Cfg = cfg
	} else {
		panic(err)
	}
}

func initMysql() {
	c, err := mysqlx.NewClient(Cfg.Mysql.Uri, Cfg.Mysql.MaxConn, Cfg.Mysql.MaxIdel, Cfg.Mysql.MaxLifeTime)
	if err != nil {
		panic(err)
	}
	Mysql = c
}

func initMinIO() {
	c, err := miniox.NewClient(Cfg.MinIO.DoMain, Cfg.MinIO.EndPoint,
		Cfg.MinIO.AccessKey, Cfg.MinIO.AccessSecret,
		Cfg.MinIO.BucketName)
	if err != nil {
		panic(err)
	}
	Minio = c
}

func initRedis() {
	c, err := redisx.NewClient(Cfg.Redis.Address, Cfg.Redis.Password, Cfg.Redis.Db)
	if err != nil {
		panic(err)
	}
	Redis = c
}

func initMongoDb() {
	c, err := mongox.NewClient(Cfg.MongoDB.URI,
		Cfg.MongoDB.DbName,
		time.Duration(Cfg.MongoDB.WithTimeout)*time.Second)

	if err != nil {
		panic(err)
	}
	Mongo = c
}

func initNats() {
	c, err := natsx.NewClient(Cfg.Nats.Url,
		Cfg.Nats.Username, Cfg.Nats.Password)
	if err != nil {
		panic(err)
	}
	Nats = c
}
