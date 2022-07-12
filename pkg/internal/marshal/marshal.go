package marshal

import "encoding/json"

func NormalizeDocument(document []byte, normalizers ...Normalizer) ([]byte, error) {
	tmp := make(map[string]interface{}, 0)
	if err := json.Unmarshal(document, &tmp); err != nil {
		return nil, err
	}
	for _, normalizer := range normalizers {
		normalizer(tmp)
	}
	return json.Marshal(tmp)
}

type Normalizer func(map[string]interface{})

func KeyAlias(alias string, aliasFor string) Normalizer {
	return func(m map[string]interface{}) {
		for k, v := range m {
			if k == alias {
				m[aliasFor] = v
				delete(m, k)
			}
		}
	}
}

func Plural(key string) Normalizer {
	return func(m map[string]interface{}) {
		if _, isSlice := m[key].([]interface{}); m[key] != nil && !isSlice {
			m[key] = []interface{}{m[key]}
		}
	}
}

func Unplural(key string) Normalizer {
	return func(m map[string]interface{}) {
		if arr, _ := m[key].([]interface{}); len(arr) == 1 {
			m[key] = arr[0]
		}
	}
}

func PluralValueOrMap(key string) Normalizer {
	return func(m map[string]interface{}) {
		value := m[key]
		if value == nil {
			return
		} else if _, isMap := value.(map[string]interface{}); isMap {
			return
		} else if _, isSlice := value.([]interface{}); !isSlice {
			m[key] = []interface{}{m[key]}
		}
	}
}
