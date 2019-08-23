package autoipv6ptr

import (
    "testing"
)

func TestRemoveIP6DotArpa(t *testing.T) {
    input := "0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa."
    expected := "0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2"
    result := RemoveIP6DotArpa(input)

    if result != expected {
        t.Fatalf("Input: %s Expected: %s Result: %s", input, expected, result)
    }
}

func TestReverseString(t *testing.T) {
    testcases := map[string]string {
        "hello": "olleh",
        "HelloWorld": "dlroWolleH",
        "20010db8000000000000000000000000": "0000000000000000000000008bd01002" }

    for input, expected := range testcases {
        result := ReverseString(input)

        if result != expected {
            t.Fatalf("Input: %s Expected: %s Result: %s", input, expected, result)
        }
    }
}
