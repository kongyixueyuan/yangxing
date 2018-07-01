package main

import(
	"fmt"
	"./blc"
)

func main(){
	blockchain := blc.CreateGenesisBlockChainWithBlock()
	blockchain.AddNewBlock("yangxing transfer 10RMB to lili")
	fmt.Println(blockchain)
}