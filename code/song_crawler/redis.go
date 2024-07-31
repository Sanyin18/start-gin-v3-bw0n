package song_crawler

import (
	"strconv"

	"github.com/go-redis/redis"
)

func ListsRB(listsid string, id string) error {
	// 将字符串添加到 Redis Hash
	err := RedisTX.HSet("lists", listsid, id).Err()
	return err
}

func SingersRB(singersid string, id string) error {
	// 将字符串添加到 Redis Hash
	err := RedisTX.HSet("singers", singersid, id).Err()
	return err
}

func MusicsRB(musicsid string, id string) error {
	// 将字符串添加到 Redis Hash
	err := RedisTX.HSet("musics", musicsid, id).Err()
	return err
}

func IfListsRB(listsid string) (bool, int, error) {
	// 在 "lists" 中查找字符串
	idstring, err := RedisDB.HGet("lists", listsid).Result()
	if err == redis.Nil {
		// 查无该值
		return false, 0, nil
	} else if err != nil {
		return false, 0, err
	}

	id, _ := strconv.Atoi(idstring)

	return true, id, nil
}

func IfSingersRB(singersid string) (bool, int, error) {
	// 在 "lists" 中查找字符串
	idstring, err := RedisDB.HGet("singers", singersid).Result()
	if err == redis.Nil {
		// 查无该值
		return false, 0, nil
	} else if err != nil {
		return false, 0, err
	}

	id, _ := strconv.Atoi(idstring)

	return true, id, nil
}

func IfMusicsRB(musicsid string) (bool, int, error) {
	// 在 "lists" 中查找字符串
	idstring, err := RedisDB.HGet("musics", musicsid).Result()
	if err == redis.Nil {
		// 查无该值
		return false, 0, nil
	} else if err != nil {
		return false, 0, err
	}

	id, _ := strconv.Atoi(idstring)

	return true, id, nil
}
