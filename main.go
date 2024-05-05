package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/exec"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	rpcURL := "https://eth-sepolia.g.alchemy.com/v2/GLYjGj3YzX6ayV_RlGe--fhdxbcuyDBv"
	ethRPC, err := rpc.DialHTTP(rpcURL)
	if err != nil {
		panic(err)
	}
	client := ethclient.NewClient(ethRPC)

	_ = client

	// CheckCorrectness("16424131")

	// BlockHash(client, big.NewInt(16424131))
	// BlockNumber(client)
	// SendETH(client)
	// BlockNumber(client)

	CallContractAtHash(client)
}

func SendETH(client *ethclient.Client) {
	sender := common.HexToAddress("0xb8AE85bF3A5C2C179498fe2093733F01746f5102")
	senderPrivateKey, err := crypto.HexToECDSA("0a88f84a4211da88db9b083d2d29e0b58ee2992e76b82cc94f47deb546f03d87")
	if err != nil {
		panic(err)
	}
	receiver := common.HexToAddress("0x43a407a1dFF44dBBd2f7cf6C61aaF2Ba6653BA59")

	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		panic(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	var data []byte
	tx := types.NewTransaction(nonce, receiver, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(chainID)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), senderPrivateKey)
	if err != nil {
		panic(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}

func BlockNumber(client *ethclient.Client) {
	number, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("block number", number)
}

func CallContractAtHash(client *ethclient.Client) {
	contract := common.HexToAddress("0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238")
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}

	bNBigInt := big.NewInt(int64(blockNumber))
	for i := 0; i < 100; i++ {
		position := fmt.Sprintf("%x", i)
		bt, err := client.StorageAt(context.Background(), contract, common.BytesToHash(Keccak256(position)), bNBigInt)
		if err != nil {
			panic(err)
		}
		if hex.EncodeToString(bt) != "0000000000000000000000000000000000000000000000000000000000000000" {
			fmt.Println(hex.EncodeToString(bt))
			return
		}
		fmt.Println("continue... ", i+1)
	}
}

func Keccak256(position string) []byte {
	missingZeroNumbers := 64 - len(position)
	prefix := ""
	for range missingZeroNumbers {
		prefix += "0"
	}
	fmt.Println(len(prefix + position))
	addressBt := "00000000000000000000000075C0c372da875a4Fc78E8A37f58618a6D18904e8" + prefix + position

	bt, err := hex.DecodeString(addressBt)
	if err != nil {
		panic(err)
	}

	bt = crypto.Keccak256(bt)

	fmt.Println(hex.EncodeToString(bt))
	return bt
}

func BlockHash(client *ethclient.Client, blockNumber *big.Int) {
	hash, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", hash)
	fmt.Println("hash", hash.Hash())
}

func GetBlock(client *ethclient.Client) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(5827167))
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Time())
}

func GetBlockOfTimeStamp(client *ethclient.Client, ts uint64) {
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	startBlock := currentBlock - 100

	for {
		midBlock := (startBlock + currentBlock) / 2
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(midBlock)))
		if err != nil {
			panic(err)
		}
		if block.Time() >= ts {
			currentBlock = midBlock
			fmt.Println(block.Time())
			fmt.Println("continued")
			continue
		} else {
			midBlockPlusOne, err := client.BlockByNumber(context.Background(), big.NewInt(int64(midBlock + 1)))
			if err != nil {
				panic(err)
			}
			if midBlockPlusOne.Time() >= ts {
				fmt.Println(block.Time())
				startBlock = midBlock
				break
			} else {
				startBlock = midBlock
				fmt.Println("continued from bottom")
				continue
			}
		}

		
	}
	fmt.Println(startBlock)
}

func BlockTime(client *ethclient.Client, bn uint64) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Time())
}

func CheckCorrectness(blockNumber string) {
	cmd := exec.CommandContext(context.Background(), "/Users/swopnilparajuli/workspace/cedro/MetalayerMaterials/zeth/target/release/zeth", "build", "--network", "ethereum",
		"--eth-rpc-url", "https://eth.llamarpc.com", "--cache", "host/testdata", " --block-number", blockNumber)

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}
