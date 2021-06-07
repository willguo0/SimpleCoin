package test

import (
	"BrunoCoin/pkg"
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txo"
	"BrunoCoin/pkg/proto"
	"BrunoCoin/pkg/utils"
	"encoding/hex"
	"testing"
	"time"
)

/*
 *  Validation tests
 */

func TestChkBlk(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(0)
	genNd.Start()
	genNd.StartMiner()

	// Two transactions are made and mined to the blockchain
	for i := 0; i < 2; i++ {
		genNd.SendTx(2, 20, node2.Id.GetPublicKeyBytes())
		// Sleep to give time for the transaction to be mined
		time.Sleep(time.Second * 3)
	}

	if genNd.ChkBlk(genNd.Chain.LastBlock.Block) != true {
		t.Errorf("Failed: Valid block on the chain is classified as invalid")
	}

	if genNd.ChkBlk(nil) != false {
		t.Errorf("Failed: Nil block is classified as valid but is actually invalid")
	}

	txInput1 := &proto.TransactionInput{
		Amount:          300,
		TransactionHash: "abcabcabc",
		OutputIndex:     0,
	}

	txOutput1 := &proto.TransactionOutput{
		Amount:        200,
		LockingScript: hex.EncodeToString(genNd.Wallet.Id.GetPublicKeyBytes()),
	}

	scr, err := txo.Deserialize(txOutput1).MkSig(genNd.Wallet.Id)

	if err != nil {
		t.Errorf("Error making signature")
	}

	txInput1.UnlockingScript = scr

	inputList1 := make([]*proto.TransactionInput, 0)
	inputList1 = append(inputList1, txInput1)

	outputList1 := make([]*proto.TransactionOutput, 0)
	outputList1 = append(outputList1, txOutput1)

	inputList2 := make([]*proto.TransactionInput, 0)

	transaction1 := tx.Deserialize(&proto.Transaction{Inputs: inputList1, Outputs: outputList1})
	transaction2 := tx.Deserialize(&proto.Transaction{Inputs: inputList2, Outputs: outputList1})

	transactions2 := make([]*tx.Transaction, 0)
	transactions2 = append(transactions2, transaction2)
	transactions2 = append(transactions2, transaction1)
	block2 := block.New(genNd.Chain.LastBlock.Hash(), transactions2, "")
	genNd.Chain.Add(block2)
	genNd.Chain.LastBlock.Utxo[txo.MkTXOLoc("abcabcabc", 0)] = txo.Deserialize(txOutput1)

	if genNd.ChkBlk(block2) != false {
		t.Errorf("Failed: The block does not satisfy the POW, but the transaction was classified as valid")
	}
}
func TestChkTx(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()

	if genNd.ChkTx(nil) != false {
		t.Errorf("Failed: Nil transaction is classified as valid but is actually invalid")
	}

	txInput1 := &proto.TransactionInput{
		Amount:          300,
		TransactionHash: "abcabcabc",
		OutputIndex:     0,
	}

	txOutput1 := &proto.TransactionOutput{
		Amount:        200,
		LockingScript: hex.EncodeToString(genNd.Wallet.Id.GetPublicKeyBytes()),
	}

	scr, err := txo.Deserialize(txOutput1).MkSig(genNd.Wallet.Id)

	if err != nil {
		t.Errorf("Error making signature")
	}

	txInput1.UnlockingScript = scr

	inputList1 := make([]*proto.TransactionInput, 0)
	inputList1 = append(inputList1, txInput1)

	outputList1 := make([]*proto.TransactionOutput, 0)
	outputList1 = append(outputList1, txOutput1)

	transaction1 := tx.Deserialize(&proto.Transaction{Inputs: inputList1, Outputs: outputList1})

	transactions1 := make([]*tx.Transaction, 0)
	transactions1 = append(transactions1, transaction1)
	block1 := block.New(genNd.Chain.LastBlock.Hash(), transactions1, "")
	genNd.Chain.Add(block1)
	genNd.Chain.LastBlock.Utxo[txo.MkTXOLoc("abcabcabc", 0)] = txo.Deserialize(txOutput1)

	if genNd.ChkTx(transaction1) != true {
		t.Errorf("Failed: Valid transaction is classified as invalid")
	}

	txOutput2 := &proto.TransactionOutput{
		Amount:        400,
		LockingScript: hex.EncodeToString(genNd.Wallet.Id.GetPublicKeyBytes()),
	}

	outputList2 := make([]*proto.TransactionOutput, 0)
	outputList2 = append(outputList2, txOutput2)

	transaction2 := tx.Deserialize(&proto.Transaction{Inputs: inputList1, Outputs: outputList2})

	inputList2 := make([]*proto.TransactionInput, 0)
	outputList3 := make([]*proto.TransactionOutput, 0)

	transaction3 := tx.Deserialize(&proto.Transaction{Inputs: inputList2, Outputs: outputList3})

	transactions2 := make([]*tx.Transaction, 0)
	transactions2 = append(transactions2, transaction2)
	transactions2 = append(transactions2, transaction3)

	block2 := block.New(genNd.Chain.LastBlock.Hash(), transactions2, "")
	genNd.Chain.Add(block2)
	genNd.Chain.LastBlock.Utxo[txo.MkTXOLoc("abcabcabc", 0)] = txo.Deserialize(txOutput2)

	if genNd.ChkTx(transaction2) != false {
		t.Errorf("Failed: Transaction's total output amount is greater than its total input amount, but it is incorrectly classified as valid")
	}

	if genNd.ChkTx(transaction3) != false {
		t.Errorf("Failed: Transaction's input and output are empty, but it is incorrectly classified as valid")
	}

	transaction4 := tx.Deserialize(&proto.Transaction{Inputs: inputList1, Outputs: outputList1})

	transactions3 := make([]*tx.Transaction, 0)
	transactions3 = append(transactions3, transaction4)

	block3 := block.New(genNd.Chain.LastBlock.Hash(), transactions3, "")
	genNd.Chain.Add(block3)
	genNd.Chain.LastBlock.Utxo[txo.MkTXOLoc("abcabcabc", 0)] = txo.Deserialize(txOutput2) //should be txo.Deserialize(txOutput1)

	if genNd.ChkTx(transaction4) != false {
		t.Errorf("Failed: Transaction's UTXO is not correct, but the transaction is incorrectly classified as valid")
	}
}
