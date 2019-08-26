package engine

import (
	"fmt"
	"go-robots/fetch"
)

func CreateWorker(in chan Request,out chan ParseResult)  {
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

	body,err := fetch.Get(r.Url)

	if err != nil{
		fmt.Println("err:",err)
		return ParseResult{},err
	}

	return r.ParserFunc(body),nil

}
