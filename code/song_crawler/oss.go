package song_crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func PutObjectToOSS(imageURL string) (string, error) {
	// 获取后缀
	ext := GetEXT(imageURL)

	// 从网络上下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 生成UUID作为Object名称
	objectKey := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 将图片上传到OSS
	err = Bucket.PutObject(objectKey, resp.Body)
	if err != nil {
		return "", err
	}

	return "https://spotify-crawl.oss-cn-shenzhen.aliyuncs.com/" + objectKey, nil
}

func GetObjectFromOSS(key string) []byte {
	// 下载文件
	body, err := Bucket.GetObject(key)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer body.Close()

	buf, err := io.ReadAll(body)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return buf

	// fmt.Println("Downloaded object content:", string(buf))
}

func GetEXT(imageURL string) string {
	// 从网络上下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("1")
		log.Fatalf("Error: %v", err)
	}

	// 获取图片后缀
	var ext string
	// 读取前几个字节检查文件头
	buf := make([]byte, 512)
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		fmt.Println("3")
		log.Fatalf("Error: %v", err)
	}

	// 根据文件头信息判断文件类型
	fileType := http.DetectContentType(buf)
	switch fileType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/bmp":
		ext = ".bmp"
	case "image/tiff":
		ext = ".tiff"
	case "image/webp":
		ext = ".webp"
	case "image/svg+xml":
		ext = ".svg"
	// 添加更多类型的判断
	default:
		ext = "" // 无法确定扩展名时设为空
	}

	resp.Body.Close()
	return ext
}

func DeleteObjectToOSS(key string) {
	_ = Bucket.DeleteObject(key)
}
