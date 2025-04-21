package tus

import (
	"github.com/tus/tusd/pkg/handler"
	"log"
	"os/exec"
	"path/filepath"
)

func listenUploadComplete(tusHandler *handler.Handler, uploadDir string) {
	for {
		event := <-tusHandler.CompleteUploads
		uploadID := event.Upload.ID
		meta := event.Upload.MetaData

		log.Printf("업로드 완료! ID: %s, 메타데이터: %+v\n", uploadID, meta)

		filename := filepath.Join(uploadDir, uploadID)

		loadCmd := exec.Command("docker", "load", "-i", filename)
		output, err := loadCmd.CombinedOutput()
		if err != nil {
			log.Printf("Docker load 실패: %v\n출력: %s", err, output)
			continue
		}

		// 여따가 db로직
	}
}
