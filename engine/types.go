package engine

import (
	"fmt"
)

//model
type Model struct {
	Name string
	Url string
	Size string
	Sum string
	Label string
	Time string
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
type Engine struct {
	Scheduler Schedulers
	WorkerCount int
}

type Schedulers interface {
	ConfigureMaterWorkerChan(chan Request)
	Submit(Request)
}


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


//多线程的爬虫
func (e * Engine) Run(seeds... Request)  {

	in := make(chan Request)
	out := make(chan ParseResult)

	e.Scheduler.ConfigureMaterWorkerChan(in)//e.Scheduler.workerChan = in

	for i:=0;i<e.WorkerCount;i++{
		CreateWorker(in,out)//开几个go-routine for循环接受in和输出out
	}
	for _,seed:= range seeds{
		e.Scheduler.Submit(seed)//e.Scheduler.workerChan <- seed（Request）
	}

	count := 0
	for {
		result := <- out
		for _,Item := range result.Items{
			count++
			//Save(Item)
			fmt.Printf("#%d-Item %v \n",count,Item)
		}
		for _,request := range result.Requests{
			e.Scheduler.Submit(request)
		}
	}
}