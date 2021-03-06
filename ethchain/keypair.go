package ethchain

import (
	"github.com/ethereum/eth-go/ethutil"
	"math/big"
)

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte

	// The associated account
	account *Account
	state   *State
}

func NewKeyPairFromValue(val *ethutil.Value) *KeyPair {
	keyPair := &KeyPair{PrivateKey: val.Get(0).Bytes(), PublicKey: val.Get(1).Bytes()}

	return keyPair
}

func (k *KeyPair) Address() []byte {
	return ethutil.Sha3Bin(k.PublicKey[1:])[12:]
}

func (k *KeyPair) Account() *Account {
	if k.account == nil {
		k.account = k.state.GetAccount(k.Address())
	}

	return k.account
}

// Create transaction, creates a new and signed transaction, ready for processing
func (k *KeyPair) CreateTx(receiver []byte, value *big.Int, data []string) *Transaction {
	tx := NewTransaction(receiver, value, data)
	tx.Nonce = k.account.Nonce

	// Sign the transaction with the private key in this key chain
	tx.Sign(k.PrivateKey)

	return tx
}

func (k *KeyPair) RlpEncode() []byte {
	return ethutil.EmptyValue().Append(k.PrivateKey).Append(k.PublicKey).Encode()
}

type KeyRing struct {
	keys []*KeyPair
}

func (k *KeyRing) Add(pair *KeyPair) {
	k.keys = append(k.keys, pair)
}

// The public "singleton" keyring
var keyRing *KeyRing

func GetKeyRing(state *State) *KeyRing {
	if keyRing == nil {
		keyRing = &KeyRing{}

		data, _ := ethutil.Config.Db.Get([]byte("KeyRing"))
		it := ethutil.NewValueFromBytes(data).NewIterator()
		for it.Next() {
			v := it.Value()
			keyRing.Add(NewKeyPairFromValue(v))
		}
	}

	return keyRing
}
