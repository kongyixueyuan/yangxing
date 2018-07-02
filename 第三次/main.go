package main

import (
	"./blc"
)

func main(){
	/*
	blc.CreateGenesisBlockChainWithBlock("Genesis Block")
	blockchain := blc.BlockchainObject()
	fmt.Println(blockchain)
	blockchain.AddBlockToBlockchain("yangxing transfer 10RMB to lili")
	fmt.Println(blc.BlockchainObject())
	*/
	cli := blc.Cli{}
	cli.Run()
}