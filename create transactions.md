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

The second thing needed to construct transaction is that you have valid input, which means you need someone send you bitcoin to spend,first you need to construct you
bitcoin address on the network like following:
```g
        p := new(big.Int)
	h256 := ecc.Hash256("puy you secret here")
	fmt.Printf("h256: %x\n", h256)
	p.SetBytes(tx.ReverseByteSlice(h256))
	fmt.Printf("p is %x\n", p)

	privateKey := ecc.NewPrivateKey(p)
	pubKey := privateKey.GetPublicKey()
	fmt.Printf("wallet address for you secret is %s\n", pubKey.Address(true, true))
```
In the above code we hash our secret, you should replace the string "puy you secret here" to your own secret here, then you will get an address on the bitcoin testnet,
for example the address created by myself is:
mpNzUycBH6SDU9amLK5raP6Qm71CWNezHv

Then go to bitcoin testnet faucet to get some bitcoin for our testing by using following link:
https://coinfaucet.eu/en/btc-testnet/
It is not always usable but you can try you luck, after receiving the testing bitcoin, we can check our account by using bitcoin testnet explorer:

![截屏2024-05-16 12 37 21](https://github.com/wycl16514/golang-bitcoin-transaction-create-validate/assets/7506958/26d6ebb0-33f5-46d4-8cb0-61feb350d9ed)

you can see I receive 0.00019756 tBTC from the faucet. Then we need to construct the scriptPubKey script as we have seen before, in util.go of transaction add
code like following:

```go
func P2pkScrit(h160 []byte) *ScriptSig {
	scriptContent := [][]byte{[]byte{OP_DUP},[]byte{OP_HASH160}, 
		h160, []byte{OP_EQUALVERIFY}, []byte{OP_CHECKSIG}}
	return InitScriptSig(scriptContent)
}
```
As you can see the above script op codes, we have seen them many times in previous video, you can use the following link :

https://blockstream.info/testnet/address/mpNzUycBH6SDU9amLK5raP6Qm71CWNezHv

scroll down and click the detail button and you will see the following:

![截屏2024-05-16 13 26 21](https://github.com/wycl16514/golang-bitcoin-transaction-create-validate/assets/7506958/397ad84f-e56e-42fd-8ff5-301686ee422c)

The string "703158ce66391f094ab2195cfe5579214073ba90997d0b98e6e410ed1b67aa8a" is the previous transaction hash, then in the part with title 
"mpNzUycBH6SDU9amLK5raP6Qm71CWNezHv" which is the address I created, you can see the scriptpubkey there, Let's see how we can construct the same binary data by using
our code:
```
    h160 := ecc.DecodeBase58("mpNzUycBH6SDU9amLK5raP6Qm71CWNezHv")
    scriptPubKey := tx.P2pkScrit(h160)
    fmt.Printf("raw data for scriptpubkey: %x\n", scriptPubKey.Serialize())
```
The above code will give output like following:
```g
raw data for scriptpubkey: 1976a9146137a18e79a0211946915549b5d155fc75c49b3388ac
```
The first byte 1a is the length of the content, ignore the first byte, the remaining is the same as we have seen in the picture above. Now we add code to construct
the transaction input as following in input.go:
```g
func InitTransactionInput(previousTx []byte, previousIndex *big.Int) *TransactinInput {
	return &TransactinInput{
		previousTransaction: previousTx,
        previousTransactionIdex: previousIndex,
	}
}
```
Then we can construct the transactin input for our current transaction as following:
```g
        prevTxHash, err := hex.DecodeString("703158ce66391f094ab2195cfe5579214073ba90997d0b98e6e410ed1b67aa8a")
	if err != nil {
		panic(err)
	}
	prevTxIndex := big.NewInt(int64(1))
	txInput := tx.InitTransactionInput(prevTxHash, prevTxIndex)
```
Now we need to construct the output, which is used to detail about how many bitcoins will received by whom.

