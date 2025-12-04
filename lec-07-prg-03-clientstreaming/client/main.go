package main

import (
	pb "lec-07-prg-03-clientstreaming/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//보낼 메시지를 만드는 generate_messages함수

// 서버에 보내고 응답을 받아올 send_message함수
func send_message(stub pb.ClientStreamingClient) {

}

// 클라이언트를 구동하는 main함수 생성
// 지금까지 해온 것과 거의 비슷 send_message함수 호출
func main() {
	//grpc통신 채널 생성
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("not connect: %v", err)
	}

	//protoc가 생성한 자동생성 파일의 stub함수를 앞서 생성한 채널을 사용하여 실행하여 stub생성
	stub := pb.NewClientStreamingClient(conn)
	//만든 stub을 send_message함수의 입력파라메터로 전달하여 호출
	send_message(stub)
}
