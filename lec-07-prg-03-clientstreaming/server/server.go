package main

import (
	"io"
	"log"
	"net"

	//grpc모듈 import
	"google.golang.org/grpc"

	//protoc자동생성 파일 import
	pb "lec-07-prg-03-clientstreaming/proto"
)

// ClientStreamingServer 구현체
type ClientStreamingServer struct {
	//protoc 자동생성 파일을 바탕으로 타입 임베딩
	pb.UnimplementedClientStreamingServer
}

// proto 파일내 정의한 rpc함수 이름에 대응하는 GetServerResponse함수
// 자동생성 파일 clientstreaming_grpc.pb.go의 type ClientStreamingServer interface를 참고한다.
func (s *ClientStreamingServer) GetServerResponse(stream grpc.ClientStreamingServer[pb.Message, pb.Number]) error {
	log.Printf("Server processing gRPC client-streaming")

	//client streaming 횟수를 셀 count변수
	count := 0
	//반복문을 이용해 클라이언트 스트림이 올 때마다 count를 증가시킨다.
	for {
		//Recv()를 통해 스트림을 받아온다.
		_, err := stream.Recv()

		//io.EOF로 클라이언트가 보낼 것을 다 보낸걸 확인하면 count값을 클라이언트에게 보내주고 반복문이 끝난다.
		if err == io.EOF {
			return stream.SendAndClose(&pb.Number{Value: int32(count)})
		}

		//스트림 받다가 문제가 생겼을 때 에러 반환하고 종료
		if err != nil {
			log.Printf("fail: %v", err)
			return err
		}

		//매 스트림을 받을 때마다 count + 1
		count++
	}
}

// 포트 넘버 설정
const portNumber = "[::]:50051"

// main함수 (gpc서버 생성 및 등록/서버 실행 및 유지)
// gRPC서버 등록하고 serve하는 코드이기 때문에 앞서 만든 server의 main.go코드와 거의 똑같다.
func main() {
	//서버에 전달할 tcp리스너 생성
	s, err := net.Listen("tcp", portNumber)
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to listen: %v", err)
	}

	//grpc서버 생성
	grpcServer := grpc.NewServer()
	//위에서 생성한 클라이언트의 정보를 받아 count를 증가시키는 ClientStreamingServer함수를 grpc서버에 등록
	pb.RegisterClientStreamingServer(grpcServer, &ClientStreamingServer{})

	//Serve()를 이용해 grpc서버 실행
	log.Println("Starting server. Listening on port 50051.")
	if err := grpcServer.Serve(s); err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to serve: %v", err)
	}
}
