package song_crawler

import (
	"time"
)

type CommonData struct {
	MetaData  interface{} `json:"meta_data" description:"元数据"`
	IsDelete  int         `json:"is_delete" description:"是否软删除"`
	CreatedAt time.Time   `json:"create_at" description:"创建时间"`
	UpdatedAt time.Time   `json:"update_at" description:"更新时间"`
	DeletedAt time.Time   `json:"delete_at" description:"删除时间"`
}

type Singers struct {
	Id    int    `json:"id" description:"自增ID"`
	Name  string `json:"name" description:"用户名称"`
	Thumb string `json:"thumb" description:"缩略图地址"`
	CommonData
}

type Musics struct {
	ID       int    `json:"id" description:"自增ID"`
	Title    string `json:"title" description:"标题"`
	Duration int    `json:"duration" description:"音乐时长"`
	MID      string `json:"mid" description:"第三方ID"`
	Thumb    string `json:"thumb" description:"缩略图地址"`
	Lyrics   string `json:"lyrics" description:"歌词"`
	SingerID int    `json:"singer_id" description:"歌手ID"`
	CommonData
}

type List struct {
	ID          int    `json:"id" description:"自增ID"`
	Title       string `json:"title" description:"标题"`
	Description string `json:"description" description:"描述"`
	Thumb       string `json:"thumb" description:"缩略图地址"`
	CommonData
}

type MusicToList struct {
	ID      int `json:"id" description:"自增ID"`
	MusicID int `json:"music_id" description:"歌曲ID"`
	ListID  int `json:"list_id" description:"列表ID"`
	CommonData
}

type SingersToMusics struct {
	ID        int   `json:"id" description:"自增ID"`
	MusicID   int   `json:"music_id" description:"歌曲ID"`
	SingersID []int `json:"singers_id" description:"歌手ID"`
	CommonData
}

type Local struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	CommonData
}

type Explain struct {
	ID      int    `json:"id"`
	MusicID int    `json:"music_id"`
	LocalID int    `json:"local_id"`
	Content []byte `json:"content"`
	CommonData
}

// ===========================================================================================================================================================
