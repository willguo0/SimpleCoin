package miner

import (
	"BrunoCoin/pkg/block/tx"
	"go.uber.org/atomic"
	"sync"
)

/*
 *  Brown University, CS1951L, Summer 2021
 *  Designed by: Colby Anderson, Parker Ljung
 */

// TxPool represents all the valid transactions
// that the miner can mine.
// CurPri is the current cumulative priority of
// all the transactions.
// PriLim is the cumulative priority threshold
// needed to surpass in order to start mining.
// TxQ is the transaction maximum priority queue
// that the transactions are stored in.
// Ct is the current count of the transactions
// in the pool.
// Cap is the maximum amount of allowed
// transactions to store in the pool.
type TxPool struct {
	CurPri   	*atomic.Uint32
	PriLim 		uint32

	TxQ 		*tx.Heap
	Ct			*atomic.Uint32
	Cap         uint32
	mutex		sync.Mutex
}


// Length returns the count of transactions
// currently in the pool.
// Returns:
// uint32 the count (Ct) of the pool
func (tp *TxPool) Length() uint32 {
	return tp.Ct.Load()
}


// NewTxPool constructs a transaction pool.
func NewTxPool(c *Config) *TxPool {
	return &TxPool{
		CurPri: atomic.NewUint32(0),
		PriLim: c.PriLim,
		TxQ:    tx.NewTxHeap(),
		Ct:     atomic.NewUint32(0),
		Cap:    c.TxPCap,
	}
}


// PriMet (PriorityMet) checks to see
// if the transaction pool has enough
// cumulative priority to start mining.
func (tp *TxPool) PriMet() bool {
	return tp.CurPri.Load() >= tp.PriLim
}


// CalcPri (CalculatePriority) calculates the
// priority of a transaction by dividing the
// fees (inputs - outputs) by the size of the
// transaction.
// TODO
func CalcPri(t *tx.Transaction) uint32 {
	return 0
}


// Add adds a transaction to the transaction pool.
// If the transaction pool is full, the transaction
// will not be added. Otherwise, the cumulative
// priority level is updated, the counter is
// incremented, and the transaction is added to the
// heap.
// TODO
func (tp *TxPool) Add(t *tx.Transaction) {
	return
}


// ChkTxs (CheckTransactions) checks for any duplicate
// transactions in the heap and removes them.
// TODO
func (tp *TxPool) ChkTxs(remover []*tx.Transaction) {
	return
}
