package test

import (
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/miner"
	"BrunoCoin/pkg/proto"
	"BrunoCoin/pkg/utils"
	"fmt"
	"testing"
)


func TestGetUTXOForAmount(t *testing.T){
	return
}
func TestBlockChainAdd(t *testing.T){
	return
}

func TestGenCBTx(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestTxPoolAdd(t *testing.T){
	utils.SetDebug(true)
	pool := miner.NewTxPool(miner.DefaultConfig(0))
	input := make([]*proto.TransactionInput, 0)
	output := make([]*proto.TransactionOutput, 0)
	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(i)})
		output = append(output, &proto.TransactionOutput{})
	}

	prioGreaterThanOne := &proto.Transaction{Inputs: input,
		Outputs: output}
	pool.Add(tx.Deserialize(prioGreaterThanOne))
	firstEle := pool.TxQ.Peek()
	if  firstEle == tx.Deserialize(prioGreaterThanOne){
		t.Errorf("Test errored when adding transaction one" +
			".\n")
	}
	input = make([]*proto.TransactionInput, 0)
	output = make([]*proto.TransactionOutput, 0)
	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{})
		output = append(output, &proto.TransactionOutput{})
	}
	prioOne := &proto.Transaction{Inputs: input,
		Outputs: output}
	newFirstEle := pool.TxQ.Peek()
	pool.Add(tx.Deserialize(prioOne))
	print(pool.TxQ.Len())
	if pool.TxQ.Len() != 2{
		t.Errorf("Test errored when adding transaction two" +
			".\n")
	}
	if newFirstEle == tx.Deserialize(prioOne){
		t.Errorf("Test errored when adding transaction two" +
			" added to wrong place.\n")
	}


	pool = miner.NewTxPool(miner.SmallTxPCapConfig(0))
	input = make([]*proto.TransactionInput, 0)
	output = make([]*proto.TransactionOutput, 0)
	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(i)})
		output = append(output, &proto.TransactionOutput{})
	}

	prioGreaterThanOne = &proto.Transaction{Inputs: input,
		Outputs: output}
	pool.Add(tx.Deserialize(prioGreaterThanOne))
	firstEle = pool.TxQ.Peek()
	if  firstEle == tx.Deserialize(prioGreaterThanOne){
		t.Errorf("Test errored when adding transaction one" +
			".\n")
	}
	input = make([]*proto.TransactionInput, 0)
	output = make([]*proto.TransactionOutput, 0)
	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{})
		output = append(output, &proto.TransactionOutput{})
	}
	prioOne = &proto.Transaction{Inputs: input,
		Outputs: output}
	newFirstEle = pool.TxQ.Peek()
	pool.Add(tx.Deserialize(prioOne))
	print(pool.TxQ.Len())
	if pool.TxQ.Len() != 2{
		t.Errorf("Test errored when adding transaction two" +
			".\n")
	}
	if newFirstEle == tx.Deserialize(prioOne){
		t.Errorf("Test errored when adding transaction two" +
			" added to wrong place.\n")
	}


}
func TestHandlBlk(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestHandlTx(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestCalcPri(t *testing.T){
	utils.SetDebug(true)
	input := make([]*proto.TransactionInput, 0)
	output := make([]*proto.TransactionOutput, 0)

	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{})
		output = append(output, &proto.TransactionOutput{})
	}
	testTransaction := &proto.Transaction{Inputs: input,
		Outputs: output}
	num := miner.CalcPri(tx.Deserialize(testTransaction))
	if num != 1 {
		t.Errorf("Test errored when calculating priority = 1" +
			".\n")
	}
	input = make([]*proto.TransactionInput, 0)
	output = make([]*proto.TransactionOutput, 0)

	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(i)})
		output = append(output, &proto.TransactionOutput{})
	}
	testTransaction = &proto.Transaction{Inputs: input,
		Outputs: output}
	num = miner.CalcPri(tx.Deserialize(testTransaction))
	if num != 7 {
		t.Errorf("Test miscalculated as " + fmt.Sprint(num) +
			" when actually supposed to be 7.\n")
	}
	input = make([]*proto.TransactionInput, 0)
	output = make([]*proto.TransactionOutput, 0)

	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(6)})
		output = append(output, &proto.TransactionOutput{Amount: uint32(i)})
	}
	testTransaction = &proto.Transaction{Inputs: input,
		Outputs: output}
	num = miner.CalcPri(tx.Deserialize(testTransaction))
	if num != 2 {
		t.Errorf("Test miscalculated as " + fmt.Sprint(num) +
			" when actually supposed to be 1.\n")
	}

}
func Test(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestMinerChkTxs(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestHndlChkBlks(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestWalletAdd(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestWalletChkTxs(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestHndlTxRq(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestHndlBlk(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestChkBlk(t *testing.T){
	utils.SetDebug(true)

	return
}
func TestChkTx(t *testing.T){
	utils.SetDebug(true)

	return
}

var t1 = tx.Transaction{

}
var t2 = tx.Transaction{

}
var t3 = tx.Transaction{

}
var t4 = tx.Transaction{

}
