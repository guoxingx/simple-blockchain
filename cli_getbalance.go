package main

import (
    "fmt"
    "log"
)

func (cli *CLI) getBalance(address string) {
    if !ValidateAddress(address) { log.Panic("ERROR: Address is not Valid") }

    bc := NewBlockchain()
    u := &UTXOSet{bc}
    defer bc.db.Close()

    balance := 0
    UTXOs := u.FindUTXO(AddressToPubKeyHash(address))

    for _, out := range UTXOs { balance += out.Value }

    fmt.Printf("getBalance of '%s': %d\n", address, balance)
}
