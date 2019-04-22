package halpal

import (
    "context"
    "encoding/json"

    "github.com/pkg/errors"
)

type ItemExtras map[string]interface{}

type HalItemOption func(*HalItem) error

type HalItem struct {
    Links *HalLinks `json:"_links,omitempty"`
    Embedded *HalEmbedded `json:"_embedded,omitempty"`

    extras map[string]interface{}
}

func (self HalItem) Embed(key string, value interface{}) error {
    if self.Embedded == nil {
        self.Embedded = new(HalEmbedded)
    }

    return self.Embedded.Set(key, value)
}

func (self *HalItem) MarshalJSON() ([]byte, error) {
    if self.Links != nil && (*self.Links) != (HalLinks{}) {
        self.extras["_links"] = self.Links
    }
    if self.Embedded != nil && len(*self.Embedded) != 0 {
        self.extras["_embedded"] = self.Embedded
    }

    return json.Marshal(self.extras)
}

func (self HalItem) EmbedMany(items ...EmbedKeyPair) error {
    for i, item := range items {
        key, val := item()
        err := self.Embed(key, val)
        if err != nil {
            return errors.Wrapf(
                err, "Error embedding item #%d", i,
            )
        }
    }

    return nil
}

func (self HalItem) AddLink(linkOptions ...HalLinksOption) {
    if self.Links == nil {
        self.Links = &HalLinks{}
    }

    for _, opt := range linkOptions {
        opt(self.Links)
    }
}

func (self HalItem) AddExtra(key string, val interface{}) {
    self.extras[key] = val
}

func (self HalItem) AddExtras(extras ItemExtras) {
    for key, val := range extras {
        self.AddExtra(key, val)
    }
}

func Links(opts ...HalLinksOption) HalItemOption {
    return func(item *HalItem) error {
        for _, opt := range opts {
            opt(item.Links)
        }

        return nil
    }
}

func NewHalItemP(ctx context.Context, opts ...HalItemOption) *HalItem {
    hal, err := NewHalItem(ctx, opts...)
    if err != nil {
        panic(err)
    }

    return hal
}

func NewHalItem(ctx context.Context, opts ...HalItemOption) (*HalItem, error) {
    hal := &HalItem{
        Links: NewHalLinks(ctx),
        Embedded: &HalEmbedded{},

        extras: make(map[string]interface{}),
    }

    for i, opt := range opts {
        err := opt(hal)
        if err != nil {
            return nil, errors.Wrapf(err, "Error processing option %d", i)
        }
    }

    return hal, nil
}
