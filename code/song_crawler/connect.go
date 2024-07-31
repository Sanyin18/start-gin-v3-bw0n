package song_crawler

import (
	"context"
	"database/sql"
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func MysqlConnect() (*sql.DB, *sql.Tx) {
	// 连接 MySQL
	mysql_connect := Mysql_userName + ":" + Mysql_password + "@tcp(" + Mysql_host + ":" + Mysql_port + ")/" + Mysql_database_name
	mysqlDB, err := sql.Open("mysql", mysql_connect)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()
	// 测试 MySQL 连接
	if err := mysqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	// 开始事务
	mysqlTX, err := mysqlDB.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	return mysqlDB, mysqlTX
}

func RedisConnect() (*redis.Client, redis.Pipeliner) {
	// 连接 Redis
	redisDB := redis.NewClient(&redis.Options{
		Addr:     Redis_host + ":" + Redis_port,
		Password: Redis_password, // 密码
		DB:       0,              // 使用默认数据库
	})

	// 测试 Redis 连接
	_, err := redisDB.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 开始事务
	redisTX := redisDB.TxPipeline()

	return redisDB, redisTX
}

func OSSConnect() *oss.Bucket {
	// 创建OSS客户端
	client, err := oss.New(OSS_Endpoint, OSS_AccessKey_ID, OSS_AccessKey_Secret)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 获取存储空间
	Bucket, err := client.Bucket(OSS_BucketName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return Bucket
}

func SpotifyConnect() *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     SpotifyClientID,
		ClientSecret: SpotifyClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("无法获取token: %v", err)
	}

	httpClient := spotifyauth.New().Client(context.Background(), token)
	SpotifyClient := spotify.New(httpClient)
	return SpotifyClient
}
