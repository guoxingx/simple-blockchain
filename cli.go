package main

import (
    "os"
    "fmt"
    "flag"
    "log"
)

// init with a blockchain
type CLI struct {}

const usage = `
Usage:
  printchain                             print all the blocks of the blockchain
  createblockchain -address ADDRESS      Create a blockchain and send genesis block reward to ADDRESS
  createwallet                           Generates a new key-pair and saves it into the wallet file
  listaddress                            Lists all addresses from the wallet file
  getbalance -address ADDRESS            Get balance of ADDRESS
  send -from FROM -to TO -amount AMOUNT  Send AMOUNT of coins from FROM address to TO
`

func (cli *CLI) Run() {
    cli.validateArgs()

    // NewFlagSet  f func(name string, errorHandling flag.ErrorHandling) *flag.FlagSet
    printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
    createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
    createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
    listAddressCmd := flag.NewFlagSet("listaddress", flag.ExitOnError)
    getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
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
    case "createwallet":
        err := createWalletCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "listaddress":
        err := listAddressCmd.Parse(os.Args[2:])
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
    if printChainCmd.Parsed() { cli.printChain() }

    if createBlockchainCmd.Parsed() {
        if *createBlockchainData == "" {
            createBlockchainCmd.Usage()
            os.Exit(1)
        }
        cli.createBlockchain(*createBlockchainData)
    }

    if createWalletCmd.Parsed() { cli.createWallet() }

    if listAddressCmd.Parsed() { cli.listAddress() }

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
