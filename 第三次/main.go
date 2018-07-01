package main

import(
	"fmt"
	"./blc"
)

func main(){
	blc.CreateGenesisBlockChainWithBlock("Genesis Block")
	blc.BlockchainObject().AddBlockToBlockchain("yangxing transfer 10RMB to lili")
	fmt.Println(blc.BlockchainObject())
}