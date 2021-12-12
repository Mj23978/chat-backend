package ep

import (
	"fmt"
	"testing"

	"github.com/thoas/go-funk"
)

func TestMergeMap(t *testing.T) {
	mp1 := map[string]interface{}{
		"1": "one",
		"2": "two",
	}
	mp2 := map[string]interface{}{
		"2": "new two",
		"3": "three",
		"4": "four",
	}
	res := mergeMaps(mp1, mp2)
	expect(t, len(res), 4)
	expect(t, res["1"], "one")
	expect(t, res["2"], "new two")
}

func TestCheckRandon(t *testing.T) {
	rand1 := randomString(15)
	rand2 := randomString(15)
	fmt.Printf("First : %v  - Second : %v\n", rand1, rand2)
	expectRev(t, rand1, rand2)
}

func TestFind(t *testing.T) {
	v := funk.Find([]string{"Hi", "Hello"}, func(ele string) bool {
		return ele == "Salam"
	})
	v2 := funk.Find([]string{"Hi", "Hello"}, func(ele string) bool {
		return ele == "Hi"
	})
	expect(t, v, nil)
	expect(t, v2, "Hi")
}
