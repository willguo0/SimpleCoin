package test

import (
	"BrunoCoin/pkg"
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/blockchain"
	"BrunoCoin/pkg/id"
	"BrunoCoin/pkg/proto"
	"BrunoCoin/pkg/utils"
	"BrunoCoin/pkg/wallet"
	"fmt"
	"testing"
	"time"
)

/*
 *  Wallet tests
 */

func TestWalletAdd(t *testing.T) {
	utils.SetDebug(true)

	var transaction1 = tx.Transaction{}
	var transaction2 = tx.Transaction{}
	var transaction3 = tx.Transaction{}
	var transaction4 = tx.Transaction{}

	lt := wallet.NewLmnlTxs(wallet.DefaultConfig())

	lt.Add(nil)

	if lt.TxQ.Len() != 0 {
		t.Errorf("Failed: Heap should have 0 txs but it has " + fmt.Sprint(lt.TxQ.Len()))
	}

	lt.Add(&transaction1)

	if lt.TxQ.Len() != 1 {
		t.Errorf("Failed: Heap should have 1 txs but it has " + fmt.Sprint(lt.TxQ.Len()))
	}

	lt.Add(&transaction2)

	if lt.TxQ.Len() != 2 {
		t.Errorf("Failed: Heap should have 2 txs but it has " + fmt.Sprint(lt.TxQ.Len()))
	}

	lt.Add(&transaction3)
	lt.Add(&transaction4)
	lt.Add(nil)
	lt.Add(nil)

	if lt.TxQ.Len() != 4 {
		t.Errorf("Failed: Heap should have 4 txs but it has " + fmt.Sprint(lt.TxQ.Len()))
	}
}
func TestWalletChkTxs(t *testing.T) {
	utils.SetDebug(true)

	var transaction1 = tx.Transaction{
		LockTime: 5,
	}

	var transaction2 = tx.Transaction{
		LockTime: 10,
	}

	var transaction3 = tx.Transaction{
		LockTime: 15,
	}

	var transaction4 = tx.Transaction{
		LockTime: 20,
	}

	transactions := make([]*tx.Transaction, 0)
	transactions = append(transactions, &transaction1)
	transactions = append(transactions, &transaction2)

	liminal := wallet.NewLmnlTxs(wallet.DefaultConfig())

	liminal.Add(&transaction1)
	liminal.Add(&transaction4)

	liminal.TxQ.Add(12345, &transaction3)

	ans1, ans2 := liminal.ChkTxs(nil)

	if ans1 != nil || ans2 != nil {
		t.Errorf("Failed: Should return nil when checking nil transaction")

	}

	txsAbove, duplicatedtxs := liminal.ChkTxs(transactions)

	if len(txsAbove) != 1 {
		t.Errorf("Failed: There should be exactly one element above the threshold, but there are actually " + fmt.Sprint(len(txsAbove)) + "\n")
	}

	if len(duplicatedtxs) != 1 {
		t.Errorf("Failed: There should be exactly one element removed, but there are actually " + fmt.Sprint(len(duplicatedtxs)) + "\n")
	}

	if duplicatedtxs[0] != &transaction1 {
		t.Errorf("Failed: Duplicated transactions were not calculated properly\n")
	}

	if liminal.TxQ.Len() != 1 {
		t.Errorf("Failed: Heap should have 1 txs but it has " + fmt.Sprint(liminal.TxQ.Len()) + "\n")
	}
}
func TestHndlTxRq(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node := pkg.New(pkg.DefaultConfig(GetFreePort()))

	genNd.Start()
	node.Start()

	genNd.ConnectToPeer(node.Addr)
	genNd.Wallet.HndlTxReq(nil)

	if genNd.Wallet.LmnlTxs.TxQ.Len() != 0 {
		t.Errorf("Failed: Heap should have 0 txs but it has " + fmt.Sprint(genNd.Wallet.LmnlTxs.TxQ.Len()))
	}

	time.Sleep(2 * time.Second)
	genNd.Wallet.HndlTxReq(&wallet.TxReq{Amt: 0})

	if genNd.Wallet.LmnlTxs.TxQ.Len() != 0 {
		t.Errorf("Failed: Heap should have 0 txs but it has " + fmt.Sprint(genNd.Wallet.LmnlTxs.TxQ.Len()))
	}

	genNd.SendTx(20, 50, node.Id.GetPublicKeyBytes()) //Calls wallet.hndltxreq inside
	node.SendTx(100, 50, genNd.Id.GetPublicKeyBytes())

	time.Sleep(4 * time.Second)

	ChkTxSeenLen(t, genNd, 1)
	ChkTxSeenLen(t, node, 1)
	ChkMnChnCons(t, []*pkg.Node{genNd, node})
}

