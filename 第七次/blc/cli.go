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
	fmt.Println("\tstartnode -miner ADDRESS -- 启动节点服务器，并且指定挖矿奖励的地址.")
}


func (cli *Cli) addBlock(data string,nodeID string)  {

	if YX_dbExists(nodeID) == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject(nodeID)

	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockchain(data)
}

func (cli *Cli) printchain(nodeID string)  {

	if YX_dbExists(nodeID) == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject(nodeID)

	defer blockchain.DB.Close()

	blockchain.Printchain()

}

func (cli *Cli) createGenesisBlockchain(data string,nodeID string)  {

	blockchain := CreateGenesisBlockChainWithBlock(data,nodeID)
	defer blockchain.DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.ResetUTXOSet()
}

func (cli *Cli) resetwallet(nodeID string)  {

	blockchain := BlockchainObject(nodeID)
	defer blockchain.DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.ResetUTXOSet()
}

func (cli *Cli) getBalance(addr string,nodeID string) {

	if YX_dbExists(nodeID) == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject(nodeID)
	defer  blockchain.DB.Close()
	//return 	blockchain.getBalance(addr)
	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.GetBalance(addr)

	fmt.Printf("%s一共有%d个Token\n",addr,amount)
}

func (cli *Cli)YX_createwallet(nodeID string) {
	wallets,_ := YX_NewWallets(nodeID)

	wallets.CreateWallets()
}

func (cli *Cli)YX_listAddress(nodeID string) {
	wallets,_ := YX_NewWallets(nodeID)

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

func (cli *Cli) startNode(nodeID string,minerAdd string)  {

	// 启动服务器

	if minerAdd == "" || YX_IsValidForAdress([]byte(minerAdd))  {
		//  启动服务器
		fmt.Printf("启动服务器:localhost:%s\n",nodeID)
		YX_startServer(nodeID,minerAdd)

	} else {

		fmt.Println("指定的地址无效....")
		os.Exit(0)
	}

}



func (cli *Cli) Run()  {
	//first you need to set env NODE_ID
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!\n")
		os.Exit(1)
	}

	fmt.Printf("NODE_ID:%s\n",nodeID)
	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addBlockToBlockchain",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createwalletCmd := flag.NewFlagSet("createwalletCmd",flag.ExitOnError)
	createblockChainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listAddress",flag.ExitOnError)
	resetwalletCmd := flag.NewFlagSet("resetwallet",flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode",flag.ExitOnError)
	//fmt.Println(sendCmd.Name())

	flagFrom := sendCmd.String("from","","源地址")
	flagTo := sendCmd.String("to","","目标地址")
	flagAmount := sendCmd.String("money","","转账金额......")
	flagMine := sendCmd.Bool("mine",false,"是否在当前节点中立即验证....")
	//
	//fmt.Println(flagFrom)
	//fmt.Println(flagTo)
	//fmt.Println(flagAmount)

	flagAddBlockData := addBlockCmd.String("data","","You need type here your trasfer")

	flagCreateBlockChainWhisData := createblockChainCmd.String("data","","创世区块")

	flagGetBalanceData := getBalanceCmd.String("addr","","想要查询的地址")
	flagMiner := startNodeCmd.String("miner","","定义挖矿奖励的地址......")
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
		case "startnode":
			err := startNodeCmd.Parse(os.Args[2:])
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
		cli.addBlock(*flagAddBlockData,nodeID)

	}

	if createwalletCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.YX_createwallet(nodeID)
	}
	if listAddressCmd.Parsed() {

		//fmt.Println("输出所有区块的数据........")
		cli.YX_listAddress(nodeID)
	}

	if createblockChainCmd.Parsed() {

		if *flagCreateBlockChainWhisData == "" {
			fmt.Println("交易数据不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockChainWhisData,nodeID)
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
		cli.send(from,to,amount,nodeID,*flagMine)
	}

	if getBalanceCmd.Parsed() {

		if *flagGetBalanceData == "" {
			fmt.Println("地址不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*flagGetBalanceData,nodeID)
		//fmt.Printf("addr = %s value == %d\n",*flagGetBalanceData,value)
	}

	if resetwalletCmd.Parsed() {
		cli.resetwallet(nodeID)
		//fmt.Printf("addr = %s value == %d\n",*flagGetBalanceData,value)
	}
	if startNodeCmd.Parsed() {


		cli.startNode(nodeID,*flagMiner)
	}
}


func (cli *Cli) send(from []string,to []string,amount []string,nodeID string,mineNow bool)  {

	blockchain := BlockchainObject(nodeID)
	defer blockchain.DB.Close()

	if mineNow {
		//
		blockchain.CreateNewBlockWithTransaction(from, to, amount,nodeID)

		utxoSet := &UTXOSet{blockchain}

		//转账成功以后，需要更新一下
		utxoSet.Update()
	} else {
		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理......")
	}
}
