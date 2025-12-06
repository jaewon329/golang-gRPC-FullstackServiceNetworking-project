package main

import (
	"io"
	"log"

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

// main함수 (grpc서버 생성 및 등록/서버 실행 및 유지)
//앞에서 계속 만들었던 것과 거의 유사하다.
