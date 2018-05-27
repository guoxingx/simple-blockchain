package main

import (
    "fmt"
)

type Transaction struct {
    ID   []byte
    Vin  []TXInput
    Vout []TXOutput
}

type TXOutput struct {
    Value        int
    ScriptPubKey string
}

type TXInput struct {
    Txid      []byte
    Vout      int
    ScriptSig string
}

const subsidy = 50

// coinbase 交易
// 即区块的奖励交易
func NewCoinbaseTX(to, data string) *Transaction {
    if data == "" {
        data = fmt.Sprintf("Reward to '%s'", to)
    }

    txin := TXInput{[]byte{}, -1, data}
    txout := TXOutput{subsidy, to}
    tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
    tx.SetID()

    return &tx
}

func NewTx(from, to, data string) *Transaction {
    return nil
}

func (tx *Transaction) SetID() {

}
