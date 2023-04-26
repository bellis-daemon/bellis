package v2api

import (
	"errors"
	"google.golang.org/grpc"
)

var connMap = make(map[string]*grpc.ClientConn)

func getGrpcConn(host string) (*grpc.ClientConn, error) {
	if c, exist := connMap[host]; !exist {
		cmdConn, err := grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			return &grpc.ClientConn{}, errors.New("cant connect grpc server")
		}
		connMap[host] = cmdConn
		return cmdConn, nil
	} else {
		return c, nil
	}
}
