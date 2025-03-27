package crypto

import (
	"errors"
)

var ErrKeyMismatch = errors.New("hashed value is not the hash of the given password")

// Crypter is the common interface implemented by all crypt functions.
type Crypter interface {
	// Generate performs the hashing algorithm, returning a full hash suitable
	// for storage and later password verification.
	//
	// If the salt is empty, a randomly-generated salt will be generated with a
	// length of SaltLenMax and number RoundsDefault of rounds.
	//
	// Any error only can be got when the salt argument is not empty.
	Generate(key, salt []byte) (string, error)

	// Verify compares a hashed key with its possible key equivalent.
	// Returns nil on success, or an error on failure; if the hashed key is
	// different, the error is "ErrKeyMismatch".
	Verify(hashedKey string, key []byte) error

	// Cost returns the hashing cost (in rounds) used to create the given hashed
	// key.
	//
	// When, in the future, the hashing cost of a key needs to be increased in
	// order to adjust for greater computational power, this function allows one
	// to establish which keys need to be updated.
	//
	// The algorithms based in MD5-crypt use a fixed value of rounds.
	Cost(hashedKey string) (int, error)
}
