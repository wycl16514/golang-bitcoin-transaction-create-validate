package main

import (
	ecc "elliptic_curve"
	"encoding/hex"
	"fmt"
	"math/big"
	tx "transaction"
)

/*
1. make sure the total amount in the inputs of transaction is more than
than ouput,
*/

func main() {
	p := new(big.Int)
	h256 := ecc.Hash256("chenyi1982")
	fmt.Printf("h256: %x\n", h256)
	p.SetBytes(tx.ReverseByteSlice(h256))
	fmt.Printf("p is %x\n", p)
	privateKey := ecc.NewPrivateKey(p)
	pubKey := privateKey.GetPublicKey()

	prevTxHash, err := hex.DecodeString("703158ce66391f094ab2195cfe5579214073ba90997d0b98e6e410ed1b67aa8a")
	if err != nil {
		panic(err)
	}
	prevTxIndex := big.NewInt(int64(1))
	txInput := tx.InitTransactionInput(prevTxHash, prevTxIndex)

	/*
		0.00019756 btc
		send back 0.0001 to myself, and set 0.00009756 as fee to miners
	*/
	changeAmount := big.NewInt(int64(0.0001 * tx.STASHI_PER_BITCOIN))
	changeH160 := ecc.DecodeBase58("mpNzUycBH6SDU9amLK5raP6Qm71CWNezHv")
	changeScript := tx.P2pkScript(changeH160)
	changeOut := tx.InitTransactionOutput(changeAmount, changeScript)

	transaction := tx.InitTransaction(big.NewInt(int64(1)), []*tx.TransactionInput{txInput},
		[]*tx.TransactionOutput{changeOut}, big.NewInt(int64(0)), true)

	fmt.Printf("%s\n", transaction)

	//sign the first transaction
	z := transaction.SignHash(0)
	zMsg := new(big.Int)
	zMsg.SetBytes(z)
	der := privateKey.Sign(zMsg).Der()
	//add the last byte as hash type
	sig := append(der, byte(tx.SIGHASH_ALL))
	_, sec := pubKey.Sec(true)
	scriptSig := tx.InitScriptSig([][]byte{sig, sec})
	txInput.SetScriptSig(scriptSig)

	rawTx := transaction.SerializeWithSign(-1)
	fmt.Printf("raw tx: %x\n", rawTx)
}

/*
1. find the scriptsig for the current input

2. replace the scriptsig data with 00

3. use the scriptpubkey from previous transaction to replace the 00

4. append hash type to the end of the transaction binary data
hash type is 4 byte in little endian format

SIGHASH_ALL 1 => 01 00 00 00

5. Do hash256 on the modified binary data

=> signature message

0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d100000000


	1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac


feffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88
ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac1943060001000000
*/
