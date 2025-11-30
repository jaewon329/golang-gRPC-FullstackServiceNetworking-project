package main

import (
	"log"

	//grpc모듈 import
	"google.golang.org/grpc"
	//insecure 통신채널 생성 위해 import
	"google.golang.org/grpc/credentials/insecure"

	//protoc자동생성 파일 import
	pb "lec-07-prg-04-serverstreaming/proto"
)

// 서버에 stream요청 값을 전달하고 stream을 받아 출력할 recv_message함수 생성
func recv_message(stub pb.ServerStreamingClient) {

}

// 클라이언트를 구동하는 main함수 생성
// hello_grpc때와 거의 비슷하지만 함수를 호출해서 응답 완료
func main() {
	//grpc통신 채널 생성
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("not connect: %v", err)
	}

	//protoc가 생성한 자동생성 파일의 stub함수를 앞서 생성한 채널을 사용하여 실행하여 stub을 생성
	stub := pb.NewServerStreamingClient(conn)
	//만든 stub을 recv_message함수의 입력파라메터로 전달하여 호출
	recv_message(stub)
}
