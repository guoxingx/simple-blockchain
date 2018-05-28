package main

import (
    "fmt"
    "log"
)

func (cli *CLI) listAddress() {
    wallets, err := NewWallets()
    if err != nil { log.Panic(err) }

    addresses := wallets.GetAddresses()

    for _, address := range addresses {
        fmt.Printf("%v, ", address)
    }
    fmt.Println()
}
