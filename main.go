package main

import (
	"go-robots/config"
	"go-robots/engine"
)

//ZOL软件下载为爬取目标 http://xiazai.zol.com.cn/word_index.html
//先爬取分类列表
//再爬取分类下软件列表
//最后爬取软件详情信息
//数据录入 elasticsearch

func main(){
	e:=engine.Engine{
		Scheduler:&engine.SimpleSchedulers{},//调度器
		WorkerCount:20,//生成20个协程
	}
	e.Run(
		engine.Request{
			Url:        config.URL,//ZOL软件下载为爬去目标
			ParserFunc: config.ListParserConfig,
		},
	)
}

