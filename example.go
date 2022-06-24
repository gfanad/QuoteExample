package main

import (
	"context"
	"example/quote_client"
	"log"
	"time"
)

func main() {
	//创建客户端
	opt := quote.Options{
		Target: "ip:8085",
		UserId: 0,
	}

	ctx := context.Background()
	client := quote.NewClient(ctx, opt)

	//获取读取channel
	ch := client.GetReadChannel()

	//创建新goroutine读取推送信息
	go func() {
		for {
			select {
			case quote := <-ch:
				log.Printf("receive: %v\n", quote)
			}
		}
	}()

	//订阅
	id := quote.Id{
		Exchange: 101,
		Code:     "600004",
	}
	client.Sub(id)
	log.Printf("subscirbe: %v\n", id)

	<-time.After(30 * time.Second)

	//取消订阅
	client.UnSub(id)
	log.Printf("unsubscirbe: %v\n", id)

	//订阅非法id不会获得推送信息
	id = quote.Id{
		Exchange: 0,
		Code:     "",
	}
	client.Sub(id)
	log.Printf("subscirbe: %v\n", id)

	<-time.After(30 * time.Second)

	log.Println("finish!")
}
