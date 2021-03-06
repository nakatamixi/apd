package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshaler_Append(t *testing.T) {
	kf := NewKeyFinder("", "")
	t.Run("json", func(t *testing.T) {
		marshaler := JSONMarshaler{}
		t.Run("cant unmarshal", func(t *testing.T) {
			input := `hoge`
			_, err := Append(kf, &marshaler, input)
			require.Error(t, err)
		})
		t.Run("can unmarshal", func(t *testing.T) {
			t.Run("can find keys", func(t *testing.T) {
				input := `
[
  {
    "created_at" : "2020-01-01 00:00:00",
    "updated_at" : "2020-01-01 00:00:01"
  }
]`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `[
  {
    "created_at": "2020-01-01 00:00:00",
    "duration": "1s",
    "updated_at": "2020-01-01 00:00:01"
  }
]`, output)

			})
			t.Run("nested map", func(t *testing.T) {
				input := `
[
  {
    "hoge" : {
      "created_at" : "2020-01-01 00:00:00",
      "updated_at" : "2020-01-01 00:00:01"
    }
  }
]`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `[
  {
    "hoge": {
      "created_at": "2020-01-01 00:00:00",
      "duration": "1s",
      "updated_at": "2020-01-01 00:00:01"
    }
  }
]`, output)
			})
			t.Run("nested list", func(t *testing.T) {
				input := `
[
  {
    "hoge" : [
      {
        "created_at" : "2020-01-01 00:00:00",
        "updated_at" : "2020-01-01 00:00:01"
      }
   ]
  }
]`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `[
  {
    "hoge": [
      {
        "created_at": "2020-01-01 00:00:00",
        "duration": "1s",
        "updated_at": "2020-01-01 00:00:01"
      }
    ]
  }
]`, output)
			})
			t.Run("cant find from key", func(t *testing.T) {
				input := `
[
  {
    "hoge" : "2020-01-01 00:00:00",
    "updated_at" : "2020-01-01 00:00:01"
  }
]`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `[
  {
    "hoge": "2020-01-01 00:00:00",
    "updated_at": "2020-01-01 00:00:01"
  }
]`, output)
			})
			t.Run("cant find to key", func(t *testing.T) {
				input := `
[
  {
    "created_at" : "2020-01-01 00:00:00",
    "hoge" : "2020-01-01 00:00:01"
  }
]`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, "[\n  {\n    \"created_at\": \"2020-01-01 00:00:00\",\n    \"hoge\": \"2020-01-01 00:00:01\"\n  }\n]", output)
			})
			t.Run("not supported format", func(t *testing.T) {
				t.Run("cant parse time", func(t *testing.T) {
					input := `
[
  {
    "created_at" : "hoge",
    "updated_at" : "2020-01-01 00:00:01"
  }
]`
					output, err := Append(kf, &marshaler, input)
					require.NoError(t, err)
					require.Equal(t, `[
  {
    "created_at": "hoge",
    "updated_at": "2020-01-01 00:00:01"
  }
]`, output)
				})
			})
		})
	})
	t.Run("yaml", func(t *testing.T) {
		marshaler := YAMLMarshaler{}
		t.Run("cant unmarshal", func(t *testing.T) {
			input := `hoge`
			_, err := Append(kf, &marshaler, input)
			require.Error(t, err)
		})
		t.Run("can unmarshal", func(t *testing.T) {
			t.Run("can find keys", func(t *testing.T) {
				input := `
-  "created_at":  "2020-01-01 00:00:00"
   "updated_at":  "2020-01-01 00:00:01"
`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `- created_at: "2020-01-01 00:00:00"
  duration: 1s
  updated_at: "2020-01-01 00:00:01"
`, output)
			})
			t.Run("nested map", func(t *testing.T) {
				input := `
- hoge:
    "created_at":  "2020-01-01 00:00:00"
    "updated_at":  "2020-01-01 00:00:01"
`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `- hoge:
    created_at: "2020-01-01 00:00:00"
    duration: 1s
    updated_at: "2020-01-01 00:00:01"
`, output)
			})
			t.Run("nested list", func(t *testing.T) {
				input := `
- hoge:
  - "created_at":  "2020-01-01 00:00:00"
    "updated_at":  "2020-01-01 00:00:01"
`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `- hoge:
  - created_at: "2020-01-01 00:00:00"
    duration: 1s
    updated_at: "2020-01-01 00:00:01"
`, output)
			})
			t.Run("cant find from key", func(t *testing.T) {
				input := `
-  "hoge":  "2020-01-01 00:00:00"
   "updated_at":  "2020-01-01 00:00:01"
`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `- hoge: "2020-01-01 00:00:00"
  updated_at: "2020-01-01 00:00:01"
`, output)
			})
			t.Run("cant find to key", func(t *testing.T) {
				input := `
- "created_at":  "2020-01-01 00:00:00"
  "hoge":  "2020-01-01 00:00:01"
`
				output, err := Append(kf, &marshaler, input)
				require.NoError(t, err)
				require.Equal(t, `- created_at: "2020-01-01 00:00:00"
  hoge: "2020-01-01 00:00:01"
`, output)
			})
			t.Run("not supported format", func(t *testing.T) {
				t.Run("cant parse time", func(t *testing.T) {
					input := `
-  "created_at":  "hoge"
   "updated_at":  "2020-01-01 00:00:01"
`
					output, err := Append(kf, &marshaler, input)
					require.NoError(t, err)
					require.Equal(t, `- created_at: hoge
  updated_at: "2020-01-01 00:00:01"
`, output)
				})
			})
		})
	})
}
