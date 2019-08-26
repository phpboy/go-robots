package main

import (
	"go-robots/engine"
	"go-robots/parser"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)


func main(){
	e:=engine.Engine{
		Scheduler:&engine.SimpleSchedulers{},
		WorkerCount:20,
	}
	e.Run(
		engine.Request{
			Url:"http://xiazai.zol.com.cn/word_index.html",
			ParserFunc:parser.ListParser,
		},
	)
}




func CreateRpcPools(hosts [] string) chan *rpc.Client{
	var clients [] *rpc.Client
	for _,host:=range hosts{

		conn,err:=net.Dial("tcp",host)
		if err!=nil{
			log.Printf("err dial tcp %v",host)
			continue
		}

		clients = append(clients,jsonrpc.NewClient(conn))
	}
	out:= make(chan *rpc.Client)

	go func() {
		for{
			for _,c:=range clients{
				out<-c
			}
		}
	}()

	return out
}

