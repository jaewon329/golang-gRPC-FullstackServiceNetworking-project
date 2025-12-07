package main

import (
	"context"
	pb "lec-07-prg-03-clientstreaming/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 보낼 메시지를 만드는 generate_messages함수
func generate_messages() []string {
	return []string{
		"message #1",
		"message #2",
		"message #3",
		"message #4",
		"message #5",
	}
}

// 서버에 보내고 응답을 받아올 send_message함수
func send_message(stub pb.ClientStreamingClient) {
	//stub을 이용해서 스트림을 보내기 위한 연결을 생성
	responses, err := stub.GetServerResponse(context.Background())
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("fail: %v", err)
	}

	// 스트림으로 보낼 문자열을 generate_messages함수로 받아온다.
	msg := generate_messages()

	// 받아온 것 하나씩 꺼내서 보낸다.
	// 스트림으로 보낼 것이 문자열. 인덱스는 무시하고 문자열 값만으로 반복문 실행
	for _, messages := range msg {
		// 서버로 보낼 것 먼저 출력
		log.Printf("[Client to server] %s\n", messages)

		//messages(보낼 문자열)을 proto 파일에서 정의한 객체 형태로
		request := &pb.Message{Message: messages}
		//서버로 Send를 이용해서 보낸다.
		if err := responses.Send(request); err != nil {
			//보내는 과정에서 에러시 err내용을 출력하고 프로그램 즉시 종료
			log.Fatalf("send fail: %v", err)
		}
	}
	// 다 보내면 다 보냈다고 서버에게 알리고(Close), 서버로 부터 응답을 받아온다(Recv)
	response, err := responses.CloseAndRecv()
	if err != nil {
		// 에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("receive fail: %v", err)
	}

	// 받아온 서버로부터의 결과를 출력
	log.Printf("[Server to client] %d\n", response.Value)
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
