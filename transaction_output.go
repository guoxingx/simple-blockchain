package main

import (
    "log"
    "bytes"
    "encoding/gob"
)

type TXOutput struct {
    Value      int
    PubKeyHash []byte // Pubkey hash
}

// 根据address 设置 out.pubKeyHash
func (out *TXOutput) Lock(address []byte) {
    pubKeyHash := Base58Decode(address)
    pubKeyHash = pubKeyHash[1:len(pubKeyHash) - 4]
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

type TXOutputs struct {
    Outputs []TXOutput
}

// Serialize serializes TXOutputs
func (outs TXOutputs) Serialize() []byte {
    var buff bytes.Buffer

    enc := gob.NewEncoder(&buff)
    err := enc.Encode(outs)
    if err != nil { log.Panic(err) }

    return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TXOutputs {
    var outputs TXOutputs

    dec := gob.NewDecoder(bytes.NewReader(data))
    err := dec.Decode(&outputs)
    if err != nil { log.Panic(err) }

    return outputs
}
