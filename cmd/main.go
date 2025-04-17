package main

import (
	"gin/api/route"
	"gin/config"
	"log"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatal(".env 파일을 로드하지 못했습니다")
	}

	router := route.SetupRouter()
	router.Run(":8080")
}
