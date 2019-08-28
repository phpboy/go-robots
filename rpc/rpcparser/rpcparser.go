package rpcparser

import (
	"fmt"
	"go-robots/config"
	"go-robots/engine"
	"go-robots/fetch"
	"go-robots/parser"
)

type Parser struct {
	Url string
	Method string
}

func (Parser) ParseFunc(args Parser,result *engine.ParseResult) error {

	body,err := fetch.Get(args.Url)

	if err != nil{
		fmt.Println("Err Fetch url:",err)
		return err
	}

	switch args.Method {
	case config.ListParserConfig:
		*result= parser.ListParser(body)
	case config.SoftParserConfig:
		*result= parser.SoftParser(body)
	case config.DetailParserConfig:
		*result= parser.DetailParser(body,"",args.Url)
	}
	return nil

}