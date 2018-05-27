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
  printchain                             print all the blocks of the blockchain
  createblockchain -address ADDRESS      Create a blockchain and send genesis block reward to ADDRESS
  getbalance -address ADDRESS            Get balance of ADDRESS
  send -from FROM -to TO -amount AMOUNT  Send AMOUNT of coins from FROM address to TO
`

func (cli *CLI) Run() {
    cli.validateArgs()

    // NewFlagSet  f func(name string, errorHandling flag.ErrorHandling) *flag.FlagSet
    printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
    createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
    getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
    sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

    // flag.FlagSet.String  f func(name string, value string, usage string) *string
    createBlockchainData := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
    getBalanceData := getBalanceCmd.String("address", "", "The address to get balance for")
    sendFrom := sendCmd.String("from", "", "Source wallet address")
    sendTo := sendCmd.String("to", "", "Destination wallet address")
    sendAmount := sendCmd.Int("amount", 0, "Amount to send")

    switch os.Args[1] {
    case "printchain":
        err := printChainCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "createblockchain":
        err := createBlockchainCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "getbalance":
        err := getBalanceCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "send":
        err := sendCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    default:
        cli.printUsage()
        os.Exit(1)
    }

    // flag.FlagSet.Parsed f func() bool
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

    if getBalanceCmd.Parsed() {
        if *getBalanceData == "" {
            getBalanceCmd.Usage()
            os.Exit(1)
        }
        cli.getBalance(*getBalanceData)
    }

    if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
        cli.send(*sendFrom, *sendTo, *sendAmount)
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

// print each block and validate pow.
func (cli *CLI) printChain() {
    bc := NewBlockchain()
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
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("Done!")
}

func (cli *CLI) getBalance(address string) {
    bc := NewBlockchain()
    defer bc.db.Close()

    balance := 0
    UTXOs := bc.FindUTXO(address)

    for _, out := range UTXOs { balance += out.Value }

    fmt.Printf("getBalance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
    bc := NewBlockchain()
    defer bc.db.Close()

    tx := NewUTXOTransaction(from, to, amount, bc)
    bc.MineBlock([]*Transaction{tx})
    fmt.Println("success!")
}
