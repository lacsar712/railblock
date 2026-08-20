package codec

import "hash/crc32"

// CRC32IEEE computes the IEEE CRC-32 checksum over the given bytes.
// This matches the polynomial used by Ethernet and is stored big-endian
// in the final four bytes of a frame.
func CRC32IEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// AppendCRC writes the IEEE CRC-32 of prefix into the last four bytes of dst.
// dst must be at least len(prefix)+4 bytes; only bytes[0:len(prefix)] are checksummed.
func AppendCRC(dst, prefix []byte) {
	sum := CRC32IEEE(prefix)
	dst[len(prefix)+0] = byte(sum)
	dst[len(prefix)+1] = byte(sum >> 8)
	dst[len(prefix)+2] = byte(sum >> 16)
	dst[len(prefix)+3] = byte(sum >> 24)
}

// VerifyCRC compares the trailing CRC in data against a freshly computed checksum.
func VerifyCRC(data []byte, payloadLen int) bool {
	if len(data) < payloadLen+4 {
		return false
	}
	expected := CRC32IEEE(data[:payloadLen])
	actual := uint32(data[payloadLen])<<24 |
		uint32(data[payloadLen+1])<<16 |
		uint32(data[payloadLen+2])<<8 |
		uint32(data[payloadLen+3])
	return expected == actual
}
