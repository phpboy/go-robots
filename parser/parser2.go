package parser

import (
	"fmt"
	"go-robots/config"
	"go-robots/engine"
	"regexp"
)



func ListParser2(body []byte) engine.ParseResult{
	var result = engine.ParseResult{}
	for i:=1;i<1000;i++{
		url := fmt.Sprintf("https://studygolang.com/articles/%d",i)
		result.Requests = append(result.Requests,engine.Request{
			Url:        string(url),
			ParserFunc: config.DetailParserConfig,
		})
	}
	return result
}

//商品详情解析器
func DetailParser2(body []byte,url string) engine.ParseResult  {

	var expSize = `<h1 id="title" data-id="[^>]+">([^<]+)</h1>`

	name := regexp.MustCompile(expSize).FindAllSubmatch(body,-1)

	var model = engine.Model{}
	model.Url = url

	if len(name)>=1{
		model.Name = string(name[0][1])
	}

	return engine.ParseResult{Items:[]interface{}{model}}
}


