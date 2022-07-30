package types

import (
	"errors"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtype "github.com/tendermint/tendermint/types"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

type NodeID string

// ProtocolVersion contains the protocol versions for the software.
type ProtocolVersion struct {
	P2P   uint64 `json:"p2p"`
	Block uint64 `json:"block"`
	App   uint64 `json:"app"`
}

//-------------------------------------------------------------

// NodeInfo is the basic node information exchanged
// between two peers during the Tendermint P2P handshake.
type NodeInfo struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`

	// Authenticate
	NodeID     NodeID `json:"id"`          // authenticated identifier
	ListenAddr string `json:"listen_addr"` // accepting incoming

	// Check compatibility.
	// Channels are HexBytes so easier to read as JSON
	Network  string         `json:"network"`  // network/chain ID
	Version  string         `json:"version"`  // major.minor.revision
	Channels bytes.HexBytes `json:"channels"` // channels this node knows about

	// ASCIIText fields
	Moniker string        `json:"moniker"` // arbitrary moniker
	Other   NodeInfoOther `json:"other"`   // other application specific data
}

// NodeInfoOther is the misc. applcation specific data
type NodeInfoOther struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

// ResultBlockchainInfo return List of block metadata
type ResultBlockchainInfo struct {
	LastHeight int64               `json:"last_height"`
	BlockMetas []*tmtype.BlockMeta `json:"block_metas"`
}

type ResultBlock struct {
	BlockID tmtype.BlockID `json:"block_id"`
	Block   *tmtype.Block  `json:"block"`
}

// ResultHeader represents the response for a Header RPC Client query
type ResultHeader struct {
	Header *tmtype.Header `json:"header"`
}

// ResultCommit is Commit and Header
type ResultCommit struct {
	tmtype.SignedHeader `json:"signed_header"`
	CanonicalCommit     bool `json:"canonical"`
}

// ResultBlockResults return ABCI results from a block
type ResultBlockResults struct {
	Height                int64                     `json:"height"`
	TxsResults            []*abci.ResponseDeliverTx `json:"txs_results"`
	TotalGasUsed          int64                     `json:"total_gas_used"`
	BeginBlockEvents      []abci.Event              `json:"begin_block_events"`
	EndBlockEvents        []abci.Event              `json:"end_block_events"`
	ValidatorUpdates      []abci.ValidatorUpdate    `json:"validator_updates"`
	ConsensusParamUpdates *tmproto.ConsensusParams  `json:"consensus_param_updates"`
}

// SyncInfo is a info about the node's syncing state
type SyncInfo struct {
	LatestBlockHash   bytes.HexBytes `json:"latest_block_hash"`
	LatestAppHash     bytes.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight int64          `json:"latest_block_height"`
	LatestBlockTime   time.Time      `json:"latest_block_time"`

	EarliestBlockHash   bytes.HexBytes `json:"earliest_block_hash"`
	EarliestAppHash     bytes.HexBytes `json:"earliest_app_hash"`
	EarliestBlockHeight int64          `json:"earliest_block_height"`
	EarliestBlockTime   time.Time      `json:"earliest_block_time"`

	MaxPeerBlockHeight int64 `json:"max_peer_block_height"`

	CatchingUp bool `json:"catching_up"`

	TotalSyncedTime time.Duration `json:"total_synced_time"`
	RemainingTime   time.Duration `json:"remaining_time"`

	TotalSnapshots      int64         `json:"total_snapshots"`
	ChunkProcessAvgTime time.Duration `json:"chunk_process_avg_time"`
	SnapshotHeight      int64         `json:"snapshot_height"`
	SnapshotChunksCount int64         `json:"snapshot_chunks_count"`
	SnapshotChunksTotal int64         `json:"snapshot_chunks_total"`
	BackFilledBlocks    int64         `json:"backfilled_blocks"`
	BackFillBlocksTotal int64         `json:"backfill_blocks_total"`
}

// ValidatorInfo is a info about the node's validator
type ValidatorInfo struct {
	Address     bytes.HexBytes `json:"address"`
	PubKey      crypto.PubKey  `json:"pub_key"`
	VotingPower int64          `json:"voting_power"`
}

// ResultStatus is Node Status
type ResultStatus struct {
	NodeInfo      NodeInfo      `json:"node_info"`
	SyncInfo      SyncInfo      `json:"sync_info"`
	ValidatorInfo ValidatorInfo `json:"validator_info"`
}

func (s *ResultStatus) TxIndexEnabled() bool {
	if s == nil {
		return false
	}
	return s.NodeInfo.Other.TxIndex == "on"
}

// ResultNetInfo is a info about peer connections
type ResultNetInfo struct {
	Listening bool     `json:"listening"`
	Listeners []string `json:"listeners"`
	NPeers    int      `json:"n_peers"`
	Peers     []Peer   `json:"peers"`
}

// ResultDialSeeds is log from dialing seeds
type ResultDialSeeds struct {
	Log string `json:"log"`
}

// ResultDialPeers is Log from dialing peers
type ResultDialPeers struct {
	Log string `json:"log"`
}

type Peer struct {
	ID  NodeID `json:"node_id"`
	URL string `json:"url"`
}

// ResultValidators is validators for a height.
type ResultValidators struct {
	BlockHeight int64               `json:"block_height"`
	Validators  []*tmtype.Validator `json:"validators"`
	// Count of actual validators in this result
	Count int `json:"count"`
	// Total number of validators
	Total int `json:"total"`
}

// ResultBroadcastTx is the CheckTx result
type ResultBroadcastTx struct {
	Code         uint32         `json:"code"`
	Data         bytes.HexBytes `json:"data"`
	Log          string         `json:"log"`
	Codespace    string         `json:"codespace"`
	MempoolError string         `json:"mempool_error"`

	Hash bytes.HexBytes `json:"hash"`
}

// ResultBroadcastTxCommit is CheckTx and DeliverTx results
type ResultBroadcastTxCommit struct {
	CheckTx   abci.ResponseCheckTx   `json:"check_tx"`
	DeliverTx abci.ResponseDeliverTx `json:"deliver_tx"`
	Hash      bytes.HexBytes         `json:"hash"`
	Height    int64                  `json:"height"`
}

// ResultCheckTx wraps abci.ResponseCheckTx.
type ResultCheckTx struct {
	abci.ResponseCheckTx
}

// ResultTx is a result of querying for a tx
type ResultTx struct {
	Hash     bytes.HexBytes         `json:"hash"`
	Height   int64                  `json:"height"`
	Index    uint32                 `json:"index"`
	TxResult abci.ResponseDeliverTx `json:"tx_result"`
	Tx       tmtype.Tx              `json:"tx"`
	Proof    tmtype.TxProof         `json:"proof,omitempty"`
}

// ResultTxSearch is result of searching for txs
type ResultTxSearch struct {
	Txs        []*ResultTx `json:"txs"`
	TotalCount int         `json:"total_count"`
}

// ResultBlockSearch defines the RPC response type for a block search by events.
type ResultBlockSearch struct {
	Blocks     []*ResultBlock `json:"blocks"`
	TotalCount int            `json:"total_count"`
}

// ResultUnconfirmedTxs is list of mempool txs
type ResultUnconfirmedTxs struct {
	Count      int         `json:"n_txs"`
	Total      int         `json:"total"`
	TotalBytes int64       `json:"total_bytes"`
	Txs        []tmtype.Tx `json:"txs"`
}

// ResultABCIInfo is a info abci msg
type ResultABCIInfo struct {
	Response abci.ResponseInfo `json:"response"`
}

// ResultABCIQuery is response of querying abci msg
type ResultABCIQuery struct {
	Response abci.ResponseQuery `json:"response"`
}

// ResultBroadcastEvidence is the result of broadcasting evidence
type ResultBroadcastEvidence struct {
	Hash []byte `json:"hash"`
}

type (
	ResultUnsafeFlushMempool struct{}
	ResultUnsafeProfile      struct{}
	ResultSubscribe          struct{}
	ResultUnsubscribe        struct{}
	ResultHealth             struct{}
)

// ResultEvent is the event data from a subscription
type ResultEvent struct {
	Query  string             `json:"query"`
	Data   tmtype.TMEventData `json:"data"`
	Events TmEvents           `json:"events"`
}

type TmEvents map[EventName][]string
type EventName = string
