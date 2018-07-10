package blc

import (
	"os"
	"bytes"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
)

const walletFile  = "Wallets.dat"

type Wallets struct {
	Walletmap map[string] *Wallet
}


func NewWallets() (*Wallets,error) {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		log.Printf("yangxing2\n")
		wallets := &Wallets{}
		wallets.Walletmap = make(map[string]*Wallet)
		return wallets,err
	}

	filecontent ,err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(filecontent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets,nil
}

func (w *Wallets)CreateWallets() {
	wallet := NewWallet()
	w.Walletmap[string(wallet.getAddress())] = wallet
	log.Printf(".walletmap[string(wallet.getAddress())]  = %s",wallet.getAddress())
	w.SaveWallet()
}

func (w *Wallets)SaveWallet()  {
	var content bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err != nil {
		log.Panic(err)
	}
}