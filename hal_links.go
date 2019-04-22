package halpal

import (
    "context"
    "net/http"
)


type Link struct {
    Href string `json:"href"`
}

type HalLinks struct {
    Self *Link `json:"self,omitempty"`
    Next *Link `json:"next,omitempty"`
    Prev *Link `json:"prev,omitempty"`
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