func TestHndlBlk(t *testing.T) {
	utils.SetDebug(true)

	node1 := pkg.New(pkg.DefaultConfig(1))
	block1 := node1.Chain.GetLastBlock()
	wallet1 := node1.Wallet
	len1 := wallet1.LmnlTxs.TxQ.Len()
	wallet1.HndlBlk(nil) //tests nil case

	if len1 != wallet1.LmnlTxs.TxQ.Len() {
		t.Errorf("Failed: Heap should have 0 txs but it has " + fmt.Sprint(wallet1.LmnlTxs.TxQ.Len()))
	}

	wallet1.HndlBlk(block1) //tests case where no transaction passes priority

	if len1 != wallet1.LmnlTxs.TxQ.Len() {
		t.Errorf("Failed: Heap should have 0 txs but it has " + fmt.Sprint(wallet1.LmnlTxs.TxQ.Len()))
	}

	input := make([]*proto.TransactionInput, 0)
	output := make([]*proto.TransactionOutput, 0)

	for i := 0; i < 2; i++ {
		input = append(input, &proto.TransactionInput{Amount: uint32(i)})
		output = append(output, &proto.TransactionOutput{})
	}

	transaction1 := &proto.Transaction{Inputs: input, Outputs: output}

	transaction2 := &proto.Transaction{Inputs: input, Outputs: output}

	wallet1.LmnlTxs.Add(tx.Deserialize(transaction2))
	wallet1.HndlBlk(block1) // none of them are above the limit

	if wallet1.LmnlTxs.TxQ.Len() != 1 {
		t.Errorf("Failed: Heap should have 1 txs but it has " + fmt.Sprint(wallet1.LmnlTxs.TxQ.Len()))

	}

	transactions3 := block1.Transactions[0]
	wallet1.LmnlTxs.Add(transactions3)
	wallet1.HndlBlk(block1) //one of them is in block

	if wallet1.LmnlTxs.TxQ.Len() != 1 {
		t.Errorf("Failed: Heap should have 1 txs but it has " + fmt.Sprint(wallet1.LmnlTxs.TxQ.Len()))

	}

	transactions := make([]*tx.Transaction, 0)
	transactions = append(transactions, tx.Deserialize(transaction1), tx.Deserialize(transaction2))
	block0 := block.New("", transactions, "")
	id1, _ := id.CreateSimpleID()

	w := wallet.New(wallet.DefaultConfig(), id1, blockchain.New(blockchain.DefaultConfig()))
	w.LmnlTxs.TxQ.Add(100, tx.Deserialize(transaction1))

	if w.LmnlTxs.TxQ.Len() != 1 {
		t.Errorf("Failed: Heap should have 1 txs but it has " + fmt.Sprint(wallet1.LmnlTxs.TxQ.Len()))
	}

	w.HndlBlk(block0)

	if w.LmnlTxs.TxQ.Len() != 0 {
		t.Errorf("Failed: Heap should have 1 txs but it has " + fmt.Sprint(wallet1.LmnlTxs.TxQ.Len()))
	}
}
