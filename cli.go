package main

import (
    "os"
    "fmt"
    "flag"
    "strconv"
    "log"
)

// init with a blockchain
type CLI struct {}

const usage = `
Usage:
  addblock -data BLOCK_DATA           add a block to the blockchain
  printchain                          print all the blocks of the blockchain
  createblockchain -address ADDRESS   Create a blockchain and send genesis block reward to ADDRESS
`

func (cli *CLI) Run() {
    cli.validateArgs()

    // NewFlagSet  f func(name string, errorHandling flag.ErrorHandling) *flag.FlagSet
    addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
    printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
    createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

    // flag.FlagSet.String  f func(name string, value string, usage string) *string
    addBlockData := addBlockCmd.String("data", "", "Block data")
    createBlockchainData := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

    switch os.Args[1] {
    case "addblock":
        // flag.FlagSet.Parse  f func(arguments []string) error
        err := addBlockCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "printchain":
        err := printChainCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "createblockchain":
        err := createBlockchainCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    default:
        cli.printUsage()
        os.Exit(1)
    }

    // flag.FlagSet.Parsed f func() bool
    if addBlockCmd.Parsed() {
        if *addBlockData == "" {
            addBlockCmd.Usage()
            os.Exit(1)
        }
        cli.addBlock(*addBlockData)
    }

    if printChainCmd.Parsed() {
        cli.printChain()
    }

    if createBlockchainCmd.Parsed() {
        if *createBlockchainData == "" {
            createBlockchainCmd.Usage()
            os.Exit(1)
        }
        cli.createBlockchain(*createBlockchainData)
    }
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
    fmt.Println(usage)
}

// call blockchain.AddBlock() with data
func (cli *CLI) addBlock(data string) {
    bc := NewBlockchain("")
    bc.AddBlock(data)
    fmt.Println("Success!")
}

// print each block and validate pow.
func (cli *CLI) printChain() {
    bc := NewBlockchain("")
    bci := bc.Iterator()

    for {
        block := bci.Next()

        fmt.Printf("Number: %v\n", block.Number)
        fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
        fmt.Printf("Transactions: %s\n", block.Transactions)
        fmt.Printf("Hash: %x\n", block.Hash)

        pow := NewProofOfWork(block)
        fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
        fmt.Println()

        if len(block.PrevBlockHash) == 0 { break }
    }
}

func (cli *CLI) createBlockchain(address string) {
	bc := NewBlockchain(address)
	bc.db.Close()
	fmt.Println("Done!")
}