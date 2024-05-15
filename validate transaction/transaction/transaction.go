package transaction

import (
	"bufio"
	"bytes"
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

const (
	SIGHASH_ALL = 1
)

type Transaction struct {
	version   *big.Int
	txInputs  []*TransactionInput
	txOutputs []*TransactionOutput
	lockTime  *big.Int
	testnet   bool
}

func getInputCount(bufReader *bufio.Reader) *big.Int {
	/*
		if the first byte of input is 0, then witness transaction,
		we need to skip the first two bytes(0x00, 0x01)
	*/
	firstByte, err := bufReader.Peek(1)
	if err != nil {
		panic(err)
	}
	if firstByte[0] == 0x00 {
		//skip the first two bytes
		skipBuf := make([]byte, 2)
		_, err = bufReader.Read(skipBuf)
		if err != nil {
			panic(err)
		}
	}

	count := ReadVarint(bufReader)
	fmt.Printf("input count is: %x\n", count)
	return count
}

func (t *Transaction) SignHash(inputIdx int) []byte {
	/*
		construct signature message for the given input indicate by inputIdx,
		we need to change the given scriptsig with the scriptpubkey from the
		output of previous transaction, and then do hash256 on the binary transaction
		data
	*/
	signBinary := make([]byte, 0)
	signBinary = append(signBinary, BigIntToLittleEndian(t.version, LITTLE_ENDIAN_4_BYTES)...)

	inputCount := big.NewInt(int64(len(t.txInputs)))
	signBinary = append(signBinary, EncodeVarint(inputCount)...)
	/*
		serialize inputs, need to replace the scritSig of the given input
		to scriptPubKey of previous transaction
	*/
	for i := 0; i < len(t.txInputs); i++ {
		if i == inputIdx {
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

	h256 := ecc.Hash256(string(signBinary))
	return h256
}

func (t *Transaction) VerifyInput(inputIndex int) bool {
	verifyScript := t.GetScript(inputIndex, t.testnet)
	z := t.SignHash(inputIndex)
	return verifyScript.Evaluate(z)
}

func (t *Transaction) Verify() bool {
	/*
		1. verify fee
		2. verify each transaction input
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

func ParseTransaction(binary []byte) *Transaction {
	transaction := &Transaction{}
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	verBuf := make([]byte, 4)
	bufReader.Read(verBuf)

	version := LittleEndianToBigInt(verBuf, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("transaction version:%x\n", version)
	transaction.version = version

	inputs := getInputCount(bufReader)
	transactionInputs := []*TransactionInput{}
	for i := 0; i < int(inputs.Int64()); i++ {
		input := NewTractionInput(bufReader)
		transactionInputs = append(transactionInputs, input)
	}
	transaction.txInputs = transactionInputs

	//read output counts
	outputs := ReadVarint(bufReader)
	transactionOutputs := []*TransactionOutput{}
	for i := 0; i < int(outputs.Int64()); i++ {
		output := NewTractionOutput(bufReader)
		transactionOutputs = append(transactionOutputs, output)
	}
	transaction.txOutputs = transactionOutputs

	//get last four bytes for lock time
	lockTimeBytes := make([]byte, 4)
	bufReader.Read(lockTimeBytes)
	transaction.lockTime = LittleEndianToBigInt(lockTimeBytes, LITTLE_ENDIAN_4_BYTES)

	return transaction

}

func (t *Transaction) GetScript(idx int, testnet bool) *ScriptSig {
	if idx < 0 || idx >= len(t.txInputs) {
		panic("invalid idx for transaction input")
	}

	txInput := t.txInputs[idx]
	return txInput.Script(testnet)
}

func (t *Transaction) Fee() *big.Int {
	//amount of input - amount of ouptput > 0
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
