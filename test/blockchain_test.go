package test

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txo"
	"BrunoCoin/pkg/blockchain"
	"BrunoCoin/pkg/proto"
	"fmt"
	"testing"
)

func TestGetUTXOForAmount(t *testing.T) {
	utxoMap := make(map[string]*txo.TransactionOutput)
	transaction1 := txo.TransactionOutput{
		Amount:        5,
		LockingScript: "trans1",
		Liminal:       false,
	}
	transaction2 := txo.TransactionOutput{
		Amount:        27,
		LockingScript: "trans1",
		Liminal:       true,
	}
	transaction3 := txo.TransactionOutput{
		Amount:        6,
		LockingScript: "trans1",
		Liminal:       false,
	}
	transaction4 := txo.TransactionOutput{
		Amount:        99,
		LockingScript: "trans2",
		Liminal:       false,
	}
	utxoMap[txo.MkTXOLoc(transaction1.Hash(), 0)] = &transaction1
	utxoMap[txo.MkTXOLoc(transaction2.Hash(), 0)] = &transaction2
	utxoMap[txo.MkTXOLoc(transaction3.Hash(), 0)] = &transaction3
	utxoMap[txo.MkTXOLoc(transaction4.Hash(), 0)] = &transaction4

	block1 := blockchain.BlockchainNode{
		Utxo: utxoMap,
	}

	bc := blockchain.Blockchain{
		LastBlock: &block1,
	}
	utxoinfo, change, hasenough := bc.GetUTXOForAmt(0, "trans1")
	if !hasenough {
		t.Errorf("Test Failed: There should be enough utxo")
	}
	if change != 0 {
		t.Errorf("Test Failed: Change should be 0 but it is " + fmt.Sprint(change))
	}
	if len(utxoinfo) != 0 {
		t.Errorf("Test Failed: Length of UTXOinfo list should be 0 but it is " + fmt.Sprint(len(utxoinfo)))
	}
	utxoinfo, change, hasenough = bc.GetUTXOForAmt(8, "trans1")
	if !hasenough {
		t.Errorf("Test Failed: There should be enough utxo")
	}
	if change != 3 {
		t.Errorf("Test Failed: Change should be 3 but it is " + fmt.Sprint(change))
	}
	if len(utxoinfo) != 2 {
		t.Errorf("Test Failed: Length of UTXOinfo list should be 2 but it is " + fmt.Sprint(len(utxoinfo)))
	}

	_, _, hasenough2 := bc.GetUTXOForAmt(100, "trans2")
	if hasenough2 {
		t.Errorf("Test Failed: There shouldn't be enough utxo")
	}
}

func TestBlockChainAdd(t *testing.T) {

	bc := blockchain.New(blockchain.DefaultConfig())
	input := make([]*proto.TransactionInput, 0)
	output := make([]*proto.TransactionOutput, 0)
	for i := 0; i < 2; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(i)})
		output = append(output, &proto.TransactionOutput{})
	}
	transaction1 := &proto.Transaction{Inputs: input,
		Outputs: output}

	transaction2 := &proto.Transaction{Inputs: input,
		Outputs: output}
	txs := make([]*proto.Transaction, 0)
	tx1 := make([]*tx.Transaction, 0)
	tx1 = append(tx1, tx.Deserialize(transaction2), tx.Deserialize(transaction1))
	block2 := block.New(bc.LastBlock.Hash(), tx1, "")
	txs = append(txs, transaction1, transaction2)
	block1 := block.New(bc.LastBlock.Hash(), tx1, "")
	bc.Add(nil)
	bc.Add(block1)
	if bc.Length() != 2 {
		t.Errorf("Blockchain length wrong, block not added properly")
	}
	if bc.GetLastBlock() != block1 {
		t.Errorf("Block not added properly, last block wrong")
	}
	bc.Add(block2)
	if bc.Length() != 2 {
		t.Errorf("Blockchain length wrong, block not added properly")
	}
	block0 := block.New(bc.LastBlock.Hash(), tx1, "")
	bc.Add(block0)
	if bc.Length() != 3 {
		t.Errorf("Blockchain length wrong, block not added properly")
	}
	if bc.GetLastBlock() != block0 {
		t.Errorf("Block not added properly, last block wrong")
	}

}
