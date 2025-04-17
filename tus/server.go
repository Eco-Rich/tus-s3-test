package tus

import (
	"gin/config"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
	"log"
	"net/http"
	"os"
)

var TusHandler http.Handler

func InitTusHandler() http.Handler {
	sess, err := config.NewAWSSession()
	if err != nil {
		log.Fatal("AWS 세션 생성실패:", err)
	}

	s3Client := s3.New(sess)
	store := s3store.New(os.Getenv("S3_BUCKET"), s3Client)
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

	go listenUploadComplete(tusHandler)

	TusHandler = http.StripPrefix("/files/", tusHandler)
	return TusHandler
}
