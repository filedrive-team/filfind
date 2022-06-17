package filclient

import (
	"context"
	"encoding/json"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs/go-cid"
	logger "github.com/sirupsen/logrus"
)

type MinerInfo struct {
	Owner address.Address // Must be an ID-address.
}

// StateMinerInfo returns info about the indicated miner
func (fc *FilClient) StateMinerInfo(ctx context.Context, miner address.Address, tsk TipSetKey) (data *MinerInfo, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.StateMinerInfo", miner, tsk)
	if err != nil {
		return
	}

	data = &MinerInfo{}
	err = json.Unmarshal(resp, data)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return data, nil
}

// StateLookupID retrieves the ID address of the given address
func (fc *FilClient) StateLookupID(ctx context.Context, addr address.Address, tsk TipSetKey) (addressId address.Address, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.StateLookupID", addr, tsk)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &addressId)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return
}

// StateAccountKey returns the public key address of the given ID address
func (fc *FilClient) StateAccountKey(ctx context.Context, addressId address.Address, tsk TipSetKey) (addr address.Address, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.StateAccountKey", addressId, tsk)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &addr)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return
}

// WalletVerify takes an address, a signature, and some bytes, and indicates whether the signature is valid.
// The address does not have to be in the wallet.
func (fc *FilClient) WalletVerify(ctx context.Context, addr address.Address, msg []byte, sig *crypto.Signature) (valid bool, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.WalletVerify", addr, msg, sig)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &valid)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return
}

type DealState struct {
	SectorStartEpoch abi.ChainEpoch // -1 if not yet included in proven sector
	LastUpdatedEpoch abi.ChainEpoch // -1 if deal state never updated
	SlashEpoch       abi.ChainEpoch // -1 if deal never slashed
}

type DealProposal struct {
	PieceCID             cid.Cid
	PieceSize            abi.PaddedPieceSize
	VerifiedDeal         bool
	Client               address.Address
	Provider             address.Address
	Label                string
	StartEpoch           abi.ChainEpoch
	EndEpoch             abi.ChainEpoch
	StoragePricePerEpoch abi.TokenAmount
	ProviderCollateral   abi.TokenAmount
	ClientCollateral     abi.TokenAmount
}

type MarketDeal struct {
	Proposal DealProposal
	State    DealState
}

// StateMarketStorageDeal returns information about the indicated deal
func (fc *FilClient) StateMarketStorageDeal(ctx context.Context, dealId abi.DealID, tsk TipSetKey) (deal *MarketDeal, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.StateMarketStorageDeal", dealId, tsk)
	if err != nil {
		return
	}

	deal = &MarketDeal{}
	err = json.Unmarshal(resp, deal)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return
}

type TipSet struct {
	Cids   []cid.Cid
	Height abi.ChainEpoch
}

func (ts *TipSet) Key() TipSetKey {
	if ts == nil {
		return EmptyTSK
	}
	return NewTipSetKey(ts.Cids...)
}

// ChainGetTipSetByHeight looks back for a tipset at the specified epoch.
// If there are no blocks at the specified epoch, a tipset at an earlier epoch
// will be returned.
func (fc *FilClient) ChainGetTipSetByHeight(ctx context.Context, height abi.ChainEpoch, tsk TipSetKey) (ts *TipSet, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.ChainGetTipSetByHeight", height, tsk)
	if err != nil {
		return
	}

	ts = &TipSet{}
	err = json.Unmarshal(resp, ts)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return
}

// ChainHead returns the current head of the chain.
func (fc *FilClient) ChainHead(ctx context.Context) (ts *TipSet, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.ChainHead")
	if err != nil {
		return
	}

	ts = &TipSet{}
	err = json.Unmarshal(resp, ts)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return
}

// StateVerifiedClientStatus returns the data cap for the given address.
// Returns nil if there is no entry in the data cap table for the
// address.
func (fc *FilClient) StateVerifiedClientStatus(ctx context.Context, addr address.Address, tsk TipSetKey) (dataCap *abi.StoragePower, err error) {
	resp, err := fc.rpcCall(ctx, "Filecoin.StateVerifiedClientStatus", addr, tsk)
	if err != nil {
		return
	}

	zero := abi.NewStoragePower(0)
	dataCap = &zero
	if string(resp) == "null" {
		return
	}
	err = json.Unmarshal(resp, dataCap)
	if err != nil {
		logger.WithField("resp", string(resp)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	return
}
