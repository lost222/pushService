package model

import "github.com/jinzhu/gorm"

type Record struct {
	gorm.Model
	//Recid int `gorm:"type:int;not null " json:"recid"`
	Username string `gorm:"type:varchar(20);not null " json:"username"`
	Rssurl string `gorm:"type:varchar(256);not null " json:"rssurl"`
	Fav string `gorm:"type:varchar(256);not null " json:"fav"`

}

func SearchUserSubRecord(username string) []string{
	var ans []Record

	err := db.Table("record").Where("username = ?", username).Find(&ans).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}

	var result []string
	for _, a := range ans{
		result = append(result, a.Rssurl)
	}

	return result
}

func SearchRecordUser(rssURL string) []string{
	var ans []Record

	err := db.Table("record").Where("rssurl = ?", rssURL).Find(&ans).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}

	var result []string
	for _, a := range ans{
		result = append(result, a.Username)
	}

	return result
}