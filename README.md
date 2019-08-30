# go-robots

GO语言并发分布式爬虫


1 生成多台RPC服务器 :

 go run rpc/rpcserver.go --port=:9001 &

 go run rpc/rpcserver.go --port=:9002 &

 go run rpc/rpcserver.go --port=:9003 &


2 然后运行爬虫 go run main.go

3 查看爬取数据 先运行 go run frontend/main.go 
    浏览器查看  127.0.0.1:8888/search

代码结构

go-robots

    engine   -- 爬虫引擎 派生多个线程goroutine 生成多个worker去执行 
    
    fetch    -- 从远程url获取html源码 
    
    frontend -- 从数据库读取获取的爬虫内容显示
    
    parser   -- 爬虫html源码分析 并获取爬取字段
    
    rpc      -- 本地模拟多机器rpc 分布式处理爬虫任务
        
    save     -- 保存数据到elasticsearch
    
    main.go  -- 代码入口
