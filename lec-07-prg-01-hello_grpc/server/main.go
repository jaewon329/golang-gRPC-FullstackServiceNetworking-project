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
	_ "google.golang.org/grpc"

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
}
