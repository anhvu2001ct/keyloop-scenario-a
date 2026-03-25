package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func String(key string, trimSpace bool) string {
	return get(key, func(val string) (string, error) {
		if trimSpace {
			val = strings.TrimSpace(val)
		}
		return val, nil
	})
}

func StringOrDefault(key, defaultVal string, trimSpace bool) string {
	return getOrDefault(key, defaultVal, func(val string) (string, error) {
		if trimSpace {
			val = strings.TrimSpace(val)
		}
		return val, nil
	})
}

func Int(key string) int {
	return get(key, func(val string) (int, error) { return strconv.Atoi(val) })
}

func IntOrDefault(key string, defaultVal int) int {
	return getOrDefault(key, defaultVal, func(val string) (int, error) {
		return strconv.Atoi(val)
	})
}

func get[T any](key string, converter func(val string) (T, error)) T {
	val, exist := os.LookupEnv(key)
	if !exist {
		panic(fmt.Sprintf("env '%s' not exist", key))
	}

	res, err := converter(val)
	if err != nil {
		panic(fmt.Sprintf("cannot convert env '%s' to type %T", key, res))
	}
	return res
}

func getOrDefault[T any](key string, defaultVal T, converter func(val string) (T, error)) T {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defaultVal
	}

	res, err := converter(val)
	if err != nil {
		panic(fmt.Sprintf("cannot convert env '%s' to type %T", key, res))
	}
	return res
}
