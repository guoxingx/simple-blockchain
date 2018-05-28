package main

import (
    "fmt"
)

func (cli *CLI) send(from, to string, amount int) {
    bc := NewBlockchain()
    defer bc.db.Close()

    tx := NewUTXOTransaction(from, to, amount, bc)
    bc.MineBlock([]*Transaction{tx})
    fmt.Println("success!")
}
