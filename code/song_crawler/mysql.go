package song_crawler

import "encoding/json"

func ListsDB(list List) (int, error) {
	// 使用事务 tx 执行插入操作
	result, err := MysqlTX.Exec(
		"INSERT INTO lists (title, description, thumb, is_delete, create_date, update_date) VALUES (?, ?, ?, ?, ?, ?)",
		list.Title, []byte(list.Description), list.Thumb, list.CommonData.IsDelete, list.CommonData.CreatedAt, list.CommonData.UpdatedAt)

	if err != nil {
		return 0, err
	}

	listid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(listid), nil
}
func SingersDB(singers Singers) (int, error) {
	// 使用事务 tx 执行插入操作
	result, err := MysqlTX.Exec(
		"INSERT INTO singers (name, thumb, is_delete, create_date, update_date) VALUES (?, ?, ?, ?, ?)",
		singers.Name, singers.Thumb, singers.CommonData.IsDelete, singers.CommonData.CreatedAt, singers.CommonData.UpdatedAt)

	if err != nil {
		return 0, err
	}

	singerid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(singerid), nil
}
func MusicsDB(musics Musics) (int, error) {
	// 使用事务 tx 执行插入操作
	result, err := MysqlTX.Exec(
		"INSERT INTO musics (title, duration, mid, thumb, lyrics, is_delete, create_date, update_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		musics.Title, musics.Duration, musics.MID, musics.Thumb, []byte(musics.Lyrics), musics.CommonData.IsDelete, musics.CommonData.CreatedAt, musics.CommonData.UpdatedAt)

	if err != nil {
		return 0, err
	}

	musicid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(musicid), nil
}
func MusicsToListsDB(musicToList MusicToList) error {
	// 使用事务 tx 执行插入操作
	_, err := MysqlTX.Exec(
		"INSERT INTO musics_to_lists (music_id, list_id) VALUES (?, ?)",
		musicToList.MusicID, musicToList.ListID)

	return err
}

func SingersToMusicsDB(singersToMusics SingersToMusics) error {
	singersIDsJSON, err := json.Marshal(singersToMusics.SingersID)
	if err != nil {
		return err
	}
	// 使用事务 tx 执行插入操作
	_, err = MysqlTX.Exec(
		"INSERT INTO singers_to_musics (music_id, singers_id, is_delete, create_date, update_date) VALUES (?, ?, ?, ?, ?)",
		singersToMusics.MusicID, singersIDsJSON, singersToMusics.CommonData.IsDelete, singersToMusics.CommonData.CreatedAt, singersToMusics.CommonData.UpdatedAt)

	return err
}
