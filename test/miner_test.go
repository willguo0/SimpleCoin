package test

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/blockchain"
	"BrunoCoin/pkg/id"
	"BrunoCoin/pkg/miner"
	"BrunoCoin/pkg/proto"
	"BrunoCoin/pkg/utils"
	"fmt"
	"go.uber.org/atomic"
	"testing"
)

func TestGenCBTx(t *testing.T) {
	utils.SetDebug(true)
	id1, _ := id.New(id.DefaultConfig())
	miner1 := miner.New(miner.DefaultConfig(0), id1)
	transactions := make([]*tx.Transaction, 0)

	if miner1.GenCBTx(transactions) != nil {
		t.Errorf("Test errored when calculating empty transaction" +
			".\n")
	}
	if miner1.GenCBTx(nil) != nil {
		t.Errorf("Test errored when calculating nil transaction list" +
			".\n")
	}
	transactions = append(transactions, nil)
	if miner1.GenCBTx(transactions) != nil {
		t.Errorf("Test errored when there is a nil transaction" +
			".\n")
	}

	transactions = make([]*tx.Transaction, 0) //refresh transactions
	input := make([]*proto.TransactionInput, 0)
	output := make([]*proto.TransactionOutput, 0)
	for i := 0; i < 11; i++ {
		for j := 0; j < 2; j++ {
			input = append(input, &proto.TransactionInput{Amount: uint32(i)})
			output = append(output, &proto.TransactionOutput{})
		}
		transactions = append(transactions, tx.Deserialize(&proto.Transaction{Inputs: input, Outputs: output}))
	}
	finalTransaction := miner1.GenCBTx(transactions)
	if finalTransaction.Sz() != 28 {
		t.Errorf("Test errored when calculating the coinbase transaction" +
			" Differing lengths were found.\n")
	}
	if finalTransaction.SumInputs() != 0 {
		t.Errorf("Test errored when calculating the coinbase transaction" +
			" Differing inputs were found.\n")
	}
	if finalTransaction.SumOutputs() != 450 {
		t.Errorf("Test errored when calculating the coinbase transaction" +
			" Differing outputs were found.\n")
	}
	id2, _ := id.New(id.DefaultConfig())
	miner2 := miner.Miner{
		Conf:        miner.DefaultConfig(0),
		Id:          id2,
		TxP:         miner.NewTxPool(miner.DefaultConfig(0)),
		MiningPool:  []*tx.Transaction{},
		PrvHsh:      blockchain.GenesisBlock(blockchain.DefaultConfig()).Hash(),
		ChnLen:      atomic.NewUint32(200),
		SendBlk:     make(chan *block.Block),
		PoolUpdated: make(chan bool),
		Mining:      atomic.NewBool(false),
		Active:      atomic.NewBool(false),
	} //this miner is mining on something with a long chnlen so there is no more mint reward
	finalTransaction2 := miner2.GenCBTx(transactions)
	if finalTransaction2.SumOutputs() != 440 {
		t.Errorf("Test errored when calculating the coinbase transaction" +
			" Differing outputs were found. Make sure there the mint reward is equal to 0\n")
	}

}
func TestTxPoolAdd(t *testing.T) {
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
	if firstEle == tx.Deserialize(prioGreaterThanOne) {
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
	if pool.TxQ.Len() != 2 {
		t.Errorf("Test errored when adding transaction two" +
			".\n")
	}
	if newFirstEle == tx.Deserialize(prioOne) {
		t.Errorf("Test errored when adding transaction two" +
			" added to wrong place.\n")
	}
	pool.Add(nil)
	if pool.TxQ.Len() != 2 {
		t.Errorf("Test errored when adding nil" +
			".\n")
	}

	pool = miner.NewTxPool(miner.SmallTxPCapConfig(0))
	pool.Add(tx.Deserialize(prioGreaterThanOne))
	pool.Add(tx.Deserialize(prioOne))
	if pool.TxQ.Len() != 1 {
		t.Errorf("Test errored when adding too many transactions" +
			".\n")
	}
	if pool.TxQ.Peek() != tx.Deserialize(prioOne) {
		t.Errorf("Test errored. kept wrong transaction" +
			".\n")
	}

}
func TestHandlBlk(t *testing.T) {
	utils.SetDebug(true)

	return
}
func TestHandlTx(t *testing.T) {
	utils.SetDebug(true)
	id1, _ := id.New(id.DefaultConfig())
	miner1 := miner.New(miner.DefaultConfig(0), id1)
	miner1.HndlTx(nil)
	input := make([]*proto.TransactionInput, 0)
	output := make([]*proto.TransactionOutput, 0)
	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(i)})
		output = append(output, &proto.TransactionOutput{})
	}
	transaction1 := &proto.Transaction{Inputs: input,
		Outputs: output}
	transaction2 := &proto.Transaction{Inputs: input,
		Outputs: output}
	miner1.HndlTx(tx.Deserialize(transaction1))
	if miner1.TxP.Ct.Load() != 1 {
		t.Errorf("Test errored. didn't add transaction" +
			".\n")
	}
	if miner1.TxP.CurPri.Load() != miner.CalcPri(tx.Deserialize(transaction1)) {
		t.Errorf("Test errored. didn't update total priority" +
			".\n")
	}

	val := len(miner1.PoolUpdated)
	if val != 0 {
		t.Errorf("Test errored. sent value to poolupdated when shouldn't have" +
			".\n")
	}
	miner1.StartMiner()

	miner1.HndlTx(tx.Deserialize(transaction2))
	if miner1.TxP.Ct.Load() != 2 {
		t.Errorf("Test errored. didn't add transaction" +
			".\n")
	}

}
func TestMinerChkTxs(t *testing.T) {
	utils.SetDebug(true)
	pool1 := miner.NewTxPool(miner.DefaultConfig(0))
	input := make([]*proto.TransactionInput, 0)
	output := make([]*proto.TransactionOutput, 0)
	remover := make([]*tx.Transaction, 0)
	for i := 0; i < 10; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(i)})
		output = append(output, &proto.TransactionOutput{})
	}
	transaction1 := &proto.Transaction{Inputs: input,
		Outputs: output}
	transaction2 := &proto.Transaction{Inputs: input,
		Outputs: output}
	pool1.Add(tx.Deserialize(transaction1))
	pool1.Add(tx.Deserialize(transaction2))
	pool1.ChkTxs(nil)
	remover = append(remover, tx.Deserialize(transaction1))
	pool1.ChkTxs(remover)
	if pool1.Ct.Load() != uint32(1) {
		t.Errorf("Test errored. didn't remove the duplicates properly/didn't update the length" +
			".\n")
	}
	if pool1.CurPri.Load() != miner.CalcPri(tx.Deserialize(transaction2)) {
		t.Errorf("Test errored. didn't remove the duplicates properly/didn't update the priority" +
			".\n")
	}
	return
}
func TestHndlChkBlks(t *testing.T) {
	utils.SetDebug(true)

	return
}

func TestCalcPri(t *testing.T) {
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
