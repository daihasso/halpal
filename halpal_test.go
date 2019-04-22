package halpal

import (
    "context"
    "encoding/json"
    "net/http"
    "testing"

    gm "github.com/onsi/gomega"
)

type fooObjectTest struct {
    Id int
    Name string
}

var basicHal = []byte(`
{
    "_links": {
        "self": {
            "href": "https://myapp/test?page=2"
        },
        "next": {
            "href": "https://myapp/test?page=3"
        },
        "prev": {
            "href": "https://myapp/test"
        }
    },
    "_embedded": {
        "foos": [
            {
                "id": 5,
                "name": "bar"
            },
            {
                "id": 37,
                "name": "baz"
            }
        ]
    }
}
`)

func TestUnmarshalBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    hal := HalItem{}

    err := json.Unmarshal(basicHal, &hal)
    g.Expect(err).To(gm.BeNil())

    g.Expect(hal.Links.Self.Href).To(gm.Equal("https://myapp/test?page=2"))
    g.Expect(hal.Links.Next.Href).To(gm.Equal("https://myapp/test?page=3"))
    g.Expect(hal.Links.Prev.Href).To(gm.Equal("https://myapp/test"))
    var foos []fooObjectTest
    ok, err := hal.Embedded.Get("foos", &foos)
    g.Expect(err).To(gm.BeNil())
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(foos).To(gm.HaveLen(2))
    g.Expect(foos[0].Id).To(gm.Equal(5))
    g.Expect(foos[0].Name).To(gm.Equal("bar"))
    g.Expect(foos[1].Id).To(gm.Equal(37))
    g.Expect(foos[1].Name).To(gm.Equal("baz"))
}

func TestMissingEmbedded(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    hal := HalItem{}

    err := json.Unmarshal(basicHal, &hal)
    g.Expect(err).To(gm.BeNil())

    var foos []fooObjectTest
    ok, err := hal.Embedded.Get("nonexistant", &foos)
    g.Expect(err).To(gm.BeNil())
    g.Expect(ok).To(gm.BeFalse())
}

func TestBadMatchEmbedded(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    hal := HalItem{}
    var basicHal = []byte(`
    {
        "_embedded": {
            "foos": ["notafoo"]
        }
    }
    `)


    err := json.Unmarshal(basicHal, &hal)
    g.Expect(err).To(gm.BeNil())

    var foos []fooObjectTest
    ok, err := hal.Embedded.Get("foos", &foos)
    g.Expect(err).ToNot(gm.BeNil())
    g.Expect(ok).To(gm.BeFalse())
}

func TestSingleEmbedded(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    hal := HalItem{}
    var basicHal = []byte(`
    {
        "_embedded": {
            "foo": {
                "id": 99,
                "name": "singlefoo"
            }
        }
    }
    `)


    err := json.Unmarshal(basicHal, &hal)
    g.Expect(err).To(gm.BeNil())

    var foo fooObjectTest
    ok, err := hal.Embedded.Get("foo", &foo)

    g.Expect(err).To(gm.BeNil())
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(foo.Id).To(gm.Equal(99))
    g.Expect(foo.Name).To(gm.Equal("singlefoo"))
}

func TestCreateHalItemBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected := `{"_links":{"self":{"href":"/test?page=2"}},` +
        `"_embedded":{"foo":{"Id":76,"Name":"footest"}}}`

    req, err := http.NewRequest("GET", "https://myapp/test?page=2", nil)
    g.Expect(err).To(gm.BeNil())

    t.Log(req.Host)

    ctxWithReq := context.WithValue(context.Background(), "request", req)

    foo := fooObjectTest{
        Id: 76,
        Name: "footest",
    }

    hal, err := NewHalItem(ctxWithReq)
    g.Expect(err).To(gm.BeNil())

    err = hal.Embed("foo", foo)
    g.Expect(err).To(gm.BeNil())

    marshaledHal, err := json.Marshal(hal)
    g.Expect(err).To(gm.BeNil())

    t.Log("Marshaled hal data:", string(marshaledHal))

    g.Expect(marshaledHal).To(gm.MatchJSON(expected))
}

func TestCreateHalItemMultipleEmbedded(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected := `{"_links":{"self":{"href":"/test?page=2"}},"` +
        `_embedded":{"foos":[{"Id":76,"Name":"footest"},` +
        `{"Id":53,"Name":"footest2"}]}}`

    req, err := http.NewRequest("GET", "https://myapp/test?page=2", nil)
    g.Expect(err).To(gm.BeNil())

    t.Log(req.Host)

    ctxWithReq := context.WithValue(context.Background(), "request", req)

    foos := []fooObjectTest{
        fooObjectTest{
            Id: 76,
            Name: "footest",
        },
        fooObjectTest{
            Id: 53,
            Name: "footest2",
        },
    }

    hal, err := NewHalItem(ctxWithReq)
    g.Expect(err).To(gm.BeNil())

    err = hal.Embed("foos", foos)
    g.Expect(err).To(gm.BeNil())

    marshaledHal, err := json.Marshal(hal)
    g.Expect(err).To(gm.BeNil())

    t.Log("Marshaled hal data:", string(marshaledHal))

    g.Expect(marshaledHal).To(gm.MatchJSON(expected))
}

