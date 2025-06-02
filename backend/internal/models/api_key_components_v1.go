package models

import (
	"errors"
	"fmt"
	"github.com/lucsky/cuid"
	"hash/crc32"
	"log/slog"
	"strings"
)

type ApiKeyType string

const (
	ApiV1TypePrefix ApiKeyType = "sk" // prefix for V1 keys
	separator       string     = "-"
)

type ApiKeyComponentsV1 struct {
	Type     ApiKeyType // only "sk" as of now
	Checksum string     // CRC32 IEEE hash of fmt.Sprintf("%s%s%s", Type, Indent, Secret)
	Indent   string     // Section of the key that is stored and referenced internally
	Secret   string     // Section of the key that will never be stored internally
}

func NewApiKeyComponentsV1() *ApiKeyComponentsV1 {

	t := ApiV1TypePrefix
	index := cuid.New()
	sensitive := cuid.New()
	hash := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s%s%s", t, index, sensitive)))

	return &ApiKeyComponentsV1{
		Type:     t,
		Checksum: fmt.Sprintf("%08x", hash),
		Indent:   index,
		Secret:   sensitive,
	}
}

func (a *ApiKeyComponentsV1) ToApiKeyV1() (string, error) {
	if a.Secret == "" {
		slog.Error("ApikeyV1 should not be reconstructed from censored API information")
		return "", errors.New("reconstruction of ApiKeyV1")
	}

	fields := []string{
		string(a.Type),
		a.Checksum,
		a.Indent,
		a.Secret,
	}
	result := strings.Join(fields, separator)

	return result, nil
}
