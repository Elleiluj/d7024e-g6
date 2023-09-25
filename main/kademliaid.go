package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

// the static number of bytes in a KademliaID
const IDLength = 20
const maxValueInInt = 256

// type definition of a KademliaID
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// RandomKademliaIDInBucket returns the ID of a node in the range of the specified bucket
func RandomKademliaIDInBucket(currentId *KademliaID, bucketIndex int) *KademliaID {
	newKademliaID := NewKademliaID("0000000000000000000000000000000000000000")
	wholeBytes := (IDLength - 1) - (bucketIndex / 8)
	leftOverBits := bucketIndex % 8

	newKademliaID[wholeBytes] = 1 << leftOverBits
	newKademliaID[wholeBytes] |= uint8(rand.Intn(int(newKademliaID[wholeBytes])))

	for i := IDLength - 1; i >= 0; i-- {
		if i <= wholeBytes {
			newKademliaID[i] |= currentId[i]
		} else {
			newKademliaID[i] = currentId[i] | uint8(rand.Intn(maxValueInInt))
		}
	}
	fmt.Printf("Bucket: %v", bucketIndex)
	fmt.Printf("Distance: %v", currentId.CalcDistance(newKademliaID).String())
	return newKademliaID
}

// RandomKademliaIDInBucket
// func RandomKademliaIDInBucket(currentId *KademliaID, bucketIndex int) *KademliaID {
// 	randomDistance := KademliaID{}
// 	for i := 0; i < IDLength; i++ {
// 		if i < bucketIndex {
// 			randomDistance[i] = 0
// 		} else if i == bucketIndex {
// 			randomDistance[i] = 255
// 		} else {
// 			randomDistance[i] = uint8(rand.Intn(256))
// 		}
// 	}
// 	return currentId.Add(&randomDistance)
// }

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// func (kademliaID KademliaID) Add(otherKademliaID *KademliaID) *KademliaID {
// 	sumIDs := KademliaID{}
// 	overflow := byte(0)
// 	for i := IDLength - 1; i >= 0; i-- {
// 		sum := kademliaID[i] + otherKademliaID[i]
// 		if sum > 255 {
// 			sum = 256 - sum
// 			// overflow = sum - 255
// 			// sum = 255
// 		}
// 		sumIDs[i] = sum
// 	}
// 	return &sumIDs
// }

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
