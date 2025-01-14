package main

import (
	"context"
	"fmt"
	client "mosn.io/layotto/sdk/go-sdk/client"
	"time"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	testHello(cli)
	testSet(cli)
	testDelete(cli)
	testWatch(cli)
	cli.Close()
}

func testHello(cli client.Client) {
	req := &client.SayHelloRequest{ServiceName: "helloworld"}
	resp, err := cli.SayHello(context.Background(), req)
	if err != nil {
		fmt.Printf("say hello error: %+v", err)
		return
	}
	fmt.Printf("receive hello response: %+v \n", resp.Hello)
}

func testWatch(cli client.Client) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	item1 := &client.ConfigurationItem{Key: "hello1", Content: "world1"}
	item2 := &client.ConfigurationItem{Key: "hello2", Content: "world2"}
	saveRequest := &client.SaveConfigurationRequest{StoreName: "etcd", AppId: "sofa"}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)

	getRequest := &client.ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if cli.SaveConfiguration(ctx, saveRequest) != nil {
				fmt.Println("save key failed")
				continue
			}
			err = cli.DeleteConfiguration(ctx, getRequest)
			if err != nil {
				fmt.Printf("delete key failed, %+v \n", err)
				continue
			}
			cancel()
			return
		}
	}()

	item := &client.ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1"}}
	ch := cli.SubscribeConfiguration(ctx, item)
	for wc := range ch {
		for _, v := range wc.Item.Items {
			fmt.Printf("receive watch event, %+v \n", v)
		}
	}
}

func testSet(cli client.Client) {
	item1 := &client.ConfigurationItem{Key: "hello1", Content: "world1"}
	item2 := &client.ConfigurationItem{Key: "hello2", Content: "world2"}
	saveRequest := &client.SaveConfigurationRequest{StoreName: "etcd", AppId: "sofa"}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)

	getRequest := &client.ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}

	if cli.SaveConfiguration(context.Background(), saveRequest) != nil {
		fmt.Println("save key failed")
		return
	}

	items, err := cli.GetConfiguration(context.Background(), getRequest)
	if err != nil {
		fmt.Printf("get configuration failed %+v \n", err)
	}

	for _, item := range items {
		fmt.Printf("get configuration after save, %+v \n", item)
	}
}

func testDelete(cli client.Client) {
	item1 := &client.ConfigurationItem{Key: "hello1", Content: "world1"}
	item2 := &client.ConfigurationItem{Key: "hello2", Content: "world2"}
	saveRequest := &client.SaveConfigurationRequest{StoreName: "etcd", AppId: "sofa"}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)

	getRequest := &client.ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	if cli.DeleteConfiguration(context.Background(), getRequest) != nil {
		fmt.Println("save key failed")
		return
	}

	items, err := cli.GetConfiguration(context.Background(), getRequest)
	if err != nil {
		fmt.Printf("get configuration failed %+v \n", err)
	}

	for _, item := range items {
		fmt.Printf("get configuration after save, %+v \n", item)
	}
}
