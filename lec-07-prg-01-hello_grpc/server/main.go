package main

import (
	//서버를 위한 패키지 import
	//요청 과정에서의 신호
	"context"
	//로깅
	"log"
	//네트워크 기능 도구
	"net"

	//grpc모듈 import
	"google.golang.org/grpc"

	//원격 호출될 함수 import
	"lec-07-prg-01-hello_grpc/server/function"
	//protoc자동생성 파일 import
	pb "lec-07-prg-01-hello_grpc/proto"
)

// MyServiceServer구현체
type MyServiceServer struct {
	//protoc 자동생성 파일을 바탕으로 타입 임베딩
	pb.UnimplementedMyServiceServer
}

// proto 파일내 정의한 rpc함수 이름에 대응하는 MyFunction함수
func (s *MyServiceServer) MyFunction(ctx context.Context, MyNumber *pb.MyNumber) (*pb.MyNumber, error) {
	//원격 호출할 함수에게 client로 부터 받은 입력 파라메터를 전달하고 결과를 가져옴
	request := MyNumber.GetValue()
	response := function.My_func(int(request))
	//원격 함수 호출 결과를 client에게 돌려줌
	return &pb.MyNumber{Value: int32(response)}, nil
}

// 포트 넘버 설정
const portNumber = "[::]:50051"

// main함수 (grpc서버 생성 및 등록/ 서버 실행 및 유지)
func main() {
	//서버에 전달할 tcp리스너 생성 : 파이썬에서는 add_insecure_port()를 호출하면 내부적으로 자동 처리 했었음
	s, err := net.Listen("tcp", portNumber)
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to listen: %v", err)
	}

	//grpc서버 생성
	grpcServer := grpc.NewServer()
	//위에서 생성한 클라이언트로 부터 값을 받아 My_func를 실행해 결과를 돌려주는 MyServiceServer함수를 grpc서버에 등록
	pb.RegisterMyServiceServer(grpcServer, &MyServiceServer{})
	//grpc서버 실행
	//Serve()를 사용해서 앞서 생성한 리스너를 grpc서버로 전달해줌, 또한 Serve()에서 동시성 처리
	//thread pool을 지정해 줬던 파이썬과 달리 go에서는 고루틴으로 연결마다 내부적으로 새 고루틴 생성해 동시성 처리를 한다
	log.Println("Starting server. Listening on port 50051.")
	if err := grpcServer.Serve(s); err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("failed to server: %v", err)
	}
}
