package tus

import (
	"github.com/tus/tusd/pkg/filestore"
	"github.com/tus/tusd/pkg/handler"
	"log"
	"net/http"
	"os"
)

var TusHandler http.Handler

func InitTusHandler() http.Handler {
	uploadDir := "/tmp/tus_uploads"
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		log.Fatal("업로드 디렉토리 생성 실패 : ", err)
	}

	store := filestore.FileStore{
		Path: uploadDir,
	}

	composer := handler.NewStoreComposer()
	store.UseIn(composer)

	tusHandler, err := handler.NewHandler(handler.Config{
		BasePath:                "/files/",
		StoreComposer:           composer,
		NotifyCompleteUploads:   true,
		RespectForwardedHeaders: true,
	})

	if err != nil {
		log.Fatal("tus 핸들러 생성 실패:", err)
	}

	// 비동기 처리 진행하므로 Main 서버에서 db 처리 할 경우 api 요청 보내야 함.
	go listenUploadComplete(tusHandler, uploadDir)

	TusHandler = http.StripPrefix("/files/", tusHandler)
	return TusHandler
}
