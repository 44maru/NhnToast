package flow

import (
	"log"

	"nhn-toast/pkg/api/toast"
)

func ListImage(token string) {

	imageList, err := toast.GetImageList(token)
	if err != nil {
		log.Printf("ERROR: イメージリスト取得失敗: %s\n", err.Error())
		return
	}

	for _, imageInfo := range imageList.Images {
		log.Printf("%s %s\n", imageInfo.ID, imageInfo.Name)
	}

}
