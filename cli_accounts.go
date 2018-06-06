package main

import (
    "fmt"
    "log"
)

func (cli *CLI) accounts() {
    wallets, err := NewWallets()
    if err != nil { log.Panic(err) }

    accounts := wallets.GetAddresses()

    for _, account := range accounts {
        fmt.Printf("%v, ", account)
    }
    fmt.Println()
}
