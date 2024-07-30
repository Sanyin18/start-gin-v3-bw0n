package song_crawler

import "os"

func Getenv() {
	Mysql_host = os.Getenv("Mysql_host")
	Mysql_port = os.Getenv("Mysql_port")
	Mysql_userName = os.Getenv("Mysql_userName")
	Mysql_password = os.Getenv("Mysql_password")
	Mysql_database_name = os.Getenv("Mysql_database_name")
	Redis_host = os.Getenv("Redis_host")
	Redis_port = os.Getenv("Redis_port")
	Redis_userName = os.Getenv("Redis_userName")
	Redis_password = os.Getenv("Redis_password")
	OSS_Endpoint = os.Getenv("OSS_Endpoint")
	OSS_AccessKey_ID = os.Getenv("OSS_AccessKey_ID")
	OSS_AccessKey_Secret = os.Getenv("OSS_AccessKey_Secret")
	OSS_BucketName = os.Getenv("OSS_BucketName")
	SpotifyClientID = os.Getenv("SpotifyClientID")
	SpotifyClientSecret = os.Getenv("SpotifyClientSecret")
	JWTKeyString = os.Getenv("JWTKeyString")
	Admin = os.Getenv("Admin")
	Password = os.Getenv("Password")
}
