package parser

import (
	"go-robots/engine"
	"regexp"
)

//获取首页分类
func ListParser(body []byte) engine.ParseResult {
	var exp = `<a href="/([^>]+)" [^>]+>[^<]+</a>`
	matches := regexp.MustCompile(exp).FindAllSubmatch(body,-1)
	count:=0
	var result = engine.ParseResult{}
	for _,m:= range matches{
		count++
		if count>5{
			break
		}
		var url = "http://xiazai.zol.com.cn/"+string(m[1])
		result.Requests = append(result.Requests,engine.Request{
			Url:url,
			ParserFunc: SoftParser,
		})
	}
	return result
}

//列表解析器
func SoftParser(body []byte) engine.ParseResult {
	var exp = `<a target="_blank" href="/([^>]+)">([^<]+)</a>`
	matches := regexp.MustCompile(exp).FindAllSubmatch(body,-1)
	result := engine.ParseResult{}
	for _,m:= range matches{
		var url = "http://xiazai.zol.com.cn/"+string(m[1])
		var name = m[2]
		result.Requests = append(result.Requests,engine.Request{
			Url:url,
			ParserFunc: func(bytes []byte) engine.ParseResult {
				result := DetailParser(bytes,string(name),url)
				return result
			},
		})
	}

	return result
}

//商品详情解析器
func DetailParser(body []byte,name,url string) engine.ParseResult  {

	var expSize = `<li><span>资源大小：</span>([^<]+)</li>`
	var expSum = `<li><span>月下载量：</span>([^<]+)</li>`
	var expLabel = `<li><span>软件属性：</span>([^<]+)</li>`
	var expTime = `<li><span>更新时间：</span>([^<]+)</li>`

	size := regexp.MustCompile(expSize).FindAllSubmatch(body,-1)
	sum := regexp.MustCompile(expSum).FindAllSubmatch(body,-1)
	label := regexp.MustCompile(expLabel).FindAllSubmatch(body,-1)
	time := regexp.MustCompile(expTime).FindAllSubmatch(body,-1)

	var model = engine.Model{}

	model.Name = name
	model.Url = url
	if len(size)>=1{
		model.Size = string(size[0][1])
		model.Sum = string(sum[0][1])
		model.Label = string(label[0][1])
		model.Time = string(time[0][1])
	}
	var result = engine.ParseResult{
		Items:[]interface{}{model},
	}

	return result
}

