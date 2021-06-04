package pkg

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
)

/*
 *  Brown University, CS1951L, Summer 2021
 *  Designed by: Colby Anderson, John Roy,
 *	Parker Ljung
 *
 */

// ChkBlk (CheckBlock) validates a block based on multiple
// conditions.
// To be valid:
// The block must be syntactically (ChkBlkSyn), semantically
// (ChkBlkSem), and configurally (ChkBlkConf) valid.
// Each transaction on the block must be syntactically (ChkTxSyn),
// semantically (ChkTxSem), and configurally (ChkTxConf) valid.
// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// Inputs:
// b *block.Block the block to be checked for validity
// Returns:
// bool True if the block is valid. false
// otherwise
// TODO:
// to be valid

// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// The block's size must be less than or equal to the largest
// allowed block size.
// The block hash must be less than the difficulty target.
// The block's first transaction must be of type Coinbase.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// t.IsCoinbase()
// b.SatisfiesPOW(...)
// n.Conf.MxBlkSz
// b.Sz()
// n.Chain.ChkChainsUTXO(...)
func (n *Node) ChkBlk(b *block.Block) bool {
	if n == nil || b == nil {
		return false
	}

	for i, v := range b.Transactions {
		if i == 0 {
			if !v.IsCoinbase() || len(v.Outputs) == 0 {
				return false
			}

			for _, vo := range v.Outputs {
				if vo.Amount <= 0 {
					return false
				}
			}
		} else {
			if v.IsCoinbase() || !n.ChkTx(v) {
				return false
			}
		}
	}

	return n.Chain.ChkChainsUTXO(b.Transactions, b.Hdr.PrvBlkHsh) && b.Sz() <= n.Conf.MxBlkSz && b.SatisfiesPOW(b.Hdr.DiffTarg)
}

// ChkTx (CheckTransaction) validates a transaction.
// Inputs:
// t *tx.Transaction the transaction to be checked for validity
// Returns:
// bool True if the transaction is syntactically valid. false
// otherwise
// TODO:
// to be valid:

// The transaction's inputs and outputs must not be empty.
// The transaction's output amounts must be larger than 0.
// The sum of the transaction's inputs must be larger
// than the sum of the transaction's outputs.
// The transaction must not double spend any UTXO.
// The unlocking script on each of the transaction's
// inputs must successfully unlock each of the corresponding
// UTXO.
// The transaction must not be larger than the
// maximum allowed block size.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// note: let u be a transaction output object
// n.Conf.MxBlkSz
// t.Sz()
// u.IsUnlckd(...)
// n.Chain.GetUTXO(...)
// n.Chain.IsInvalidInput(...)
// t.SumInputs()
// t.SumOutputs()
func (n *Node) ChkTx(t *tx.Transaction) bool {
	if n == nil || t == nil || len(t.Outputs) == 0 {
		return false
	}

	for _, v := range t.Inputs {
		if n.Chain.IsInvalidInput(v) || !n.Chain.GetUTXO(v).IsUnlckd(v.UnlockingScript) {
			return false
		}
	}

	for _, v := range t.Outputs {
		if v.Amount <= 0 {
			return false
		}
	}

	return t.SumInputs() >= t.SumOutputs() && t.Sz() <= n.Conf.MxBlkSz // check with TAs about whether it should be greater / greater or equal (for sum inputs/outputs)
}
