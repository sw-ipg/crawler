package ds

type KeyValuePair[TKey any, TValue any] struct {
	Key   TKey
	Value TValue
}

func NewKeyValuePair[TKey any, TValue any](key TKey, val TValue) KeyValuePair[TKey, TValue] {
	return KeyValuePair[TKey, TValue]{
		Key:   key,
		Value: val,
	}
}
