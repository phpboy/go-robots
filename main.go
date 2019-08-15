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

func main(){
	var result = CityListParser("https://www.zhenai.com/zhenghun")
	for i,mm:= range result.Items{
		fmt.Printf("%s\n",mm)
		fmt.Printf("%s\n",result.Requests[i].Url)
	}
}

func CityListParser(url string) ParseResult {
	resp,err := http.Get(url)
	if err != nil{
		panic(err)
	}

	//defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		fmt.Println("get url error CODE",resp.StatusCode)
		panic(err)
	}

	all,err := ioutil.ReadAll(resp.Body)

	if err!=nil{
		panic(err)
	}

	re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)

	matches := re.FindAllSubmatch(all,-1)

	var result = ParseResult{}
	//循环获取每个城市的第一页的用户
	for _,m:= range matches{

		var url = string(m[1])

		fmt.Printf("city:%s url %s\n\n",m[2],url)

		result.Items = append(result.Items,"City:"+string(m[2]))
		result.Requests = append(result.Requests,Request{
			Url:url,
			ParserFunc: func(bytes []byte) ParseResult {
				return CityParser(url)
			},
		})
	}
	return result
}

//城市解析器 -解析第一页的用户
func CityParser(url string) ParseResult {
	cityResp,err := http.Get(url)
	if err != nil{
		panic(err)
	}
	cityBody,err := ioutil.ReadAll(cityResp.Body)
	if err!=nil{
		panic(err)
	}

	re := regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	matches := re.FindAllSubmatch(cityBody,-1)
	profile := ParseResult{}
	for _,m:= range matches{
		var url = string(m[1])
		profile.Items = append(profile.Items,"User:"+string(m[2]))
		profile.Requests = append(profile.Requests,Request{
			Url:url,
			ParserFunc: func(bytes []byte) ParseResult {
				return UserParser(url)
			},
		})
	}
	return profile
}

//用户解析器 获取用户的个人信息 目前因为对方反爬虫 返回403页面
func UserParser(url string) ParseResult  {

	userResp,err := http.Get(url)
	if err != nil{
		panic(err)
	}
	userBody,err := ioutil.ReadAll(userResp.Body)
	if err!=nil{
		panic(err)
	}

	var userinfo = Profile{}

	userinfo.Name = string(userBody)

	var result = ParseResult{
		Items:[]interface{}{userinfo},
	}

	return result
}


