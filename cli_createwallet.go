package main

import (
    "fmt"
)

// 创建新账号
func (cli *CLI) createWallet() {
    wallets, _ := NewWallets()
    address := wallets.CreateWallet()
    wallets.SaveToFile()

    fmt.Printf("Your new address: %s\n", address)
}
