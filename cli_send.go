package main

import (
    "fmt"
)

func (cli *CLI) send(from, to string, amount int) {
    bc := NewBlockchain()
    defer bc.db.Close()

    tx := NewUTXOTransaction(from, to, amount, bc)

    // 给矿工的奖励交易
    cbTx := NewCoinbaseTX(from, "")

    bc.MineBlock([]*Transaction{cbTx, tx})
    fmt.Println("success!")
}