func TestCreateHalItemMultipleEmbeddedBulk(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected := `{"_links":{"self":{"href":"/test?page=2"}},` +
        `"_embedded":{"foo1":[{"Id":76,"Name":"footest"}],` +
        `"foo2":{"Id":53,"Name":"footest2"}}}`

    req, err := http.NewRequest("GET", "https://myapp/test?page=2", nil)
    g.Expect(err).To(gm.BeNil())

    t.Log(req.Host)

    ctxWithReq := context.WithValue(context.Background(), "request", req)

    foo := []fooObjectTest{
        fooObjectTest{
            Id: 76,
            Name: "footest",
        },
    }

    foo2 := fooObjectTest{
        Id: 53,
        Name: "footest2",
    }

    hal, err := NewHalItem(ctxWithReq)
    g.Expect(err).To(gm.BeNil())

    err = hal.EmbedMany(
        EmbeddedItem("foo1", foo),
        EmbeddedItem("foo2", foo2),
    )
    g.Expect(err).To(gm.BeNil())

    marshaledHal, err := json.Marshal(hal)
    g.Expect(err).To(gm.BeNil())

    t.Log("Marshaled hal data:", string(marshaledHal))

    g.Expect(marshaledHal).To(gm.MatchJSON(expected))
}

func TestCreateHalItemEmbeddedInStruct(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected := `{"_links":{"self":{"href":"/test?page=2"}},` +
        `"_embedded":{"foo1":[{"Id":76,"Name":"footest"}]},` +
        `"count":5,"total":6}`

    req, err := http.NewRequest("GET", "https://myapp/test?page=2", nil)
    g.Expect(err).To(gm.BeNil())

    t.Log(req.Host)

    ctxWithReq := context.WithValue(context.Background(), "request", req)

    halStruct := struct{
        HalItem

        Count int `json:"count"`
        Total int `json:"total"`
    }{
        HalItem: *NewHalItemP(ctxWithReq),

        Count: 5,
        Total: 6,
    }

    foo := []fooObjectTest{
        fooObjectTest{
            Id: 76,
            Name: "footest",
        },
    }

    err = halStruct.Embed("foo1", foo)
    g.Expect(err).To(gm.BeNil())

    marshaledHal, err := json.Marshal(halStruct)
    g.Expect(err).To(gm.BeNil())

    t.Log("Marshaled hal data:", string(marshaledHal))

    g.Expect(marshaledHal).To(gm.MatchJSON(expected))
}

func TestCreateHalItemWithExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected := `{"count":5,"total":6}`

    req, err := http.NewRequest("GET", "https://myapp/test?page=2", nil)
    g.Expect(err).To(gm.BeNil())

    t.Log(req.Host)

    hal, err := NewHalItem(context.Background())
    g.Expect(err).To(gm.BeNil())

    hal.AddExtras(ItemExtras{
        "count": 5,
        "total": 6,
    })

    marshaledHal, err := json.Marshal(hal)
    g.Expect(err).To(gm.BeNil())

    t.Log("Marshaled hal data:", string(marshaledHal))

    g.Expect(marshaledHal).To(gm.MatchJSON(expected))
}

func TestAddLinkToHalItem(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected := `{"_links":{"next":{"href":"https://myapp/test?page=3"}}}`

    req, err := http.NewRequest("GET", "https://myapp/test?page=2", nil)
    g.Expect(err).To(gm.BeNil())

    t.Log(req.Host)

    hal, err := NewHalItem(context.Background(), Links(Next("/test?page=3")))
    g.Expect(err).To(gm.BeNil())

    hal.AddLink(Next("https://myapp/test?page=3"))

    marshaledHal, err := json.Marshal(hal)
    g.Expect(err).To(gm.BeNil())

    t.Log("Marshaled hal data:", string(marshaledHal))

    g.Expect(marshaledHal).To(gm.MatchJSON(expected))
}

func TestLinksOptionHalItem(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected := `{"_links":{"self":{"href":"/test?page=2"},` +
        `"next":{"href":"/test?page=3"},"prev":{"href":"/test"}}}`

    req, err := http.NewRequest("GET", "https://myapp/test?page=2", nil)
    g.Expect(err).To(gm.BeNil())

    t.Log(req.Host)

    ctxWithReq := context.WithValue(context.Background(), "request", req)

    hal, err := NewHalItem(
        ctxWithReq, Links(Next("/test?page=3"), Prev("/test")),
    )
    g.Expect(err).To(gm.BeNil())

    t.Logf("Links: %#+v", hal.Links)

    marshaledHal, err := json.Marshal(hal)
    g.Expect(err).To(gm.BeNil())

    t.Log("Marshaled hal data:", string(marshaledHal))

    g.Expect(marshaledHal).To(gm.MatchJSON(expected))
}
