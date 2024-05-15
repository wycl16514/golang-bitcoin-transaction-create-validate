For a bitcoin node, one of its major task is to velidate a transaction, there are several steps to take for it, the first thing is to check the output can match to the transaction. For example if a transaction
is about "jim using 10 dollars to by a cup of coffee with price of 3 dollars", then we need to check :

1, jim really has 10 dollars

2, the amount left after buying the coffee should be 7 dollars

If the transaction is honest, then the input of the transaction(10 dollars) should greater than the output of the transaction(7 dollars), that is when we use the amount of input minus the amount of the output
the result should be positive, if the result is negative, then the transaction is "dishonest" it want to fake money from air. We use following code to compare the input amount and output amont:
```g
func (t *Transaction) Fee() *big.Int {
	inputSum := big.NewInt(int64(0))
	outputSum := big.NewInt(int64(0))

	for i := 0; i < len(t.txInputs); i++ {
		addOp := new(big.Int)
		value := t.txInputs[i].Value(t.testnet)
		inputSum = addOp.Add(inputSum, value)
	}

	for i := 0; i < len(t.txOutputs); i++ {
		addOp := new(big.Int)
		outputSum = addOp.Add(outputSum, t.txOutputs[i].amount)
	}

	opSub := new(big.Int)
	return opSub.Sub(inputSum, outputSum)
}
```
Now we can construct a transaction and check its fee like following:
```g
//legacy transaction
	binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}
	transaction := tx.ParseTransaction(binary)
	fmt.Printf("Fee of transaction is :%v\n", transaction.Fee())
```

Running the above code we can get the following result:
```g
Fee of transaction is :40000
```

This means there are 40000 stashi left after the transaction, and no fake money created by this transaction.

The second thing for validation of transaction is to verify signature, In previous section, we combined scriptpubkey and  scriptsig together,
and we can run it to validate the transaction, but the problem is we don't know the message for the signature, now here we give the ways to
construct it by the following ways:

1,find the scriptsig from the input, take the transaction above, we use { and } to show the part of scriptsig binary data:

0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d100000000
{
6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b027745
7c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed
01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a
}

feffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88
ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600

2. remove data for the scriptsig and change it to a single byte with value 00:

0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d100000000
{
00
}
```
feffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88
ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600

3. As we have seen in previous section, we need to get the scriptpubkey from the output of last transaction as show in following:
   
![bitcoin_script](https://github.com/wycl16514/golang-bitcoin-transaction-create-validate/assets/7506958/26675d48-8900-4113-b5e6-78a817a71493)

we get the scriptpubkey from previous transaction output, the following binary data is the scirptpubkey of from the previous transaction
of our transaction above:

1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac

then we replace the data above to the 00 we put in last step as show in the following:

0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d100000000
{
1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac
}

feffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88
ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600

4. appeend hash type to the end, there are several hash types for different purpose, SIGHASN_ALL used to authorize this input to go with all other inputs and outputs,
, SIGHASH_SINGLE used to authorize the input the specific output, SIGHASH_NONE authorize the input to any output, value for SIGHASH_ALL is 1 and need to be encode as
4 bytes with littel endian format, therefore we append 4 bytes of value 1 with previous raw data as following:

0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d100000000
{
1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac
}

feffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88
ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600 [01000000]

Notices the hash type value is enclose with [ and ] 

We do a hash256 on the aboved modified transaction, the result is the signature message z.

Let's put the above modified transaction raw data with previous construct script for a run:
```g
package main

import (
	ecc "elliptic_curve"
	"encoding/hex"
	"fmt"
	tx "transaction"
)

func main() {
	//legacy transaction
	binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}
	transaction := tx.ParseTransaction(binary)
	fmt.Printf("Fee of transaction is :%v\n", transaction.Fee())
	script := transaction.GetScript(0, false)
	//modified transaction for signature
	modifiedTx, err := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000001976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88acfeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac1943060001000000")
	if err != nil {
		panic("decode modified transaction err")
	}
	hash256 := ecc.Hash256(string(modifiedTx))
	fmt.Printf("hash256 of modified transaction is %x\n", hash256)
	res := script.Evaluate((hash256))
	fmt.Printf("the evaluation result is: %v\n", res)
}
```
Run the above code and you can see the final step that is OP_CHECKSIG can be passed and the evaluate result returns true. In aboved code, we use hand to modify the 
transaction binary data, let's see how to do it by code instead of hand, first we need to make some change to BigIntToLittleEndian in util.go:
```g
func BigIntToLittleEndian(v *big.Int, length LITTLE_ENDIAN_LENGTH) []byte {
	switch length {
	case LITTLE_ENDIAN_2_BYTES:
		bin := make([]byte, 2)
		binary.LittleEndian.PutUint16(bin, uint16(v.Uint64()))
		return bin
	case LITTLE_ENDIAN_4_BYTES:
		bin := make([]byte, 4)
		binary.LittleEndian.PutUint32(bin, uint32(v.Uint64()))
		return bin
	case LITTLE_ENDIAN_8_BYTES:
		bin := make([]byte, 8)
		binary.LittleEndian.PutUint64(bin, v.Uint64())
		return bin
	}

	return nil
}
```
Then we goto input.go, We need to add two methods to Transaction input one for getting script pub key from output of previous tranction, the other is to replace
the scriptSig with scriptPubKey as we metioned in step 3:
```g
func (t *TransactinInput) ScriptPubKey(testnet bool) *ScriptSig {
	tx := t.getPreviousTx(testnet)
	return tx.txOutputs[t.previousTransactionIdex.Int64()].scriptPubKey
}

func (t *TransactinInput) ReplaceWithScriptPubKey(testnet bool) {
	//use scriptpubkey of previous transaction to replace current scriptsig
	tx := t.getPreviousTx(testnet)
	t.scriptSig = tx.txOutputs[t.previousTransactionIdex.Int64()].scriptPubKey
}
```

Then in transaction.go, we serialize the transaction into binary data and replace the scriptSig in given input with the scriptPubKey from output of previous 
transaction:
```g
func (t *Transaction) SignHash(inputIdx int) []byte {
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

	h256 := ecc.Hash256(string(signBinary))
	return h256
}

func (t *Transaction) VerifyInput(inputIndex int) bool {
	// txIn := t.txInputs[inputIndex]
	// scriptPubKey := txIn.ScriptPubKey(t.testnet)
	// verifyScript := txIn.scriptSig.Add(scriptPubKey)
	verifyScript := t.GetScript(inputIndex, t.testnet)
	z := t.SignHash(inputIndex)
	return verifyScript.Evaluate(z)
}

func (t *Transaction) Verify() bool {
	/*
		1. verify fee
		2. verify signature of each input
	*/
	if t.Fee().Cmp(big.NewInt(int64(0))) < 0 {
		return false
	}

	for i := 0; i < len(t.txInputs); i++ {
		if t.VerifyInput(i) != true {
			return false
		}
	}

	return true
}
```

In above code method SignHash is implementing 4 steps above by code, and method VerifyInput constructs the verify script by using the GetScript method of 
TransactionInput, then call SignHash to get the signature message and evaluate the message by using the verify script and return the final result.

Let's goto main.go and use code to call those code we just add:
```g
package main

import (
	"encoding/hex"
	"fmt"
	tx "transaction"
)

func main() {
	//legacy transaction
	binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}
	transaction := tx.ParseTransaction(binary)
	res := transaction.Verify()
	fmt.Printf("the evaluation result is: %v\n", res)
}
```
Running above code will have the following result:
```g
the evaluation result is: true
```
This shows we can verify transaction signature successfully!


