package phpserialize_test

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"

	elliotchance_phpserialize "github.com/elliotchance/phpserialize"
	"github.com/trim21/go-phpserialize"
	"github.com/trim21/go-phpserialize/internal/encoder"
)

func BenchmarkMarshal_type(b *testing.B) {
	for _, data := range testCase {
		data := data
		b.Run(data.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				phpserialize.Marshal(data.Data)
			}
		})
	}
}

func BenchmarkMarshal_ifce(b *testing.B) {
	for _, data := range testCase {
		data := data
		b.Run(data.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				phpserialize.Marshal(data)
			}
		})
	}
}

func BenchmarkMarshal_map_type(b *testing.B) {
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]uint, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(m)
				}
			})
		})
	}
}

func BenchmarkMarshal_map_as_ifce(b *testing.B) {
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]uint, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				var v = struct{ Value any }{m}
				for pb.Next() {
					encoder.Marshal(v)
				}
			})
		})
	}
}

func BenchmarkMarshal_map_with_ifce_value(b *testing.B) {
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make(map[int]any, i)
			for j := 0; j < i; j++ {
				m[j+1] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(m)
				}
			})
		})
	}
}

func BenchmarkMarshal_slice_of_value(b *testing.B) {
	type D struct {
		Value []uint
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]uint, i)
			for j := 0; j < i; j++ {
				m[j] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(D{m})
				}
			})
		})
	}
}

func BenchmarkMarshal_ifce_slice_as_value(b *testing.B) {
	type D struct {
		Value any
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]uint, i)
			for j := 0; j < i; j++ {
				m[j] = uint(j + 2)
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(D{m})
				}
			})
			runtime.KeepAlive(m)
		})
	}
}

func BenchmarkMarshal_slice_of_type(b *testing.B) {
	type D struct {
		Value []User
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]User, i)
			for j := 0; j < i; j++ {
				m[j] = User{ID: uint64(j + 2), Name: "u-" + strconv.Itoa(j+2)}
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(D{m})
				}
			})
			runtime.KeepAlive(m)
		})
	}
}

func BenchmarkMarshal_ifce_slice_of_type(b *testing.B) {
	type D struct {
		Value any
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]User, i)
			for j := 0; j < i; j++ {
				m[j] = User{ID: uint64(j + 2), Name: "u-" + strconv.Itoa(j+2)}
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(D{Value: m})
				}
			})
			runtime.KeepAlive(m)
		})
	}
}

func BenchmarkMarshal_ifce_slice_of_ifce(b *testing.B) {
	type D struct {
		Value any
	}
	for i := 1; i <= 1000; i = i * 10 {
		i := i
		b.Run(fmt.Sprintf("len-%d", i), func(b *testing.B) {
			var m = make([]any, i)
			for j := 0; j < i; j++ {
				m[j] = User{ID: uint64(j + 2), Name: "u-" + strconv.Itoa(j+2)}
			}
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					encoder.Marshal(D{Value: m})
				}
			})
			runtime.KeepAlive(m)
		})
	}
}

func Benchmark_marshal_compare(b *testing.B) {
	type Obj struct {
		V int `php:"v"`
		S int `php:"s"`
	}

	type TestData struct {
		Users []User `php:"users"`
		Obj   Obj    `php:"obj"`
	}

	var data = TestData{
		Users: []User{
			{ID: 1, Name: "sai"},
			{ID: 2, Name: "trim21"},
		},
		Obj: Obj{V: 2, S: 3},
	}

	for i := 0; i < b.N; i++ {
		phpserialize.Marshal(data)
	}
}

func Benchmark_elliotchance_phpserialize_marshal(b *testing.B) {
	var data = map[any]any{
		"users": []map[any]any{
			{"id": 1, "name": "sai"},
			{"id": 2, "name": "trim21"},
		},
		"obj": map[string]int{"v": 2, "s": 3},
	}
	for i := 0; i < b.N; i++ {
		elliotchance_phpserialize.Marshal(data, nil)
	}
}
