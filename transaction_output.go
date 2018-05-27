package main

import (
    "bytes"
)

type TXOutput struct {
    Value      int
    PubKeyHash []byte // Pubkey hash
}

func (out *TXOutput) Lock(address []byte) {
    pubKeyHash := Base58Decode(address)
    pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
    out.PubKeyHash = pubKeyHash
}

// 是否和 out.PubKeyHash 相等
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
    return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
    txo := &TXOutput{ value, nil }
    txo.Lock([]byte(address))

    return txo
}
