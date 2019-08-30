package main

import (
	"flag"
	"fmt"
	"go-robots/rpc/rpcparser"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)


//模拟多台机器
//分别运行3次 开启3台RPC服务

var port = flag.String("port","","端口必填")

func main()  {

	flag.Parse()
	if *port == ""{
		fmt.Println("端口必填")
		return
	}

	err:=rpc.Register(rpcparser.Parser{})

	fmt.Println("Register:",err)

	listener,err:=net.Listen("tcp",*port)

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






