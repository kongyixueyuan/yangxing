package main

import (
	"./blc"
	"crypto/sha256"
	"fmt"
)

func main(){
	/*
	blc.CreateGenesisBlockChainWithBlock("Genesis Block")
	blockchain := blc.BlockchainObject()
	fmt.Println(blockchain)
	blockchain.AddBlockToBlockchain("yangxing transfer 10RMB to lili")
	fmt.Println(blc.BlockchainObject())
	*/
	hash := sha256.New()
	hash.Write([]byte("yangxing"))
	bytes := hash.Sum(nil)
	fmt.Printf("%x\n",bytes)
	cli := blc.Cli{}
	cli.Run()
}