package main

import (
	"fmt"
	"os"
	"strconv"
	"flag"
	"log"
)

type CLI struct {
	bc *Blockchain  // CLI 结构体中放一个 bc 对象是否合适？
	           // 命令行中应该放什么？
}


func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("  addblock")
	fmt.Println("    -data BLOCK_DATA    add a blcok to the blockchain")
	fmt.Println("  printchain            print all the blocks of the blockchain")
}


func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
}


func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}


func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 注意，上面的 FlagSet 只是一个集合 addblock，然后下面可以定义一系列标志，比如-p 等
	// 如果没有定义其他的，则默认只有上面的
	// 返回值是一个字符串指针，String 函数的一个功能是将这个 data 标志添加
	// 到了 flag set 当中
	// 即一个是子命令，一个是子命令的子标志
	addBclockData := addBlockCmd.String("data", "", "Blockdata")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	// 注意 flagSet 有 parsed 方法，然后子标志没有
	if addBlockCmd.Parsed() {
		if *addBclockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBclockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}
