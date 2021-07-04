package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

// ESDemo:Go操作Elasticsearch客户端

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func main() {
	//创建一个客户端
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("connect to es success")
	p1 := Person{Name: "xili", Age: 23, Married: false}
	//插入数据
	put1, err := client.Index().
		Index("user").
		BodyJson(p1).  //Json化数据
		Do(context.Background()) //执行插入操作
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
