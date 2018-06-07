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
  createchain -account ACCOUNT      Create a blockchain and send genesis block reward to ACCOUNT
  createwallet                           Generates a new key-pair and saves it into the wallet file
  accounts                               Lists all accounts
  getbalance -account ACCOUNT            Get balance of ACCOUNT
  send -from FROM -to TO -amount AMOUNT  Send AMOUNT of coins from FROM account to TO
`

func (cli *CLI) Run() {
    cli.validateArgs()

    // NewFlagSet  f func(name string, errorHandling flag.ErrorHandling) *flag.FlagSet
    printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
    createChainCmd := flag.NewFlagSet("createchain", flag.ExitOnError)
    createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
    accountsCmd := flag.NewFlagSet("accounts", flag.ExitOnError)
    getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
    sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

    // flag.FlagSet.String  f func(name string, value string, usage string) *string
    createChainData := createChainCmd.String("account", "", "The account to send genesis block reward to")
    getBalanceData := getBalanceCmd.String("account", "", "The account to get balance for")
    sendFrom := sendCmd.String("from", "", "Source wallet account")
    sendTo := sendCmd.String("to", "", "Destination wallet account")
    sendAmount := sendCmd.Int("amount", 0, "Amount to send")

    switch os.Args[1] {
    case "printchain":
        err := printChainCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "createchain":
        err := createChainCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "createwallet":
        err := createWalletCmd.Parse(os.Args[2:])
        if err != nil { log.Panic(err) }
    case "accounts":
        err := accountsCmd.Parse(os.Args[2:])
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

    if createChainCmd.Parsed() {
        if *createChainData == "" {
            createChainCmd.Usage()
            os.Exit(1)
        }
        cli.createChain(*createChainData)
    }

    if createWalletCmd.Parsed() { cli.createWallet() }

    if accountsCmd.Parsed() { cli.accounts() }

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
