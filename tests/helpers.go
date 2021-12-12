package ep

import (
	"math/rand"
	"path"
	"reflect"
	"time"
)

// Test returns boolean value to indicate that given pattern is valid.
//
// What is it for?
// Internally `emitter` uses `path.Match` function to find matching. But
// as this functionality is optional `Emitter` don't indicate that the
// pattern is invalid. You should check it separately explicitly via
// `Test` function.
func Test(pattern string) bool {
	_, err := path.Match(pattern, "---")
	return err == nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(10000)) - int64(10000))
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func mergeMapsSelect(maps ...map[string]reflect.SelectCase) map[string]reflect.SelectCase {
	result := make(map[string]reflect.SelectCase)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
