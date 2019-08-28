package main

import (
	"fmt"
	"go-robots/config"
	"go-robots/rpc/rpcparser"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)


//模拟多台机器
//分别运行3次 开启3台RPC服务
func main()  {

	err:=rpc.Register(rpcparser.Parser{})

	fmt.Println("Register:",err)

	listener,err:=net.Listen("tcp",config.RPCWorker1)
	//listener,err:=net.Listen("tcp",config.RPCWorker2)
	//listener,err:=net.Listen("tcp",config.RPCWorker3)

	if err!=nil{
		fmt.Println("err listen tcp:",err)
	}

	for {
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Printf("conn error:%v",err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}






