package blc

type TXInput struct {
	Txhash []byte
	Vout      int
	ScriptSig string
}

func (txinput *TXInput) UnLockWithAddress(addr string) bool {
	return txinput.ScriptSig == addr
}

