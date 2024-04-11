package main

import (
	"fmt"
	"log"
	//"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {
	// 创建 Fabric SDK 实例
	configProvider := config.FromFile("config.yaml") // 替换为你的配置文件路径
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}
	defer sdk.Close()

	// 创建通道客户端
	clientChannelContext := sdk.ChannelContext("mychannel", fabsdk.WithUser("user1"))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		log.Fatalf("Failed to create channel client: %v", err)
	}

	// 构造 gRPC 请求
	request := channel.Request{
		ChaincodeID: "fabcar",   // 替换为你的链码 ID
		Fcn:         "queryCar", // 替换为你的链码中的函数名
		Args: [][]byte{
			[]byte("CAR1"), // 替换为你的链码中的参数
		},
	}

	// 发起 gRPC 请求
	response, err := channelClient.Query(request)
	if err != nil {
		log.Fatalf("Failed to query chaincode: %v", err)
	}

	// 处理链码响应
	fmt.Printf("Response: %s\n", response.Payload)
}
