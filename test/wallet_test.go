package test

import (
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/utils"
	"BrunoCoin/pkg/wallet"
	"fmt"
	"testing"
)

// WALLET TESTS //

// Contents of transaction don't matter but the versions are different to make them distinct
var transaction1 = tx.Transaction{
	Version: 1,
}
var transaction2 = tx.Transaction{
	Version: 2,
}
var transaction3 = tx.Transaction{
	Version: 3,
}
var transaction4 = tx.Transaction{
	Version: 4,
}

func TestWalletAdd(t *testing.T) {
	utils.SetDebug(true)
	lt := wallet.NewLmnlTxs(wallet.DefaultConfig())

	lt.Add(nil)
	if lt.TxQ.Len() != 0 {
		t.Errorf("Failed: Length of heap should be 0 but it is " + fmt.Sprint(lt.TxQ.Len()))
	}
	lt.Add(&transaction1)
	if lt.TxQ.Len() != 1 {
		t.Errorf("Failed: Length of heap should be 1 but it is " + fmt.Sprint(lt.TxQ.Len()))
	}
	lt.Add(&transaction2)
	if lt.TxQ.Len() != 2 {
		t.Errorf("Failed: Length of heap should be 2 but it is " + fmt.Sprint(lt.TxQ.Len()))
	}

}
func TestWalletChkTxs(t *testing.T) {
	utils.SetDebug(true)
	transactions := make([]*tx.Transaction, 0)
	transactions = append(transactions, &transaction1)
	transactions = append(transactions, &transaction2)
	liminal := wallet.NewLmnlTxs(wallet.DefaultConfig())
	liminal.Add(&transaction1)
	liminal.Add(&transaction4)
	liminal.TxQ.Add(500, &transaction3)
	ans1, ans2 := liminal.ChkTxs(nil)
	if ans1 != nil || ans2 != nil {
		t.Errorf("Failed: When checking nil transactions, should return nil")

	}
	txsAbove, duplicatedtxs := liminal.ChkTxs(transactions)
	if len(txsAbove) != 1 {
		t.Errorf("Failed: number of elements above the threshold should be 1 but instead is " + fmt.Sprint(len(txsAbove)))
	}
	if len(duplicatedtxs) != 1 {
		t.Errorf("Failed: number of elements removed should be 1 but instead is " + fmt.Sprint(len(duplicatedtxs)))
	}
	if duplicatedtxs[0] != &transaction1 {
		t.Errorf("Failed: The duplicated transaction was returned wrong")
	}
	if liminal.TxQ.Len() != 1 {
		t.Errorf("Failed: Length of heap should be 1 but it is " + fmt.Sprint(liminal.TxQ.Len()))
	}

}
func TestHndlTxRq(t *testing.T) {
	utils.SetDebug(true)

	return
}
func TestHndlBlk(t *testing.T) {
	utils.SetDebug(true)

	//id, _ := id.CreateSimpleID()
	//blockchain1 := blockchain.New(blockchain.DefaultConfig())
	//wallet1 := wallet.New(wallet.DefaultConfig(),id,blockchain1)
	//wallet1.HndlBlk(nil) //tests nil case
	//input := make([]*proto.TransactionInput, 0)
	//output := make([]*proto.TransactionOutput, 0)
	//for i := 0; i < 10; i++ {
	//	input = append(input, &proto.TransactionInput{Amount: uint32(i)})
	//	output = append(output, &proto.TransactionOutput{})
	//}
	//transactions := make([]*tx.Transaction, 0)
	//transaction1 := &proto.Transaction{Inputs: input,
	//	Outputs: output}
	//transaction2 := &proto.Transaction{Inputs: input,
	//	Outputs: output}
	//wallet1.LmnlTxs.Add(tx.Deserialize(transaction2))
	//transactions = append(transactions, tx.Deserialize(transaction1),tx.Deserialize(transaction2))
	//print("1")
	//
	//block1 := block.New(blockchain.GenesisBlock(blockchain.DefaultConfig()).Hash(), nil, "2")
	//wallet1.HndlBlk(block1)//nil transactions
	//print("1")
	//block1 = block.New(blockchain.GenesisBlock(blockchain.DefaultConfig()).Hash(), transactions, "2")
	//wallet1.HndlBlk(block1)//no transactions above priority threshold
	//wallet1.LmnlTxs.TxQ.Add(500, tx.Deserialize(transaction1))
	//wallet1.HndlBlk(block1)

}
