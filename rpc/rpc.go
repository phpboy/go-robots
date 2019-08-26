package main

import (
	"fmt"
	"go-robots/rpcdemo"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main()  {
	RpcServer()
}

func RpcServer(){

	rpc.Register(rpcdemo.DivRpc{})

	listener,err:=net.Listen("tcp",":1234")

	if err!=nil{
		fmt.Println("err:",err)
	}
	fmt.Println("aaaaa:",listener)
	for {
		conn,err2:=listener.Accept()
		if err2!=nil{
			fmt.Printf("con error:%v",err2)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}

}