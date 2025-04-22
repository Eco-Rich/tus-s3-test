# tus-s3-test

## 로직 구상
1. client -> main server로 업로드 허가 요청
2. Main Server에서 유저 정보 확인 후 해당 요청에 대한 UUID 생성(UUID는 일종의 허가 토큰과 같은 것)
3. TUS 서버로 UUID 전송, TUS 서버는 UUID 보관
4. client의 업로드 허가 요청에 대한 응답으로 UUID 전송
5. client -> TUS 서버로 UUID 및 docker image tar 파일 업로드 요청
6. TUS 서버 내에 요청에 대한 UUID가 존재하는지 확인, 확인 시 docker image tar 파일 ECR에 저장
7. ECR 저장 완료 시 docker image metadata 저장 요청(TUS 서버 -> Main 서버)
8. Main 서버는 DB에 docker image metadata 저장