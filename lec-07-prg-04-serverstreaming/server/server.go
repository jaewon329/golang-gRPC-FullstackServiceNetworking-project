package main

import (
	"fmt"
	"log"
	"net"

	//grpc모듈 import
	"google.golang.org/grpc"

	//protoc자동생성 파일 import
	pb "lec-07-prg-04-serverstreaming/proto"
)

// ServerStreamingServer구현체
type ServerStreamingServer struct {
	//protoc 자동생성 파일을 바탕으로 타입 임베딩
	pb.UnimplementedServerStreamingServer
}

// proto 파일내 정의한 rpc함수 이름에 대응하는 GetServerResponse함수
// 자동생성 파일 serverstreaming_grpc.pb.go의 type ServerStreamingServer interface를 참고한다.
func (s *ServerStreamingServer) GetServerResponse(request *pb.Number, stream grpc.ServerStreamingServer[pb.Message]) error {
	//클라이언트에서 결정한 스트리밍의 수를 가져온다.
	req_Count := request.GetValue()

	//해야할 스트리밍만큼 반복문을 이용해 스트리밍
	for i := int32(1); i <= req_Count; i++ {
		//각 스트리밍마다(반복마다) 내보낼 문자열
		make_message := fmt.Sprintf("message #%d", i)

		//응답 메시지 생성
		response := &pb.Message{
			Message: make_message,
		}

		//스트리밍으로 전송
		//파이썬에서는 yield를 사용했지만 Go에서는 stream.Sent()를 사용
		//yield에서는 message 리스트를 먼저 만들고 하나씩 기다렸다가 보내는 방식
		//Send는 바로 보내고, 보낸거 받으면 다음거 보내고 err도 바로 확인
		if err := stream.Send(response); err != nil {
			//에러시 err내용을 출력한다.
			log.Printf("failed to send: %v", err)
			return err
		}
	}
	return nil
}

// 포트 넘버 설정
const portNumber = "[::]:50051"

// main함수 (grpc서버 생성 및 등록/서버 실행 및 유지)
// hello_grpc 때랑 거의 유사하다.
func main() {
	//서버에 전달할 tcp리스너 생성 : 파이썬에서는 add_insecure_port()를 호출하면 내부적으로 자동 처리 했었음.
	s, err := net.Listen("tcp", portNumber)
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to listen: %v", err)
	}

	//grpc서버 생성
	grpcServer := grpc.NewServer()
	//위에서 생성한 Streaming을 진행하는 ServerStreamingServer함수를 grpc서버에 등록
	pb.RegisterServerStreamingServer(grpcServer, &ServerStreamingServer{})
	//grpc서버 실행
	//Serve()를 사용해서 앞서 생성한 리스너를 grpc서버로 전달해줌, 또한 Serve()에서 동시성 처리
	//thread pool을 지정해 줬던 파이썬과 달리 go에서는 고루틴으로 연결마다 내부적으로 새 고루틴 생성해 동시성 처리를 한다
	log.Println("Starting server. Listening on port 50051.")
	if err := grpcServer.Serve(s); err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to server: %v", err)
	}
}
