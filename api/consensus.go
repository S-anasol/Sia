package api

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/NebulousLabs/Sia/crypto"
	"github.com/NebulousLabs/Sia/types"
	"github.com/NebulousLabs/Sia/encoding"
	"github.com/NebulousLabs/Sia/modules"


	"github.com/NebulousLabs/bolt"
	"github.com/julienschmidt/httprouter"
)

// ConsensusGET contains general information about the consensus set, with tags
// to support idiomatic json encodings.
type ConsensusGET struct {
	Synced       bool              `json:"synced"`
	Height       types.BlockHeight `json:"height"`
	CurrentBlock types.BlockID     `json:"currentblock"`
	Target       types.Target      `json:"target"`
	Difficulty   types.Currency    `json:"difficulty"`
}

type ConsensusFileContract struct {
	FileSize           uint64                         `json:"filesize"`
	FileMerkleRoot     crypto.Hash                    `json:"filemerkleroot"`
	WindowStart        types.BlockHeight              `json:"windowstart"`
	WindowEnd          types.BlockHeight              `json:"windowend"`
	Payout             types.Currency                 `json:"payout"`
	ValidProofOutputs  map[string]types.SiacoinOutput `json:"validproofoutputs"`
	MissedProofOutputs map[string]types.SiacoinOutput `json:"missedproofoutputs"`
	UnlockHash         types.UnlockHash               `json:"unlockhash"`
	RevisionNumber     uint64                         `json:"revisionnumber"`
}

type ConsensusFileContractRevision struct {
	ParentID          types.FileContractID   `json:"parentid"`
	UnlockConditions  types.UnlockConditions `json:"unlockconditions"`
	NewRevisionNumber uint64                 `json:"newrevisionnumber"`

	NewFileSize           uint64                         `json:"newfilesize"`
	NewFileMerkleRoot     crypto.Hash                    `json:"newfilemerkleroot"`
	NewWindowStart        types.BlockHeight              `json:"newwindowstart"`
	NewWindowEnd          types.BlockHeight              `json:"newwindowend"`
	NewValidProofOutputs  map[string]types.SiacoinOutput `json:"newvalidproofoutputs"`
	NewMissedProofOutputs map[string]types.SiacoinOutput `json:"newmissedproofoutputs"`
	NewUnlockHash         types.UnlockHash               `json:"newunlockhash"`
}

type ConsensusTransaction struct {
	SiacoinInputs         map[string]types.SiacoinInput            `json:"siacoininputs"`
	SiacoinOutputs        map[string]types.SiacoinOutput           `json:"siacoinoutputs"`
	FileContracts         map[string]ConsensusFileContract         `json:"filecontracts"`
	FileContractRevisions map[string]ConsensusFileContractRevision `json:"filecontractrevisions"`
	StorageProofs         map[string]types.StorageProof            `json:"storageproofs"`
	SiafundInputs         map[string]types.SiafundInput            `json:"siafundinputs"`
	SiafundOutputs        map[string]types.SiafundOutput           `json:"siafundoutputs"`
	MinerFees             map[string]types.Currency                `json:"minerfees"`
	ArbitraryData         [][]byte                                 `json:"arbitrarydata"`
	TransactionSignatures map[string]types.TransactionSignature    `json:"transactionsignatures"`
}

// ConsensusBlockGET is the object returned by a GET request to
// /consensus/block.
type ConsensusBlock struct {
    BlockID             types.BlockID                    `json:"id"`
	BlockHeight         types.BlockHeight                `json:"blockheight"`
	BlockHeader         types.BlockHeader                `json:"blockheader"`
	Target              types.Target                     `json:"target"`
	Difficulty          types.Currency                   `json:"difficulty"`
    TotalCoins          types.Currency                   `json:"totalcoins"`
    EstimatedHashrate   types.Currency                   `json:"estimatedhashrate"`

	MinerPayouts map[string]types.SiacoinOutput  `json:"minerpayouts"`
	Transactions map[string]ConsensusTransaction `json:"transactions"`
}

type Scods struct {
	scods []modules.SiacoinOutputDiff  `json:"scods"`
}

// consensusHandler handles the API calls to /consensus.
func (api *API) consensusHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	cbid := api.cs.CurrentBlock().ID()
	currentTarget, _ := api.cs.ChildTarget(cbid)
	WriteJSON(w, ConsensusGET{
		Synced:       api.cs.Synced(),
		Height:       api.cs.Height(),
		CurrentBlock: cbid,
		Target:       currentTarget,
		Difficulty:   currentTarget.Difficulty(),
	})
}

// consensusValidateTransactionsetHandler handles the API calls to
// /consensus/validate/transactionset.
func (api *API) consensusValidateTransactionsetHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var txnset []types.Transaction
	err := json.NewDecoder(req.Body).Decode(&txnset)
	if err != nil {
		WriteError(w, Error{"could not decode transaction set: " + err.Error()}, http.StatusBadRequest)
		return
	}
	_, err = api.cs.TryTransactionSet(txnset)
	if err != nil {
		WriteError(w, Error{"transaction set validation failed: " + err.Error()}, http.StatusBadRequest)
		return
	}
	WriteSuccess(w)
}
