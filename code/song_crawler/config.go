package song_crawler

import (
	"database/sql"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis"
	"github.com/zmb3/spotify/v2"
)

var (
	Mysql_host           string
	Mysql_port           string
	Mysql_userName       string
	Mysql_password       string
	Mysql_database_name  string
	Redis_host           string
	Redis_port           string
	Redis_userName       string
	Redis_password       string
	OSS_Endpoint         string
	OSS_AccessKey_ID     string
	OSS_AccessKey_Secret string
	OSS_BucketName       string
	SpotifyClientID      string
	SpotifyClientSecret  string
	JWTKeyString         string
	Admin                string
	Password             string

	JWTKey = []byte(JWTKeyString)

	AdminUser = map[string]string{
		Admin: Password,
	}

	MysqlDB       *sql.DB
	RedisDB       *redis.Client
	Bucket        *oss.Bucket
	SpotifyClient *spotify.Client

	MysqlTX *sql.Tx
	RedisTX redis.Pipeliner
)
