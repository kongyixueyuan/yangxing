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
	/*
	hash := sha256.New()
	hash.Write([]byte("yang00000xing"))
	bytes := hash.Sum(nil)
	fmt.Printf("%x\n",bytes)
	test := blc.Base58Encode([]byte("sdfaasdf"))
	fmt.Printf("%x\n",test)
	back := blc.Base58Decode(test)
	fmt.Printf("%x\n",back)
	fmt.Println("origin = ",test)
	fmt.Println("back =",back)
*/
	cli := blc.Cli{}
	cli.Run()

}