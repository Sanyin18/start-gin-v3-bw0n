package main

import (
	"song/song_crawler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 读取环境变量等信息.
	song_crawler.Getenv()

	// 连接
	song_crawler.MysqlDB, song_crawler.MysqlTX = song_crawler.MysqlConnect()
	song_crawler.RedisDB, song_crawler.RedisTX = song_crawler.RedisConnect()
	song_crawler.Bucket = song_crawler.OSSConnect()
	song_crawler.SpotifyClient = song_crawler.SpotifyConnect()

	r := gin.Default()

	// 定义 GET 方法,只有拥有有效 JWT 的用户才能访问
	r.GET("/crawl", song_crawler.JWTAuthMiddleware(), song_crawler.Crawl)
	r.POST("/login", song_crawler.Login)

	r.Run()
}
