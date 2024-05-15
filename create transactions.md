A transaction is recording an event of bitcoin transition, it need to make sure where the bitcoins transfer to, is the input for this transaction legal or valid, and how quickly the transaction can on chain which
means the transaction is legally accepted.

Let's see how to construct a valid transaction and send it to the network, this process may be like you first go to the bank, deposit some amount of money in your account and transfer some of them to your friend.
The first thing we need to do is convert an wallet address from base58 encode, that is we need a process that can get its orignal content when the input is encodes by base58, let's check the code first, remember
we have EncodeBase58 first in util.go of elliptic-curve, let's add the method of base58 decode there too:
```go
func DecodeBase58(s string) []byte {
	BASE58_ALPHABET := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	num := big.NewInt(int64(0))
	for _, char := range s {
		mulOp := new(big.Int)
		num = mulOp.Mul(num, big.NewInt(int64(58)))

		idx := strings.Index(BASE58_ALPHABET, string(char))
		if idx == -1 {
			panic("can't find char in base58 alphabet")
		}
		addOp := new(big.Int)
		num = addOp.Add(num, big.NewInt(int64(idx)))
	}

	combined := num.Bytes()
	checksum := combined[len(combined)-4:]
	h256 := Hash256(string(combined[0 : len(combined)-4]))
	if bytes.Equal(h256[0:4], checksum) != true {
		panic("decode base58 checksum error")
	}

	//the first byte is network prefix
	return combined[1 : len(combined)-4]
}

```
Now let's test the above code in main.go:
```g
package main

import (
	ecc "elliptic_curve"
	"fmt"
)

func main() {
	res := ecc.DecodeBase58("mzx5YhAH9kNHtcN481u6WkjeHjYtVeKVh2")
	fmt.Printf("decode result is %x\n", res)
}
```
The aboved code gives the following result:
```g
decode result is d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f
```
