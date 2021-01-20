package env

import (
	"fmt"
	"os"
	"strconv"
)

func MustInt(env string) int {
	v, err := Int(env)
	if err != nil {
		panic(err)
	}
	return v
}

func Int(env string) (int, error) {
	v, err := String(env)
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func MustString(env string) string {
	v, err := String(env)
	if err != nil {
		panic(err)
	}
	return v
}

func String(env string) (string, error) {
	v := os.Getenv(env)
	if v == "" {
		return "", fmt.Errorf("%s is required", env)
	}
	return v, nil
}
