package main

import (
    "fmt"
)

func (cli *CLI) send(from, to string, amount int) {
    bc := NewBlockchain()
    u := &UTXOSet{bc}
    defer u.Blockchain.db.Close()

    tx := NewUTXOTransaction(from, to, amount, u)

    // 给矿工的奖励交易
    cbTx := NewCoinbaseTX(from, "")

    newBlock := bc.MineBlock([]*Transaction{cbTx, tx})
    u.Update(newBlock)
    fmt.Println("success!")
}
