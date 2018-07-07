package blc

import (
	"math/big"
	"math/bits"
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func base58Encode(input []byte) []byte {

	var result []byte
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(b58Alphabet)))
	mod := &big.Int{}

	zero := big.NewInt(0)


	for x.Cmp(zero) != 0{
		x.DivMod(x,base,mod)
		result = append(result,b58Alphabet[mod.Int64()])
	}
	ReverseByte(result)





	return result
}

func base58Decode(input []byte) []byte {

}
