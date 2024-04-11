package main

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"os"
	"strconv"
	"testing"
	"time"
)

func FuzzCar(f *testing.F) {
	f.Add("CreateAsset", "1", "yellow", "1", "Tom", "1300")
	f.Add("CreateAsset", "2", "yellow", "1", "Tom", "1300")
	f.Add("CreateAsset", "3", "yellow", "1", "Tom", "1300")
	f.Add("CreateAsset", "4", "yellow", "1", "Tom", "1300")
	f.Add("CreateAsset", "600", "red", "1", "To", "130")
	f.Add("CreateAsset", "100", "red", "1", "To", "130")

	//需要针对  已经输入的id排除  不能重复 
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific chaincodeSupportClient identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	initLedger(contract)

	f.Fuzz(func(t *testing.T, contractName string, id string, color string, num string, person string, value string) {
		// 检查 id 是否是在范围 [1, 1000] 内的整数
		//之后可以使用 转化数  提高效率
		idInt, err := strconv.Atoi(id)
		if err != nil || idInt < 1 || idInt > 1000 {
			t.Logf("Skipping test with id: %s", id)
			return
		}

		// 检查 color 是否是合法的颜色
		if color != "red" && color != "blue" && color != "green" {
			t.Logf("Skipping test with color: %s", color)
			return
		}

		// 检查 num 是否是在范围 [1, 9999] 内的整数
		numInt, err := strconv.Atoi(num)
		if err != nil || numInt < 1 || numInt > 9999 {
			t.Logf("Skipping test with num: %s", num)
			return
		}

		// 检查 value 是否是在范围 [1, 999999] 内的整数
		valueInt, err := strconv.Atoi(value)
		if err != nil || valueInt < 1 || valueInt > 999999 {
			t.Logf("Skipping test with value: %s", value)
			return
		}

		// 如果所有数据都符合要求，则进行测试 Create Asset
		fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new asset with ID, Color, Size, Owner and AppraisedValue arguments \n")
		_, err = contract.SubmitTransaction("CreateAsset", id, color, num, person, value)
		// _, err = contract.SubmitTransaction(contractName, id, color, num, person, value)
		if err != nil {
			t.Errorf("failed to submit transaction: %s，contractName = %s, id = %s,color = %s, num = %s, person = %s, value = %s\n", err,contractName, id, color, num, person, value)
			panic(fmt.Errorf("failed to submit transaction: %w", err))
		}

		t.Logf("*** Transaction committed successfully\n")
		fmt.Printf("*** Transaction committed successfully\n")

		// 测试 Read Asset
		t.Logf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")
		fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")
		_, err = contract.SubmitTransaction("ReadAsset", id)
		if err != nil {
			t.Errorf("failed to submit transaction: %s", err)
			//panic(fmt.Errorf("failed to submit transaction: %w", err))
		}
		t.Logf("*** Transaction committed successfully\n")
		fmt.Printf("*** Transaction committed successfully\n")

		//测试 Update Asset
		//测试其他
	})

}
