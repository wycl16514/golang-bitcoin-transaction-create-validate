package main

import (
	"encoding/hex"
	"fmt"
	tx "transaction"
)

/*
1. make sure the total amount in the inputs of transaction is more than
than ouput,
*/

func main() {
	//legacy transaction
	binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	//segwit transaction
	//binaryStr := "01000000000102197393122da5beff963907ff11e4041af10780c868188aad754cc73e3cc35cd9010000001716001462c61a14835b032d5acbe190291d80d0cc5ca28e00000000feae2204104ffe542f30a20012a5b8e2b54a6f61f592520b511801b2237b5ed80100000017160014b30be91e50402cda780c56a3e1c350b1086c80af000000000200a3e111000000001976a914e60c9ac5f72d1d620287a0fc35656bceae5e2ab988ac525d35130000000017a9144795995aff558cc538669ebfecffbe5c9837d5ca870247304402207dd1e7c6c596041276b5285dd3747f586ad819a24acdf0ad60b1faa82af00d3b022046a22dd57df4b72ac165e05b4a6cf8dbecfcfad8f16ae7353df56638ebbf5d1f012103a1a226c5047672af98b2e673751dc69f0140b957753d9c1a789c243100292c6f024730440220670625143c3dfc7a862659a79cbf4ad0f84ff1509bd052cfbfbcdba7adf501f9022015f14a6ee1ae7a8f9fec1070d8a97195422b76a317286c816392cb150d7eb76d012102c910a40bf5726168acc5a8318b0505375e877d4d74448f32ef48156794e657f900000000"
	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}

	transaction := tx.ParseTransaction(binary)
	//fmt.Printf("Fee of the transaction is %v\n", transaction.Fee())
	// script := transaction.GetScript(0, false)
	// //this is not our transaction and we don't have its message and private
	// modifiedTx, err := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000001976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88acfeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac1943060001000000")
	// if err != nil {
	// 	panic(err)
	// }
	// hash256 := ecc.Hash256(string(modifiedTx))
	// fmt.Printf("hash256 of modified transaction is %x\n", hash256)
	// script.Evaluate(hash256)

	res := transaction.Verify()
	fmt.Printf("The evaluation result is %v\n", res)
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
