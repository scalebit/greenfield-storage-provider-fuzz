package blocksyncer

import (
	"context"
	"encoding/json"

	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/forbole/juno/v4/database"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/parser"
	"github.com/forbole/juno/v4/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func NewIndexer(codec codec.Codec, proxy node.Node, db database.Database, modules []modules.Module) parser.Indexer {
	return &Impl{
		codec:   codec,
		Node:    proxy,
		DB:      db,
		Modules: modules,
	}
}

type Impl struct {
	Modules []modules.Module
	codec   codec.Codec
	Node    node.Node
	DB      database.Database
}

// ExportBlock accepts a finalized block and persists then inside the database.
// An error is returned if write fails.
func (i *Impl) ExportBlock(block *coretypes.ResultBlock, events *coretypes.ResultBlockResults, txs []*types.Tx, vals *coretypes.ResultValidators) error {
	return nil
}

// HandleEvent accepts the transaction and handles events contained inside the transaction.
func (i *Impl) HandleEvent(ctx context.Context, block *coretypes.ResultBlock, index int, event sdk.Event) {
	for _, module := range i.Modules {
		if eventModule, ok := module.(modules.EventModule); ok {
			err := eventModule.HandleEvent(ctx, block, index, event)
			if err != nil {
				log.Errorw("failed to handle event", "module", module, "event", event, "error", err)
			}
		}
	}
}

// Process fetches a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (i *Impl) Process(height uint64) error {
	log.Debugw("processing block", "height", height)

	block, err := i.Node.Block(int64(height))
	if err != nil {
		log.Errorf("failed to get block from node: %s", err)
		return err
	}

	events, err := i.Node.BlockResults(int64(height))
	if err != nil {
		log.Errorf("failed to get block results from node: %s", err)
		return err
	}

	err = i.ExportEvents(context.Background(), block, events)
	if err != nil {
		log.Errorf("failed to ExportEvents: %s", err)
		return err
	}

	err = i.ExportEpoch(block)
	if err != nil {
		log.Errorf("failed to ExportEpoch: %s", err)
		return err
	}

	return nil
}

// ExportEpoch accept a block result data and persist basic info into db to record current sync progress
func (i *Impl) ExportEpoch(block *coretypes.ResultBlock) error {
	// Save the block
	err := i.DB.SaveEpoch(context.Background(), &models.Epoch{
		ID:          1,
		BlockHeight: block.Block.Height,
		BlockHash:   common.HexToHash(block.BlockID.Hash.String()),
		UpdateTime:  block.Block.Time.Unix(),
	})
	if err != nil {
		log.Errorf("failed to persist block: %s", err)
		return err
	}

	return nil
}

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if write fails.
func (i *Impl) ExportTxs(txs []*types.Tx) error {
	return nil
}

// ExportValidators accepts ResultValidators and persists validators inside the database.
// An error is returned if write fails.
func (i *Impl) ExportValidators(block *coretypes.ResultBlock, vals *coretypes.ResultValidators) error {
	return nil
}

// ExportCommit accepts ResultValidators and persists validator commit signatures inside the database.
// An error is returned if write fails.
func (i *Impl) ExportCommit(block *coretypes.ResultBlock, vals *coretypes.ResultValidators) error {
	return nil
}

// ExportAccounts accepts a slice of transactions and persists accounts inside the database.
// An error is returned if write fails.
func (i *Impl) ExportAccounts(block *coretypes.ResultBlock, txs []*types.Tx) error {
	return nil
}

// ExportEvents accepts a slice of transactions and get events in order to save in database.
func (i *Impl) ExportEvents(ctx context.Context, block *coretypes.ResultBlock, events *coretypes.ResultBlockResults) error {
	// get all events in order from the txs within the block
	for _, tx := range events.TxsResults {
		// handle all events contained inside the transaction
		events := filterEventsByType(tx)
		// call the event handlers
		for idx, event := range events {
			i.HandleEvent(ctx, block, idx, event)
		}
	}
	return nil
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers in the order in which they have been registered.
func (i *Impl) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return nil
}

// HandleBlock accepts block and calls the block handlers.
func (i *Impl) HandleBlock(block *coretypes.ResultBlock, events *coretypes.ResultBlockResults, txs []*types.Tx, vals *coretypes.ResultValidators) {
	for _, module := range i.Modules {
		if blockModule, ok := module.(modules.BlockModule); ok {
			err := blockModule.HandleBlock(block, events, txs, vals)
			if err != nil {
				log.Errorw("failed to handle event", "module", module.Name(), "height", block.Block.Height, "error", err)
			}
		}
	}
}

// HandleTx accepts the transaction and calls the tx handlers.
func (i *Impl) HandleTx(tx *types.Tx) {
	log.Infof("HandleTx")
}

// HandleMessage accepts the transaction and handles messages contained inside the transaction.
func (i *Impl) HandleMessage(index int, msg sdk.Msg, tx *types.Tx) {
	log.Infof("HandleMessage")
}