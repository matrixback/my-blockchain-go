package main

import (
	"fmt"
	"os"
	"strconv"
	"flag"
)

type CLI struct {}


func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("  getbalance")
	fmt.Println("    -address ADDRESS    Get balance of ADDRESS")
	fmt.Println("  createblockchain")
	fmt.Println("    -address ADDRESS Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  send")
	fmt.Println("    -from FROM -to TO -amout AMOUNT")
	fmt.Println("     Send AMOUNT of coins from FROM address to To ")
	fmt.Println("  printchain            print all the blocks of the blockchain")
}


func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}


func (cli *CLI) printChain() {
	bc := NewBlockchain("")
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}


func(cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Done!")
}

func(cli *CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int)  {
	bc := NewBlockchain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 注意，上面的 FlagSet 只是一个集合 addblock，然后下面可以定义一系列标志，比如-p 等
	// 如果没有定义其他的，则默认只有上面的
	// 返回值是一个字符串指针，String 函数的一个功能是将这个 data 标志添加
	// 到了 flag set 当中
	// 即一个是子命令，一个是子命令的子标志
	getBalanceAddress := getBalanceCmd.String("address","", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to end genesis block reward to", )
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	// 为什么第一个参数还需要用 switch case?
	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		NilPanic(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		NilPanic(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		NilPanic(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		NilPanic(err)
	default:
		cli.printUsage()
		os.Exit(1)
	}

	// 注意 flagSet 有 parsed 方法，然后子标志没有
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}


	if createBlockchainCmd.Parsed() {
		if (*createBlockchainAddress) == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
