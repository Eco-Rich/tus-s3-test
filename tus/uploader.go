package tus

import (
	"github.com/tus/tusd/pkg/handler"
	"log"
)

func listenUploadComplete(tusHandler *handler.Handler) {
	for {
		event := <-tusHandler.CompleteUploads
		log.Printf("업로드 완료! ID: %s, 메타데이터: %+v\n", event.Upload.ID, event.Upload.MetaData)
		// 여따가 db로직
	}
}
