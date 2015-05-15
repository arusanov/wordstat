package main

import "testing"

func TestSortedKeys(t *testing.T) {
    toSort := map[string]uint64{
        "two":50,
        "three":0,
        "one":100,
    }
    keysShouldBe := []string{"one", "two", "three"}
    keys := sortedKeys(toSort)
    if len(keysShouldBe) != len(keys) {
        t.Error("len should match")
    }

    for i := range keys {
        if keys[i] != keysShouldBe[i] {
            t.Error("expected ",keysShouldBe[i]," got ",keys[i])
        }
    }
}


