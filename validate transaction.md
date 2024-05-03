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
```
Fee of transaction is :40000
```
This means there are 40000 stashi left after the transaction, and no fake money created by this transaction.
