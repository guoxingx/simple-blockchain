package main

import (
    "fmt"
)

func (cli *CLI) send(from, to string, amount int) {
    bc := NewBlockchain()
    u := &UTXOSet{bc}
    defer u.Blockchain.db.Close()

    tx := NewUTXOTransaction(from, to, amount, u)

    newBlock := bc.MineBlock(from, []*Transaction{tx})
    u.Update(newBlock)
    fmt.Println("success!")
}
