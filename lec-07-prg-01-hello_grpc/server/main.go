package main

import (
	//grpc모듈 import
	_ "google.golang.org/grpc"

	//원격 호출될 함수 import
	_ "lec-07-prg-01-hello_grpc/server/function"
	//protoc자동생성 파일 import
	_ "lec-07-prg-01-hello_grpc/proto"
)

//MyServiceServer구현체

//proto 파일내 정의한 rpc함수 이름에 대응하는 MyFunction함수

//main함수 (grpc서버 생성 및 등록/ 서버 실행 및 유지)
