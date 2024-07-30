package song_crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
)

// LyricResponse 是歌词 API 的响应结构
type LyricResponse struct {
	Lyrics string `json:"lyrics"`
}

func Crawl(c *gin.Context) {
	listsid := c.Query("listsid")
	// // "37i9dQZF1DX0kbJZpiYdZl"
	GetListsInfo(listsid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Access granted",
	})
}

// 获取歌单信息
func GetListsInfo(listID string) {
	OSSArray := []string{}

	list, err := SpotifyClient.GetPlaylist(context.Background(), spotify.ID(listID))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// iflists Redis
	exists, listsid, err := IfListsRB(listID)
	if err != nil {
		// 发生错误时回滚事务
		MysqlTX.Rollback()
		OSSRollback(OSSArray)
		return
	}
	if !exists {
		ListOSS, err := PutObjectToOSS(list.Images[0].URL)
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
		OSSArray = append(OSSArray, ListOSS)
		List := List{
			Title:       list.Name,
			Description: list.Description,
			Thumb:       ListOSS,
			CommonData: CommonData{
				IsDelete:  0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		// lists Mysql
		listsid, err = ListsDB(List)
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
		// lists Redis
		err = ListsRB(listID, strconv.Itoa(listsid))
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
	}

	MusicsArray := []Musics{}
	SingersToMusicsMap := map[string][]int{}
	SingersToMusicsArray := []SingersToMusics{}

	for _, track := range list.Tracks.Tracks {
		// 该歌曲是否已经爬取过
		exists, musicid, err := IfMusicsRB(string(track.Track.ID))
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
		if exists {
			MusicToList := MusicToList{
				MusicID: musicid,
				ListID:  listsid,
			}

			// music_to_list Mysql
			err = MusicsToListsDB(MusicToList)
			if err != nil {
				// 发生错误时回滚事务
				MysqlTX.Rollback()
				OSSRollback(OSSArray)
				return
			}
			continue
		}

		singeridlist := []int{}
		// 遍历歌手信息
		for _, artist := range track.Track.Artists {
			// ifsingers Redis
			exists, singerid, err := IfSingersRB(artist.ID.String())
			if err != nil {
				// 发生错误时回滚事务
				MysqlTX.Rollback()
				OSSRollback(OSSArray)
				return
			}
			if !exists {
				singername, thumb := GetSingersInfo(artist.ID)
				SingerOSS, err := PutObjectToOSS(thumb)
				if err != nil {
					// 发生错误时回滚事务
					MysqlTX.Rollback()
					OSSRollback(OSSArray)
					return
				}
				OSSArray = append(OSSArray, SingerOSS)
				Singers := Singers{
					Name:  singername,
					Thumb: SingerOSS,
					CommonData: CommonData{
						IsDelete:  0,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				// singers Mysql
				singerid, err = SingersDB(Singers)
				if err != nil {
					// 发生错误时回滚事务
					MysqlTX.Rollback()
					OSSRollback(OSSArray)
					return
				}
				// singers Redis
				err = SingersRB(listID, strconv.Itoa(singerid))
				if err != nil {
					// 发生错误时回滚事务
					MysqlTX.Rollback()
					OSSRollback(OSSArray)
					return
				}
			}
			singeridlist = append(singeridlist, singerid)
		}
		SingersToMusicsMap[track.Track.Name] = singeridlist

		// 获取歌词
		lyric := GetLyricsInfo(track.Track.Name, track.Track.Artists[0].Name)

		MusicOSS, err := PutObjectToOSS(track.Track.Album.Images[0].URL)
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
		OSSArray = append(OSSArray, MusicOSS)
		// 遍历歌曲信息
		Musics := Musics{
			Title:    track.Track.Name,
			Duration: int(track.Track.Duration / 1000),
			MID:      track.Track.ID.String(),
			Thumb:    MusicOSS,
			Lyrics:   lyric,
			CommonData: CommonData{
				IsDelete:  0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		MusicsArray = append(MusicsArray, Musics)
	}

	for _, Musics := range MusicsArray {
		// musics Mysql
		musicid, err := MusicsDB(Musics)
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
		// music Redis
		err = MusicsRB(Musics.MID, strconv.Itoa(musicid))
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}

		SingersToMusics := SingersToMusics{
			MusicID:   musicid,
			SingersID: SingersToMusicsMap[Musics.Title],
			CommonData: CommonData{
				IsDelete:  0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		MusicToList := MusicToList{
			MusicID: musicid,
			ListID:  listsid,
		}

		// music_to_list Mysql
		err = MusicsToListsDB(MusicToList)
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
		SingersToMusicsArray = append(SingersToMusicsArray, SingersToMusics)
	}

	// singers_to_musics Mysql
	for _, SingersToMusics := range SingersToMusicsArray {
		err = SingersToMusicsDB(SingersToMusics)
		if err != nil {
			// 发生错误时回滚事务
			MysqlTX.Rollback()
			OSSRollback(OSSArray)
			return
		}
	}

	// OSSRollback(OSSArray)
	err = MysqlTX.Commit()
	if err != nil {
		log.Fatal(err)
	}

	_, err = RedisTX.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func GetSingersInfo(artistID spotify.ID) (string, string) {
	// 获取歌手信息
	artist, err := SpotifyClient.GetArtist(context.Background(), artistID)
	if err != nil {
		log.Fatalf("couldn't get artist: %v", err)
	}

	return artist.Name, artist.Images[0].URL
}

// GetLyrics 使用 Lyrics.ovh API 获取歌词
func GetLyricsInfo(songName string, artistName string) string {
	baseURL := "https://api.lyrics.ovh/v1/"
	query := fmt.Sprintf("%s%s/%s", baseURL, url.QueryEscape(artistName), url.QueryEscape(songName))
	resp, err := http.Get(query)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var lyricResponse LyricResponse
	err = json.NewDecoder(resp.Body).Decode(&lyricResponse)
	if err != nil {
		return ""
	}

	return lyricResponse.Lyrics
}

func OSSRollback(OSS []string) {
	for _, oss := range OSS {
		DeleteObjectToOSS(oss)
	}
}
