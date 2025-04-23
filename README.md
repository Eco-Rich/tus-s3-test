# tus-s3-test

## 로직 구상
1. client -> main server로 업로드 허가 요청
2. Main Server에서 유저 정보 확인 후 해당 요청에 대한 UUID 생성(UUID는 일종의 허가 토큰과 같은 것)
3. TUS 서버 On-Demand Instance 띄우기(약 1-2분 소요) 및 TUS 서버로 UUID 전송, TUS 서버는 UUID 보관
4. client의 업로드 허가 요청에 대한 응답으로 UUID 전송
5. client -> TUS 서버로 UUID 및 docker image tar 파일 업로드 요청
6. TUS 서버 내에 요청에 대한 UUID가 존재하는지 확인, 확인 시 docker image tar 파일 ECR에 저장
7. ECR 저장 완료 시 docker image metadata 저장 요청(TUS 서버 -> Main 서버)
8. Main 서버는 DB에 docker image metadata 저장


## 고민해보아야 할 점
1. 3번 항목에서 On-Demand Instance를 띄우면서 client의 응답 대기 시간이 너무 길어지게된다.
   1. UUID 생성과 클라이언트 요청에 대한 응답 이후 인스턴스 시작 명령(시작 시 uuid 전송 동시에)
   2. 고정 ip(TUS 서버)로 health check polling 동시에 클라이언트 측 display

-> 해당 과정에서 display는 어떻게 진행할 것인가?