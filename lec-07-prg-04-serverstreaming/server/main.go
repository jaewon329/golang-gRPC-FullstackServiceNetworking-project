package main

import (
	//protoc자동생성 파일 import
	pb "lec-07-prg-04-serverstreaming/proto"
)

// ServerStreamingServer구현체
type ServerStreamingServer struct {
	//protoc 자동생성 파일을 바탕으로 타입 임베딩
	pb.UnimplementedServerStreamingServer
}
