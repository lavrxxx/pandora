package runtime

import (
	"time"

	"github.com/spacelavr/pandora/pkg/types"
	"github.com/spacelavr/pandora/pkg/utils/crypto/sha256"
)

// Runtime
type Runtime struct {
	blockchain types.Blockchain
	last       int
}

// New returns new runtime
func New() *Runtime {
	r := &Runtime{}

	r.blockchain = types.Blockchain{r.Genesis()}
	r.last = 0

	return r
}

// Genesis returns genesis block
func (r *Runtime) Genesis() *types.Block {
	block := &types.Block{
		Cert:      nil,
		PrevHash:  "",
		Index:     0,
		Timestamp: time.Now().UTC(),
	}

	block.Hash = sha256.Compute(block.String())

	return block
}

// AddBlock add block to blockchain
func (r *Runtime) AddBlock(block *types.Block) {
	r.blockchain = append(r.blockchain, block)
	r.last++
}

// LastBlock returns last blockchain block
func (r *Runtime) LastBlock() *types.Block {
	return r.blockchain[r.last]
}

// Blockchain returns blockchain
func (r *Runtime) Blockchain() types.Blockchain {
	return r.blockchain
}