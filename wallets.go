package main

import (
    "os"
    "log"
    "fmt"
    "bytes"
    "io/ioutil"
    "crypto/elliptic"
    "encoding/gob"
)

const walletFile = "data/wallet.dat"

type Wallets struct {
    Wallets map[string]*Wallet
}

// 
func NewWallets() (*Wallets, error) {
    wallets := Wallets{}
    wallets.Wallets = make(map[string]*Wallet)

    err := wallets.LoadFromFile()

    return &wallets, err
}

//
func (wallets *Wallets) CreateWallet() string {
    wallet := NewWallet()

    // wallet.GetAddress return []byte
    address := fmt.Sprintf("%s", wallet.GetAddress())
    wallets.Wallets[address] = wallet

    return address
}

//
func (wallets *Wallets) GetAddresses() []string {
    var addresses []string

    for address := range wallets.Wallets {
        addresses = append(addresses, address)
    }

    return addresses
}

//
func (wallets *Wallets) GetWallet(address string) Wallet {
    return *wallets.Wallets[address]
}

// 从文件中加载wallet
func (wallets *Wallets) LoadFromFile() error {
    if _, err := os.Stat(walletFile); os.IsNotExist(err) {
        return err
    }

    fileContent, err := ioutil.ReadFile(walletFile)
    if err != nil { log.Panic(err) }

    var wallets_loaded Wallets

    // gob.Register         f func(value interface{}) 
    // elliptic.P256        f func() elliptic.Curve
    // elliptic.Curve       t interface
    gob.Register(elliptic.P256())

    // gob.NewDecoder f func(r io.Reader) *gob.NewDecoder 
    decoder := gob.NewDecoder(bytes.NewReader(fileContent))

    // decoder.Decode      f func(e interface{}) error
    err = decoder.Decode(&wallets_loaded)
    if err != nil { log.Panic(err) }

    wallets.Wallets = wallets_loaded.Wallets

    return nil
}

func (wallets Wallets) SaveToFile() {
    var content bytes.Buffer

    gob.Register(elliptic.P256())

    // gob.NewEncoder  f func(w io.Writer) *gob.Encoder
    encoder := gob.NewEncoder(&content)

    // encoder.Encode  f func(e interface{}) error
    err := encoder.Encode(wallets)
    if err != nil { log.Panic(err) }

    // ioutil.WriteFile  f func(filename string, data []byte, perm os.FileMode) error
    err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
    if err != nil { log.Panic(err) }
}
