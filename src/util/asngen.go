package util

import (
	"fmt"
	"math/rand"
)

func RandomASN(count int) []string {
	// Generate random ASN numbers
	asn := make([]string, count)
	for i := 0; i < count; i++ {
		asn[i] = fmt.Sprintf("%d", rand.Intn(100000))
	}
	return asn
}
