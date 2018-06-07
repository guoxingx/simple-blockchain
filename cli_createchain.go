package main

import (
    "fmt"
    "log"
)

func (cli *CLI) createChain(address string) {
    if !ValidateAddress(address) { log.Panic("Error: Address is not valid") }

	bc := CreateBlockchain(address)
	defer bc.db.Close()

    UTXOSet := UTXOSet{bc}
    UTXOSet.Reindex()

	fmt.Println("Done!")
}
