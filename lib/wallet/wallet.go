package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"

	"github.com/gelembjuk/democoin/lib"
)

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type WalletBalance struct {
	Total    float64
	Approved float64
	Pending  float64
}

// MakeWallet creates Wallet. It generates new keys pair and assign to the object
func (w *Wallet) MakeWallet() {
	private, public := w.newKeyPair()
	w.PrivateKey = private
	w.PublicKey = public
}

// Returns public key of a wallet
func (w Wallet) GetPublicKey() []byte {
	return w.PublicKey
}

// Reurns private key of a wallet
func (w Wallet) GetPrivateKey() ecdsa.PrivateKey {
	return w.PrivateKey
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() []byte {
	pubKeyHash, _ := lib.HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{lib.Version}, pubKeyHash...)
	checksum := lib.Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := lib.Base58Encode(fullPayload)

	return address
}

// ValidateAddress check if address is valid, has valid format
func (w Wallet) ValidateAddress(address string) bool {
	if len(address) == 0 {
		return false
	}

	pubKeyHash := lib.Base58Decode([]byte(address))

	if len(pubKeyHash) <= lib.AddressChecksumLen {
		return false
	}
	actualChecksum := pubKeyHash[len(pubKeyHash)-lib.AddressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-lib.AddressChecksumLen]
	targetChecksum := lib.Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Generate new key pair to create new wallet
func (w *Wallet) newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
