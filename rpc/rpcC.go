package main

import (
	"fmt"
	"go-robots/rpcdemo"
	"net"
	"net/rpc/jsonrpc"
)

func main()  {
	conn,err:=net.Dial("tcp",":1234")

	if err!=nil{
		fmt.Println("xxx",err)
	}

	client:=jsonrpc.NewClient(conn)

	var result float64

	err = client.Call("DivRpc.Div",rpcdemo.DivRpc{10,3},&result)


	fmt.Println("xxx",result,err)

}