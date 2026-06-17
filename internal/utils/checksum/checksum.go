/*
** FLICK PROJECT, 2026
** flick/internal/utils/checksum/checksum.go
** File description:
** BLAKE3 file integrity helpers shared by the CLI and the API server.
 */

package checksum

import (
	"crypto/subtle"
	"encoding/hex"
	"hash"
	"io"
	"os"

	"lukechampine.com/blake3"
)

// digestSize is the BLAKE3 output length in bytes (256-bit digest). It must stay
// in sync with the web client (web/lib/checksum.ts).
const digestSize = 32

// HexLen is the length of a hex-encoded BLAKE3 digest (64 characters).
const HexLen = digestSize * 2

// IsValidHex: Report whether s is a well-formed hex-encoded BLAKE3 digest, i.e.
// the exact expected length and decodable as hexadecimal.
//
// Params:
// - s (string): The candidate hex digest.
//
// Returns:
// - result1 (bool): True when s is a valid digest string.
func IsValidHex(s string) bool {
	if len(s) != HexLen {
		return false
	}
	if _, err := hex.DecodeString(s); err != nil {
		return false
	}
	return true
}

// New: Build a fresh BLAKE3 hasher producing a 256-bit digest.
//
// Returns:
// - result1 (hash.Hash): A ready-to-use BLAKE3 hasher.
func New() hash.Hash {
	return blake3.New(digestSize, nil)
}

// Sum: Encode the current digest of a hasher as a lowercase hex string.
//
// Params:
// - h (hash.Hash): The hasher to read the digest from.
//
// Returns:
// - result1 (string): The hex-encoded digest.
func Sum(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}

// HashReader: Consume r entirely and return the hex BLAKE3 digest of its bytes.
//
// Params:
// - r (io.Reader): The stream to hash.
//
// Returns:
// - result1 (string): The hex-encoded digest.
// - result2 (error): An error if the stream could not be read.
func HashReader(r io.Reader) (string, error) {
	h := New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return Sum(h), nil
}

// HashFile: Return the hex BLAKE3 digest of the file at path, streaming it so
// large files never get loaded fully in memory.
//
// Params:
// - path (string): The path to the file to hash.
//
// Returns:
// - result1 (string): The hex-encoded digest.
// - result2 (error): An error if the file could not be opened or read.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return HashReader(f)
}

// Equal: Report whether two hex digests match. Length mismatches return false.
//
// Params:
// - a (string): The first hex digest.
// - b (string): The second hex digest.
//
// Returns:
// - result1 (bool): True when both digests are identical.
func Equal(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
