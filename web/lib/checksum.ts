import { blake3 } from "@noble/hashes/blake3.js"
import { bytesToHex } from "@noble/hashes/utils.js"

// BLAKE3 file integrity helpers, the web counterpart of the Go `checksum`
// package (internal/utils/checksum). BLAKE3 is a standard, so a digest produced
// here over some bytes is byte-for-byte identical to the one the CLI or the API
// server produces over the same bytes: a file uploaded from the website verifies
// on download by the CLI, and vice versa.

// DIGEST_SIZE: output length in bytes (256-bit digest). Must stay in sync with
// digestSize in internal/utils/checksum/checksum.go.
const DIGEST_SIZE = 32

// hashBytes: BLAKE3 digest of an in-memory byte slice, as lowercase hex.
export function hashBytes(data: Uint8Array): string {
  return bytesToHex(blake3(data, { dkLen: DIGEST_SIZE }))
}

// hashBlob: BLAKE3 digest of a Blob/File, streamed chunk by chunk so a large
// upload archive is never held in memory twice. Returns lowercase hex.
export async function hashBlob(blob: Blob): Promise<string> {
  const hasher = blake3.create({ dkLen: DIGEST_SIZE })
  const reader = blob.stream().getReader()
  for (;;) {
    const { done, value } = await reader.read()
    if (done) break
    hasher.update(value)
  }
  return bytesToHex(hasher.digest())
}

// equal: Compare two hex digests. Both stacks emit lowercase hex, so a plain
// comparison is enough to confirm a download is intact.
export function equal(a: string, b: string): boolean {
  return a.length === b.length && a === b
}
