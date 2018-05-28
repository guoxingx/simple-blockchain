package main

import (
    "fmt"
    "strconv"
)

// print each block and validate pow.
func (cli *CLI) printChain() {
    bc := NewBlockchain()
    bci := bc.Iterator()

    for {
        block := bci.Next()

        fmt.Printf("============ Block %v %x ============\n", block.Number, block.Hash)
        fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
        pow := NewProofOfWork(block)
        fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
        fmt.Printf("Transactions: ")
        for _, tx := range block.Transactions {
            fmt.Printf("%x, ", tx.ID)
        }
        fmt.Println("\n")

        if len(block.PrevBlockHash) == 0 { break }
    }
}
