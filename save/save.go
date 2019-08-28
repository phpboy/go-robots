package save

//go get gopkg.in/olivere/elastic.v5
import (
	"context"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

//数据写入elastic
func Save(item interface{}) error  {
	client, err := elastic.NewClient(elastic.SetURL("http://47.94.169.212:9201"),elastic.SetSniff(false))
	if err != nil {
		fmt.Println("error es",err)
	}
	indexService := client.Index().
		Index("nginx_log").
		Type("nginx_log").
		BodyJson(item)

	re, err := indexService.Do(context.Background())
	if err != nil {
		fmt.Println("error into es",item)
		return err
	}
	log.Println("re:",re)
	return nil
}
