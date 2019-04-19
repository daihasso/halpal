package halpal

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/pkg/errors"
)

type Link struct {
    Href string `json:"href"`
}

type HalLinks struct {
    Self *Link `json:"self,omitempty"`
    Next *Link `json:"next,omitempty"`
    Prev *Link `json:"prev,omitempty"`
}

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

type HalItem struct {
    Links *HalLinks `json:"_links,omitempty"`
    Embedded *HalEmbedded `json:"_embedded,omitempty"`
}

func (self HalItem) Embed(key string, value interface{}) error {
    if self.Embedded == nil {
        self.Embedded = new(HalEmbedded)
    }

    return self.Embedded.Set(key, value)
}

type EmbedKeyPair func() (string, interface{})

func EmbeddedItem(key string, val interface{}) EmbedKeyPair {
    return func() (string, interface{}) {
        return key, val
    }
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

func NewHalItem(ctx context.Context) *HalItem {
    return &HalItem{
        Links: NewHalLinks(ctx),
        Embedded: &HalEmbedded{},
    }
}

func NewHalLinks(ctx context.Context, opts ...HalLinksOption) *HalLinks {
    newLinks := HalLinks{}
    if req, ok := ctx.Value("request").(*http.Request); ok {
        newLinks.Self = &Link{
            Href: req.URL.RequestURI(),
        }
    }

    for _, opt := range opts {
        opt(&newLinks)
    }

    return &newLinks
}

type HalLinksOption func(*HalLinks)

func Next(nextLink string) HalLinksOption {
    return func(hl *HalLinks) {
        hl.Next = &Link{
            Href: nextLink,
        }
    }
}

func Prev(prevLink string) HalLinksOption {
    return func(hl *HalLinks) {
        hl.Prev = &Link{
            Href: prevLink,
        }
    }
}
