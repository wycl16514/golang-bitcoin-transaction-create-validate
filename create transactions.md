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

func (t *TransactinInput) String() string {
	return fmt.Sprintf("previous transaction: %x\n previous tx index: %x\n",
		t.previousTransaction, t.previousTransactionIdex)
}

func (t *TransactinInput) SetScriptSig(sig *ScriptSig) {
	t.scriptSig = sig
}
```

Now we need to construct the output, which is used to detail about how many bitcoins will received by whom. One thing to be noticed is we need to pay mining fee to miner, the more you pay, the faster they help to put you transaction on the chain,
I have 0.00019756 tBTC now, I try to pay 0.00009756 tBTC as mining fee, how many fee need to pay has not algorithm to calculated, 
it depends on experience and market situation.

After substracting the mining fee, we transfer 0.0001 tBTC back to my account, therefore we need to construct one output with amount 0.0001 * STASHI_PER_BITCOIN,
the value of STASHI_PER_BITCOIN is defined in util.go as following:
```g
const (
	STASHI_PER_BITCOIN = 100000000
)
```
Now let's go to output.go, add init function for TransactionOutput Object:

```g
func InitTransactionOutput(amount *big.Int, script *ScriptSig) *TransactionOutput {
	return &TransactionOutput{
		amount:       amount,
		scriptPubKey: script,
	}
}

func (t *TransactionOutput) String() string {
	return fmt.Sprintf("amount: %v\n scriptPubKey: %x\n", t.amount,
		t.scriptPubKey.Serialize())
}
```
goto transaction.go we add init function for Transaction object:
```g
func InitTransaction(version *big.Int, txInputs []*TransactinInput,
	txOutputs []*TransactionOutput, lockTime *big.Int, testnet bool) *Transaction {
	return &Transaction{
		version:   version,
		txInputs:  txInputs,
		txOutputs: txOutputs,
		lockTime:  lockTime,
		testnet:   testnet,
	}
}

func (t *Transaction) String() string {
	txIns := ""
	for i := 0; i < len(t.txInputs); i++ {
		txIns += t.txInputs[i].String()
		txIns += "\n"
	}

	txOuts := ""
	for i := 0; i < len(t.txOutputs); i++ {
		txOuts += t.txOutputs[i].String()
		txOuts += "\n"
	}

	return fmt.Sprintf("tx: version: %x\n transaction inputs\n:%s\n transaction outputs:\n %s\n, locktime: %x\n",
		t.version, txIns, txOuts, t.lockTime)
}
```

Finally we can construct the whole transaction as following:
```g
func main() {
	//construct transaction input by using the previous transaction id and output index
	prevTxHash, err := hex.DecodeString("703158ce66391f094ab2195cfe5579214073ba90997d0b98e6e410ed1b67aa8a")
	if err != nil {
		panic(err)
	}
	prevTxIndex := big.NewInt(int64(1))
	txInput := tx.InitTransactionInput(prevTxHash, prevTxIndex)

	/*
		construct the TransactionOutput object by setting the transfering amount
		and verify script
	*/
	changeAmount := big.NewInt(int64(0.0001 * tx.STASHI_PER_BITCOIN))
	changeH160 := ecc.DecodeBase58("mpNzUycBH6SDU9amLK5raP6Qm71CWNezHv")
	changeScript := tx.P2pkScrit(changeH160)
	changeOut := tx.InitTransactionOutput(changeAmount, changeScript)
	//Init transaction
	transaction := tx.InitTransaction(big.NewInt(int64(1)), []*tx.TransactinInput{txInput},
		[]*tx.TransactionOutput{changeOut}, big.NewInt(int64(0)), true)
	fmt.Printf("%s\n", transaction)
}
```
Running the above code we can get the following result:
```g
tx: version: 1
 transaction inputs
:previous transaction: 703158ce66391f094ab2195cfe5579214073ba90997d0b98e6e410ed1b67aa8a
 previous tx index: 1


 transaction outputs:
 amount: 10000
 scriptPubKey: 1976a9146137a18e79a0211946915549b5d155fc75c49b3388ac


, locktime: 0
```
The output looks correct, now we need to sign the transaction:
```g
//sign the first transaction because we only have one
	z := transaction.SignHash(0)
	zMsg := new(big.Int)
	zMsg.SetBytes(z)
	der := privateKey.Sign(zMsg).Der()
	//add last byte as hash type
	sig := append(der, byte(tx.SIGHASH_ALL))
	_, sec := pubKey.Sec(true)
	scritSig := tx.InitScriptSig([][]byte{sig, sec})
	txInput.SetScriptSig(scritSig)

