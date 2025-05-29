package helper

import "strings"

func IsAccessory(productId string) bool {
	return productId == "WIPING-CLOTH" || strings.HasSuffix(productId, "-CLEANNER")
}
