package quote

import (
	"context"
	pusherrpc "example/quote_client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
)

const (
	BUFFER_SIZE = 100
)

type Id struct {
	Exchange uint32 `json:"exchange"`
	Code     string `json:"code"`
}

type Options struct {
	Target     string
	UserId     uint32
	BufferSize int
}

type Quote struct {
	Id   Id
	Data string
}

type Client struct {
	opt          Options
	userId       uint32
	pusherClient pusherrpc.PusherClient
	subChan      chan Id
	unSubChan    chan Id
	readBuffer   chan Quote
	ctx          context.Context
	startOnce    sync.Once
}

func NewClient(ctx context.Context, opt Options) *Client {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()}

	conn, err := grpc.Dial(opt.Target, opts...)
	if err != nil {
		panic(err)
	}

	pusherClient := pusherrpc.NewPusherClient(conn)
	log.Println("client start!")

	c := &Client{
		ctx:          ctx,
		opt:          opt,
		userId:       opt.UserId,
		pusherClient: pusherClient,
		subChan:      make(chan Id, 10),
		unSubChan:    make(chan Id, 10),
		readBuffer:   make(chan Quote, BUFFER_SIZE),
	}

	return c
}

func (c *Client) start() {
	stream, err := c.pusherClient.RealTimeQuote(c.ctx)
	if err != nil {
		panic(err)
		return
	}

	err = stream.Send(&pusherrpc.RealTimeQuoteRequest{
		Op:     pusherrpc.Op_INIT,
		UserId: c.userId,
		Ids:    nil,
	})
	if err != nil {
		panic(err)
	}

	go c.handleRequest(stream)

	for {
		msg, err := stream.Recv()
		if err != nil {
			//log.Println(err)
			return
		}
		for _, quote := range msg.Quote {
			c.readBuffer <- Quote{
				Id: Id{
					Exchange: quote.Id.Exchange,
					Code:     quote.Id.Code,
				},
				Data: quote.Data,
			}
		}
	}
}

func (c *Client) handleRequest(stream pusherrpc.Pusher_RealTimeQuoteClient) {
	for {
		select {
		case topic := <-c.subChan:
			len := len(c.subChan)

			var ids []*pusherrpc.Id

			ids = append(ids, &pusherrpc.Id{
				Exchange: topic.Exchange,
				Code:     topic.Code,
			})

			for i := 0; i < len; i++ {
				topic := <-c.subChan
				ids = append(ids, &pusherrpc.Id{
					Exchange: topic.Exchange,
					Code:     topic.Code,
				})
			}

			err := stream.Send(&pusherrpc.RealTimeQuoteRequest{
				Op:     pusherrpc.Op_SUB,
				UserId: c.userId,
				Ids:    ids,
			})
			if err != nil {
				panic(err)
			}

		case topic := <-c.unSubChan:

			var ids []*pusherrpc.Id
			ids = append(ids, &pusherrpc.Id{
				Exchange: topic.Exchange,
				Code:     topic.Code,
			})

			len := len(c.unSubChan)
			for i := 0; i < len; i++ {
				topic := <-c.unSubChan
				ids = append(ids, &pusherrpc.Id{
					Exchange: topic.Exchange,
					Code:     topic.Code,
				})
			}

			err := stream.Send(&pusherrpc.RealTimeQuoteRequest{
				Op:     pusherrpc.Op_UNSUB,
				UserId: c.userId,
				Ids:    ids,
			})
			if err != nil {
				panic(err)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) GetReadChannel() chan Quote {
	c.startOnce.Do(func() {
		go c.start()
	})
	return c.readBuffer
}

func (c *Client) Sub(topic Id) {
	c.subChan <- topic
}

func (c *Client) UnSub(topic Id) {
	c.unSubChan <- topic
}
