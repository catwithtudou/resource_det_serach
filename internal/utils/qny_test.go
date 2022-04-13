package utils

import (
	"fmt"
	"testing"
)

func TestUploadToken(t *testing.T) {
	LocalCacheInit()
	LocalCache.Set([]byte("token"), []byte("MmC_mbt2pelp1PPnwJwGQM8pUQtiZrAHv7G1xh5c:TLz3s0eYu69JSv3I0meZJaZuC0Q=:eyJzY29wZSI6InpoZW5neXVhIiwiZGVhZGxpbmUiOjE2NDk4Mzg5NDN9"))
	LocalCache.Set([]byte("expireAt"), []byte("1649839541"))
	tk := LocalCache.Get(nil, []byte("token"))
	ex := LocalCache.Get(nil, []byte("expireAt"))
	fmt.Println(string(tk))
	fmt.Println(string(ex))
}
