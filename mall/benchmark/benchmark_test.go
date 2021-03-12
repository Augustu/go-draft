package benchmark

import (
	"net/http"
	"testing"
)

func BenchmarkClient(b *testing.B) {
	addr := "http://localhost:8888/api/order/get/1"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		resp, err := http.Get(addr)
		if err != nil {
			b.Fail()
		}
		resp.Body.Close()
	}
}
