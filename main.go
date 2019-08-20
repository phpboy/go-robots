package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

//用户信息
type Profile struct {
	Name string
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
			Url:"https://www.zhenai.com/zhenghun",
			ParserFunc:CityListParser,
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

	for {
		result := <- out
		for _,Item := range result.Items{
			fmt.Printf("Item %v \n",Item)
		}
		for _,request := range result.Requests{
			e.Scheduler.Submit(request)
		}
	}

}


//单线程的爬虫
func (e engine) Run1(seeds... Request)  {

	var requests  [] Request

	for _,seed:=range seeds{
		requests = append(requests,seed)
	}

	for len(requests)>0{
		r := requests[0]
		requests = requests[1:]

		body,err:=fetch(r.Url)
		if err!=nil{
			fmt.Println("err:",err)
			continue
		}
		parseResult := r.ParserFunc(body)

		requests = append(requests,parseResult.Requests...)

		for _,item:=range parseResult.Items{
			fmt.Printf("get %s\n",item)
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




//循环获取每个城市的第一页的用户
func CityListParser(body []byte) ParseResult {

	regexpResult := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)

	matches := regexpResult.FindAllSubmatch(body,-1)

	var result = ParseResult{}
	var count = 0
	for _,m:= range matches{
		count++
		var url = string(m[1])
		result.Items = append(result.Items,"City:"+string(m[2]))
		result.Requests = append(result.Requests,Request{
			Url:url,
			ParserFunc: CityParser,
		})
		if count>2{
			//break
		}
	}
	return result
}

//城市解析器 -解析第一页的用户
func CityParser(body []byte) ParseResult {
	re := regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	matches := re.FindAllSubmatch(body,-1)
	result := ParseResult{}
	var count = 0
	for _,m:= range matches{
		count++
		var url = string(m[1])
		result.Items = append(result.Items,"User:"+string(m[2]))
		result.Requests = append(result.Requests,Request{
			Url:url,
			ParserFunc: UserParser,
		})
		if count>10{
			break
		}
	}

	return result
}

//用户解析器 获取用户的个人信息 目前因为对方反爬虫 返回403页面
func UserParser(body []byte) ParseResult  {
	var userInfo = Profile{}

	userInfo.Name = "age:18;name:xx"

	var result = ParseResult{
		Items:[]interface{}{userInfo},
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


