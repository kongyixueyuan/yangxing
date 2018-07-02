package blc


//money
//addr
type TXOutput struct {
	Value int64
	ScriptPubKey  string
}

func (txoutput *TXOutput)UnLockScriptPubKeyWithAddress(key string) bool {

//input key == PubKEY
	return txoutput.ScriptPubKey == key
}