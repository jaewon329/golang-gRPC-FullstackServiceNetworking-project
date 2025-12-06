package main

import (
	"io"
	"log"
	"net"

	//protoc자동생성 파일 import
	pb "lec-07-prg-02-bidirectional-streaming/proto"

	//grpc모듈 import
	"google.golang.org/grpc"
)

// BidirectionalService구현체
type BidirectionalService struct {
	//protoc자동생성 파일을 바탕으로 타입 임베딩
	pb.UnimplementedBidirectionalServer
}

// proto 파일내 정의한 rpc함수 이름에 대응하는 GetServerResponse함수
// 자동생성 파일 bidirectional_grpc.pb.go의 type BidirectionalServer interface를 참고한다.
func (s *BidirectionalService) GetServerResponse(messages grpc.BidiStreamingServer[pb.Message, pb.Message]) error {
	log.Printf("Server processing gRPC bidirectional streaming.")

	//클라이언트에서 보낸 것을 받아올 무한루프
	for {
		//클라이언트에서 보내는 것을 Recv를 이용해 스트림을 하나씩 받아온다.
		message, err := messages.Recv()
		//만약 클라이언트에서 스트림을 다 보낸 것을 확인하면 반복문 빠져나옴
		if err == io.EOF {
			return nil
		}

		//받아오는데서 에러시 erer내용을 출력하고 종료
		if err != nil {
			log.Printf("faild to receive: %v", err)
			return err
		}

		//받아온 message를 Send를 이용해 스트리밍으로 전송한다.
		//보내는데서 에러시 err내용을 출력하고 종료
		if err := messages.Send(message); err != nil {
			log.Printf("failed to send: %v", err)
			return err
		}
	}
}

// 포트 넘버 설정
const portNumber = "[::]:50051"

// main함수 (grpc서버 생성 및 등록/서버 실행 및 유지)
// 앞에서 계속 만들었던 것과 거의 유사하다.
func main() {
	//서버에 전달할 tcp리스너 생성 : 피어썬에서는 add_insecure_port()를 호출하면 내부적으로 자동 처리 했었음.
	s, err := net.Listen("tcp", portNumber)
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to listen: %v", portNumber)
	}

	//grpc서버 생성
	grpcServer := grpc.NewServer()
	//위에서 생성한 Streaming을 진행하는 BiderectionalService함수를 grpc서버에 등록
	pb.RegisterBidirectionalServer(grpcServer, &BidirectionalService{})
	//grpc서버 실행
	//Serve()를 사용해서 앞서 생성한 리스너를 grpc서버로 전달해줌, 또한 Serve()에서 동시성 처리
	//thread pool을 지정해 줬던 파이썬과 달리 go에서는 고루틴으로 연결마다 내부적으로 새 고루틴 생성해 동시성 처리를 한다
	log.Println("Starting server. Listening on port 50051.")
	if err := grpcServer.Serve(s); err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to server: %v", err)
	}
}