```
Now we need to serialize the transaction input raw data, and we have done this in SighHash, the only problem is we replace the scriptSig of the given input, 
if we don't replace it, the return result is the raw data of the transaction, therefore we exract part of the SignHash out as following:
```g
func (t *Transaction) SerializeWithSign(inputIdx int) []byte {
	/*
		construct signature message for the given input,we need to change the given
		scriptsig of the input to the scriptpubkey of previous transaction, and serialize
		the transaction to binary data
	*/
	signBinary := make([]byte, 0)
	signBinary = append(signBinary, BigIntToLittleEndian(t.version, LITTLE_ENDIAN_4_BYTES)...)

	inputCount := big.NewInt(int64(len(t.txInputs)))
	signBinary = append(signBinary, EncodeVarint(inputCount)...)
	//serialize inputs, need to replace the given input scriptsig to
	//previous transaction scriptpubkey
	for i := 0; i < len(t.txInputs); i++ {
		if i == inputIdx {
			//found the given input, replace its scriptsig with the scriptpubkey of previous transaction
			t.txInputs[i].ReplaceWithScriptPubKey(t.testnet)
			signBinary = append(signBinary, t.txInputs[i].Serialize()...)
		} else {
			signBinary = append(signBinary, t.txInputs[i].Serialize()...)
		}
	}
	outputCount := big.NewInt(int64(len(t.txOutputs)))
	signBinary = append(signBinary, EncodeVarint(outputCount)...)
	for i := 0; i < len(t.txOutputs); i++ {
		signBinary = append(signBinary, t.txOutputs[i].Serialize()...)
	}
	signBinary = append(signBinary, BigIntToLittleEndian(t.lockTime, LITTLE_ENDIAN_4_BYTES)...)
	signBinary = append(signBinary,
		BigIntToLittleEndian(big.NewInt(int64(SIGHASH_ALL)), LITTLE_ENDIAN_4_BYTES)...)
	//compute hash256 for the modified transaction binary
	return signBinary
}

func (t *Transaction) SignHash(inputIdx int) []byte {
	signBinary := t.SerializeWithSign(inputIdx)
	//compute hash256 for the modified transaction binary
	h256 := ecc.Hash256(string(signBinary))
	return h256
}
```

When all things done, we call Verify of transaction to make sure it can verify it self, if the verification is success, we then serialize the transaction:
```g
    rawTx := transaction.SerializeWithSign(-1)
    fmt.Printf("Transaction raw data:%x\n", rawTx)

```
Notices that we put -1 to SerializeWithSign, this will prevent it to replace any transaction input script,run the above code we can get the following result:
```g
verify result: true
Transaction raw data:01000000018aaa671bed10e4e6980b7d9990ba7340217955fe5c19b24a091f3966ce583170010000006b483045022100b76200083845186983287805f1c6579c9ba861f3107691d5137f515654c992c5022015c60f592026997e79231adaf694063cd4eb8cd9addec59f5b2d0c694cff550801210326423c1bc88465bd1649c85998affb05c3238955400020f24244f925e4dc22acffffffff0110270000000000001976a9146137a18e79a0211946915549b5d155fc75c49b3388ac0000000001000000

```
Your output of transaction raw data may differ with me because the secret key, let's copy the raw data and goto following site:
https://live.blockcypher.com/btc/pushtx/
remember to select the network as bitcoin testnet,and past the transaction raw data into the eidtor:

![截屏2024-05-17 15 37 49](https://github.com/wycl16514/golang-bitcoin-transaction-create-validate/assets/7506958/2059e6bb-33a3-4dfb-bbba-a676a1f6a117)


then click the button "Broadcast Transaction", if it is success, it will show the following page:

![截屏2024-05-17 15 15 53](https://github.com/wycl16514/golang-bitcoin-transaction-create-validate/assets/7506958/f70f5674-ff63-48ee-b60d-50469b093858)

This means our creation of transaction is correct and it can be broadcast to the netwok, wait for a monent we can find a new transaction happended for the given address:
mpNzUycBH6SDU9amLK5raP6Qm71CWNezHv as following:

![截屏2024-05-17 15 28 43](https://github.com/wycl16514/golang-bitcoin-transaction-create-validate/assets/7506958/614dd351-e849-4c7d-bb8b-a077c0e5ca12)


if you achive this, please congradulate yourself!!
