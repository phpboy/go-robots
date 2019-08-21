package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"regexp"
)

//用户信息
type GoodInfo struct {
	price string
}
//请求
type Request struct {
	Url string
	ParserFunc func([]byte) ParseResult
}
//结果
type ParseResult struct {
	Requests []Request
	Items []interface{}
}

//爬虫引擎
type engine struct {
	Scheduler Schedulers
	WorkerCount int
}

type Schedulers interface {
	ConfigureMaterWorkerChan(chan Request)
	Submit(Request)
}


//SimpleSchedulers <<<
type SimpleSchedulers struct {
	workChan chan Request
}
func (s *SimpleSchedulers) ConfigureMaterWorkerChan(c chan Request)  {
	s.workChan = c
}

func (s *SimpleSchedulers) Submit(r Request){
	go func() {
		s.workChan <- r
	}()
}
//SimpleSchedulers >>>


func main(){
	e:=engine{
		Scheduler:&SimpleSchedulers{},
		WorkerCount:20,
	}
	e.Run(
		Request{
			Url:"https://www.dressbycouturier.com",
			ParserFunc:CategoryListParser,
		},
	)
}
//多线程的爬虫
func (e *engine) Run(seeds... Request)  {

	in := make(chan Request)
	out := make(chan ParseResult)

	e.Scheduler.ConfigureMaterWorkerChan(in)//e.Scheduler.workerChan = in

	for i:=0;i<e.WorkerCount;i++{
		createWorker(in,out)//开几个go-routine for循环接受in和输出out
	}
	for _,seed:= range seeds{
		e.Scheduler.Submit(seed)//e.Scheduler.workerChan <- seed（Request）
	}

	count := 0
	for {
		result := <- out
		for _,Item := range result.Items{
			count++
			fmt.Printf("#%d-Item %v \n",count,Item)
		}
		for _,request := range result.Requests{
			e.Scheduler.Submit(request)
		}
	}

}

func createWorker(in chan Request,out chan ParseResult)  {
	go func() {
		for{
			request := <- in
			result,err := worker(request)
			if err!=nil{
				continue
			}
			out <- result
		}
	}()
}

func worker(r Request) (ParseResult ,error) {

	body,err := fetch(r.Url)

	if err != nil{
		fmt.Println("err:",err)
		return ParseResult{},err
	}

	return r.ParserFunc(body),nil

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

//获取首页分类
func CategoryListParser(body []byte) ParseResult {

	var exp = `<a class="nav_a" style="[^>]+" href="(https://www.dressbycouturier.com/[^>]+)">([^<]+)</a>`
	regexpResult := regexp.MustCompile(exp)

	matches := regexpResult.FindAllSubmatch(body,-1)

	var result = ParseResult{}
	for _,m:= range matches{
		var url = string(m[1])
		result.Items = append(result.Items,"category name:"+string(m[2]))
		result.Requests = append(result.Requests,Request{
			Url:url,
			ParserFunc: GoodParser,
		})
	}
	return result
}

//商品列表解析器
func GoodParser(body []byte) ParseResult {
	var exp = `<a goods_id="[0-9]+" href="(https://www.dressbycouturier.com/[^>]+)" target="_blank" title="[^>]+">[^<]+</a>`
	re := regexp.MustCompile(exp)
	matches := re.FindAllSubmatch(body,-1)
	result := ParseResult{}
	for _,m:= range matches{
		var url = string(m[1])
		result.Items = append(result.Items,"good name:"+string(m[1]))
		result.Requests = append(result.Requests,Request{
			Url:url,
			ParserFunc: DetailParser,
		})
	}

	return result
}

//商品详情解析器
func DetailParser(body []byte) ParseResult  {
	var exp = `<span class="price_local" style="color:#c00000">([^>]+)</span>`
	re := regexp.MustCompile(exp)
	matches := re.FindAllSubmatch(body,-1)

	var good = GoodInfo{}
	for _,m:= range matches{
		good.price = string(m[1])
	}

	var result = ParseResult{
		Items:[]interface{}{good},
	}

	return result
}


//获取远程url
func fetch(url string) ([]byte,error)  {
	resp,err := http.Get(url)

	if err != nil{
		return []byte{},err
	}
	if resp.StatusCode != http.StatusOK{
		fmt.Println("get url error CODE",resp.StatusCode)
		return []byte{},err
	}

	body,err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		return body,err
	}
	return body,nil
}


