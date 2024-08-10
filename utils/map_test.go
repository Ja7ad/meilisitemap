package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var nestedMapForTest = map[string]any{
	"foo": map[string]any{
		"bar": map[string]any{
			"x": map[string]any{
				"y": map[string]any{
					"z": "found me",
					"a": []any{
						map[string]any{
							"title":  "Calligraphy",
							"domain": "art",
						},
						map[string]any{
							"title":  "Paint",
							"domain": "art",
						},
					},
				},
			},
		},
	},
}

func Test_PickByNestedKey(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]any
		key      string
		expected any
	}{
		{
			name:     "Full nested key",
			inputMap: nestedMapForTest,
			key:      "foo.bar.x.y.z",
			expected: "found me",
		},
		{
			name:     "Get Calligraphy",
			inputMap: nestedMapForTest,
			key:      "foo.bar.x.y.a.0.title",
			expected: "Calligraphy",
		},
		{
			name:     "Top level key",
			inputMap: nestedMapForTest,
			key:      "foo",
			expected: map[string]any{"bar": map[string]any{"x": map[string]any{"y": map[string]any{"z": "found me"}}}},
		},
		{
			name:     "Mid-level key",
			inputMap: nestedMapForTest,
			key:      "foo.bar",
			expected: map[string]any{"x": map[string]any{"y": map[string]any{"z": "found me"}}},
		},
		{
			name:     "Non-existent key path",
			inputMap: nestedMapForTest,
			key:      "foo,bar,x,y,z,a,b,c",
			expected: nil,
		},
		{
			name:     "Invalid key with leading dot",
			inputMap: nestedMapForTest,
			key:      ".foo",
			expected: nil,
		},
		{
			name:     "Empty key",
			inputMap: nestedMapForTest,
			key:      "",
			expected: nil,
		},
		{
			name:     "Key not present",
			inputMap: nestedMapForTest,
			key:      "baz",
			expected: nil,
		},
		{
			name:     "Deep nested but wrong",
			inputMap: nestedMapForTest,
			key:      "foo.bar.x.y.z.b",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PickByNestedKey(tt.inputMap, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Benchmark_PickByNestedKey(b *testing.B) {
	keys := []string{
		"foo.bar.x.y.z",
		"foo.bar.x.y.a.0.title",
		"foo.bar.x.y.a.0.domain",
		"foo.bar.x.y.a.1.title",
		"foo.bar.x.y.a.1.domain",
		"foo.bar.x.y",
		"foo.bar",
		"foo",
	}

	for _, key := range keys {
		b.Run(key, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				v := PickByNestedKey(nestedMapForTest, key)
				if v == nil {
					b.Fatal(v)
				}
			}
		})
	}
}
