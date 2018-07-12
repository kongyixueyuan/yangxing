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
	fmt.Println("\t send -from --Addr -to --Addr -money --value")
	fmt.Println("\t getbalance -addr --Addr")
	fmt.Println("\t createwallet --createwallet")
	fmt.Println("\t listAddress --list all the address")
	fmt.Println("\t resetwallet")
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

	blockchain := CreateGenesisBlockChainWithBlock(data)
	defer blockchain.DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.ResetUTXOSet()
}

func (cli *Cli) resetwallet()  {

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.ResetUTXOSet()
}

func (cli *Cli) getBalance(addr string) {

	if dbExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer  blockchain.DB.Close()
	//return 	blockchain.getBalance(addr)
	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.GetBalance(addr)

	fmt.Printf("%s一共有%d个Token\n",addr,amount)
}

func (cli *Cli)createwallet() {
	wallets,_ := NewWallets()

	wallets.CreateWallets()
}

func (cli *Cli)listAddress() {
	wallets,_ := NewWallets()

	for addr,_ := range wallets.Walletmap {
		fmt.Println(addr)
	}
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
	createwalletCmd := flag.NewFlagSet("createwalletCmd",flag.ExitOnError)
	createblockChainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listAddress",flag.ExitOnError)
	resetwalletCmd := flag.NewFlagSet("resetwallet",flag.ExitOnError)
	//fmt.Println(sendCmd.Name())

	flagFrom := sendCmd.String("from","","源地址")
	flagTo := sendCmd.String("to","","目标地址")
	flagAmount := sendCmd.String("money","","转账金额......")
	//
	//fmt.Println(flagFrom)
	//fmt.Println(flagTo)
	//fmt.Println(flagAmount)

	flagAddBlockData := addBlockCmd.String("data","","You need type here your trasfer")

	flagCreateBlockChainWhisData := createblockChainCmd.String("data","","创世区块")

	flagGetBalanceData := getBalanceCmd.String("addr","","想要查询的地址")
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
		case "send":
			err := sendCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "getbalance":
			err := getBalanceCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createwallet":
			err := createwalletCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "listAddress":
			err := listAddressCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "resetwallet":
			err := resetwalletCmd.Parse(os.Args[2:])
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

	if createwalletCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.createwallet()
	}
	if listAddressCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.listAddress()
	}

	if createblockChainCmd.Parsed() {

		if *flagCreateBlockChainWhisData == "" {
			fmt.Println("交易数据不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockChainWhisData)
	}

	if sendCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			printUsage()
			os.Exit(1)
		}
		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
/*
		fmt.Println(from)
		fmt.Println(to)
		fmt.Println(amount)
*/
		cli.send(from,to,amount)
	}

	if getBalanceCmd.Parsed() {

		if *flagGetBalanceData == "" {
			fmt.Println("地址不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*flagGetBalanceData)
		//fmt.Printf("addr = %s value == %d\n",*flagGetBalanceData,value)
	}

	if resetwalletCmd.Parsed() {
		cli.resetwallet()
		//fmt.Printf("addr = %s value == %d\n",*flagGetBalanceData,value)
	}

}


func (cli *Cli) send(from []string,to []string,amount []string)  {


	if dbExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	//
	blockchain.CreateNewBlockWithTransaction(from,to,amount)

	utxoSet := &UTXOSet{blockchain}

	//转账成功以后，需要更新一下
	utxoSet.Update()

}
