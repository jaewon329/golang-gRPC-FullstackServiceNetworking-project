package main

import (
	"context"
	"io"
	"log"

	//protoc자동생성 파일 import
	pb "lec-07-prg-02-bidirectional-streaming/proto"

	//grpc모듈 import
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

// 서버에 보내고 응답을 받아올 send_messages()함수
// 고루틴을 이용해 별도의 실행흐름으로 서버에서 응답을 받는 것을 기다리고 받으면 출력해야한다.
func send_message(stub pb.BidirectionalClient) {
	//stub을 이용해서 스트림을 보내기 위한 연결을 생성
	responses, err := stub.GetServerResponse(context.Background())
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("fail: %v", err)
	}

	// 스트림으로 보낼 문자열을 generate_messages함수로 받아온다.
	msg := generate_messages()

	// 채널 생성 : go func()의 작업(서버로 부터 응답을 받아와 출력하는)이
	// send_message()함수가 끝났다고 바로 종료되는 것을 막는다.
	wait_to_receive := make(chan struct{})

	//계속 서버의 응답을 기다리면서 Recv를 통해 응답을 받아와 출력하는 고루틴
	go func() {
		//무한루프
		for {
			//스트림 연결에서 Recv를 통해 응답을 받아온다
			response, err := responses.Recv()
			//만약 서버가 다 보냈음을 확인하면 채널을 종료한다.
			// -> send_message 최종 종료 시점 무한 루프를 종료하고 go func() 작업이 끝난다.
			if err == io.EOF {
				close(wait_to_receive)
				return
			}

			//받는데서 다른 에러 발생시 err출력하고 즉시 종료
			if err != nil {
				log.Fatalf("failed to receive: %v", err)
			}

			//받아온 응답 출력
			log.Printf("[server to client] %s", response.Message)
		}
	}()

	// generate_message에서 받아온 것 하나씩 꺼내서 보낸다.
	// 스트림으로 보낼 것이 문자열. 인덱스는 무시하고 문자열 값만으로 반복문 실행
	for _, messages := range msg {
		log.Printf("[Client to server] %s\n", messages)

		//messages(보낼 문자열)을 proto 파일에서 정의한 객체 형태로
		request := &pb.Message{Message: messages}
		//서버로 Send를 이용해서 보낸다.
		if err := responses.Send(request); err != nil {
			//보내는 과정에서 에러시 err내용을 출력하고 프로그램 즉시 종료
			log.Fatalf("send fail: %v", err)
		}
	}
	//다 보내면 다 보냈다고 알리고 닫는다.
	responses.CloseSend()

	//go func()가 종료되지 않는 한 send_message() 함수를 종료하지 않고 유지시킨다.
	<-wait_to_receive
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
	stub := pb.NewBidirectionalClient(conn)
	//만든 stub을 send_message함수의 입력파라메터로 전달하여 호출
	send_message(stub)
}
