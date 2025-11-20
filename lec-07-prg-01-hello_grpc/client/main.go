package main

//클라이언트를 위한 패키지 import
import (
	//요청 과정에서의 신호
	"context"
	//로깅
	"log"

	//grpc모듈 import
	"google.golang.org/grpc"
	//insecure 통신채널 생성 위해 import
	"google.golang.org/grpc/credentials/insecure"

	//protoc자동생성 파일 import
	pb "lec-07-prg-01-hello_grpc/proto"
)

func main() {
	//grpc 통신채널 생성
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("not connect: %v", err)
	}

	//protoc가 생성한 자동생성 파일의 stub함수를 앞서 생성한 채널을 사용하여 실행하여 stub을 생성
	stub := pb.NewMyServiceClient(conn)

	//protoc가 생성한 자동생성 파일의 메세지 타입에 맞춰서, 원격 함수에 전달할 메시지를 만들고, 전달할 값을 저장함
	request := &pb.MyNumber{Value: 4}
	//원격 함수를 stub을 사용하여 호출함
	response, err := stub.MyFunction(context.Background(), request)
	if err != nil {
		//에러시 err내용을 출력하고 프로그램 즉시 종료
		log.Fatalf("fail: %v", err)
	}

	//결과를 출력
	log.Println("gRPC result: ", response.Value)
}
