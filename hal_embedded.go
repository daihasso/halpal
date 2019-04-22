package halpal

import (
    "encoding/json"

    "github.com/pkg/errors"
)


type HalEmbedded map[string]*json.RawMessage

func (self HalEmbedded) Get(
    key string, marshalTarget interface{},
) (bool, error) {
    if _, ok := self[key]; !ok {
        return false, nil
    }

    err := json.Unmarshal(*self[key], marshalTarget)
    if err != nil {
        return false, errors.Wrap(err, "Error unmarshaling embedded item.")
    }

    return true, nil
}

func (self HalEmbedded) Set(
    key string, value interface{},
) error {
    marshaledValue, err := json.Marshal(value)
    if err != nil {
        return errors.Wrap(
            err, "Error marshaling provided value",
        )
    }

    rawValue := json.RawMessage(marshaledValue)

    self[key] = &rawValue

    return nil
}


type EmbedKeyPair func() (string, interface{})

func EmbeddedItem(key string, val interface{}) EmbedKeyPair {
    return func() (string, interface{}) {
        return key, val
    }
}

