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
	CurPri *atomic.Uint32
	PriLim uint32

	TxQ   *tx.Heap
	Ct    *atomic.Uint32
	Cap   uint32
	mutex sync.Mutex
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
// transaction and multiplying by a factor of 100.
// fees * factor / sz
// TODO
// 1. Calculate priority using above formula
// 2. If priority is 0, return 1
// Tip 1: Remember to do error checking on
// variables that might be nil

// some functions/fields/methods that might be helpful
// let t be a transaction object
// t.Sz()
func CalcPri(t *tx.Transaction) uint32 {
	if t == nil {
		return 0
	}

	priority := t.SumInputs() - t.SumOutputs()
	priority *= 100
	priority /= t.Sz()

	if priority == 0 {
		return 1
	}

	return priority
}

// Add adds a transaction to the transaction pool.
// If the transaction pool is full, the transaction
// will not be added. Otherwise, the cumulative
// priority level is updated, the counter is
// incremented, and the transaction is added to the
// heap.
// TODO
// 1. Don't add if capacity is reached
// 2. Add the transaction to the queue with
// the correct priority
// Tip 1: Remember this method is mutating state
// that multiple go routines concurrently have
// access to
// Tip 2: Remember to do error checking on
// variables that might be nil

// Some functions/methods/fields that might be
// helpful
// tp.mutex.Lock()
// tp.mutex.Unlock()
func (tp *TxPool) Add(t *tx.Transaction) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	if tp == nil || tp.Length() >= tp.Cap || t == nil {
		return
	}

	prio := CalcPri(t)

	tp.TxQ.Add(prio, t)
	tp.Ct.Add(1)
	tp.CurPri.Add(prio)
}

// ChkTxs (CheckTransactions) checks for any duplicate
// transactions in the heap and removes them.
// TODO
// 1. Remove duplicate transactions
// 2. update count and total priority fields
// Tip 1: Remember this method is mutating state
// that multiple go routines concurrently have
// access to
// Tip 2: Remember to do error checking on
// variables that might be nil

// Some functions/methods/fields that might be
// helpful
// tp.mutex.Lock()
// tp.mutex.Unlock()
// tp.TxQ.Rmv(...)
func (tp *TxPool) ChkTxs(remover []*tx.Transaction) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	if tp == nil || remover == nil {
		return
	}

	removed := tp.TxQ.Rmv(remover)

	for _, t := range removed {
		tp.Ct.Sub(1)
		tp.CurPri.Sub(CalcPri(t))
	}
}
