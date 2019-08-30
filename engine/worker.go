package engine

import (
	"fmt"
	"go-robots/config"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func CreateWorker(in chan Request,out chan ParseResult,client chan *rpc.Client)  {
	go func() {
		for{
			request := <- in
			result,err := worker(request,client)
			if err!=nil{
				continue
			}
			out <- result
		}
	}()
}

func worker(r Request,c chan *rpc.Client) (ParseResult ,error) {

	client:=<-c

	var result ParseResult

	err := client.Call("Parser.ParseFunc",Parser{Url: r.Url, Method: r.ParserFunc},&result)

	if err!=nil{
		fmt.Println("rpc result:",result,err)
	}

	return result,nil

}

//创建rpc池子
func CreateRpcPool() chan *rpc.Client {

	var clients  []*rpc.Client

	var workers = map[string]string{"1":config.RPCWorker1,"2":config.RPCWorker2,"3":config.RPCWorker3}

	for _,w:=range workers{
		conn,err:=net.Dial("tcp",w)
		if err!=nil{
			fmt.Println("Dial tcp error",err)
			continue
		}
		clients = append(clients,jsonrpc.NewClient(conn))
	}

	out := make(chan *rpc.Client)
	go func() {
		for {
			for _,c:=range clients{
				out<- c
			}
		}
	}()
	return out
}
