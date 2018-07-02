package blc

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type Cli struct {}

func printUsage()  {

	fmt.Println("Usage:")
	fmt.Println("\t createblockchain -data -- 交易数据.")
	fmt.Println("\t addBlockToBlockchain -data DATA -- 交易数据.")
	fmt.Println("\t printchain -- 输出区块信息.")
	fmt.Println("\t send -from --Addr -to --Addr")
}


func (cli *Cli) addBlock(data string)  {

	if dbExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockchain(data)
}

func (cli *Cli) printchain()  {

	if dbExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.Printchain()

}

func (cli *Cli) createGenesisBlockchain(data string)  {

	CreateGenesisBlockChainWithBlock(data)
}


func isValidArgs()  {
	if (len(os.Args) < 2) {
		printUsage()
		os.Exit(0)
	}
}

func (cli *Cli) Run()  {

	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addBlockToBlockchain",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createblockChainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")

	flagAddBlockData := addBlockCmd.String("data","yangxing trasfer 100RMB to lili","You need type here your trasfer")

	flagCreateBlockChainWhisData := createblockChainCmd.String("data","Genesis block","创世区块")
	//fmt.Println(*flagAddBlockData)

	switch os.Args[1] {
		case "addBlockToBlockchain":
			err := addBlockCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createblockchain":
			err := createblockChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			printUsage()
			os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}

		//fmt.Println(*flagAddBlockData)
		cli.addBlock(*flagAddBlockData)

	}

	if printChainCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.printchain()
	}

	if createblockChainCmd.Parsed() {

		if *flagCreateBlockChainWhisData == "" {
			fmt.Println("交易数据不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockChainWhisData)
	}


	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			printUsage()
			os.Exit(1)
		}



		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
	}

}


func (cli *Cli) send(from []string,to []string,amount []string)  {


	if dbExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from,to,amount)

}
