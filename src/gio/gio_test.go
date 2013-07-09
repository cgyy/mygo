package main

import (
    "testing"
    "io/ioutil"
    "net/http"
)

var vtests = []struct {
    url      string
    expected string
} {
    {"http://localhost:1234/", "404 page not found\n"},
    {"http://localhost:1234/sort", ""},
}

func TestUrl(t *testing.T) {
    c := &http.Client{}
    for _, vt := range vtests {
        res, err := c.Get(vt.url)
        if err != nil {
            t.Fatal(err)
        }
        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            t.Fatal(err)
        }
        if string(body) != vt.expected {
            t.Errorf("GET `%s` expected `%s` but `%s`", vt.url, vt.expected, body)
        }
    }
}
