package tus

import (
	"github.com/tus/tusd/pkg/handler"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
 * 순서
 * 1. tus 업로드
 * 2. 업로드 완료 시 Docker load를 통한 validation check
 * 3. 체크 완료 시 이미지 이름 추출
 * 4. ECR 정보 추출 및 이미지 태깅
 * 5. ECR Login 및 Docker Push
 * 6. 업로드 된 tar 파일 삭제
 * 7. DB에 해당 버전 정보 저장
 */
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
		log.Printf("Docker load 성공:\n%s", output)

		imageName := extractImageName(string(output))
		if imageName == "" {
			log.Println("이미지 이름 추출 실패")
			continue
		}

		accountID := os.Getenv("AWS_ACCOUNT_ID")
		region := os.Getenv("AWS_REGION")
		repoName := os.Getenv("ECR_REPO")
		ecrTag := accountID + ".dkr.ecr." + region + ".amazonaws.com/" + repoName + ":latest"

		exec.Command("docker", "tag", imageName, ecrTag).Run()

		loginCmd := exec.Command("sh", "-c", "aws ecr get-login-password --region "+region+" | docker login --username AWS --password-stdin "+accountID+".dkr.ecr."+region+".amazonaws.com")
		if loginOut, err := loginCmd.CombinedOutput(); err != nil {
			log.Printf("ECR 로그인 실패: %v\n%s", err, loginOut)
			continue
		}

		pushCmd := exec.Command("docker", "push", ecrTag)
		if pushOut, err := pushCmd.CombinedOutput(); err != nil {
			log.Printf("Docker push 실패: %v\n%s", err, pushOut)
			continue
		}
		log.Println("ECR 업로드 성공:", ecrTag)

		os.Remove(filename)

		// 여따가 db로직 or Main 서버로 넘겨서 처리

	}
}

func extractImageName(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Loaded image:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Loaded image:"))
		}
	}
	return ""
}
