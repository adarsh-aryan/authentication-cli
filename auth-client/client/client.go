package client

import "net/rpc"

var Client *rpc.Client

type RPCClient struct {
	client *rpc.Client
}

func NewRpcClient(client *rpc.Client) {
	Client = client
}
