package models

import (
	"errors"
	"fmt"
	"hash/crc32"
	"log/slog"
	"strings"
)

type ApiKeyV1 string

func NewApiKeyV1() string {
	c := NewApiKeyComponentsV1()
	s, _ := c.ToApiKeyV1() // a completely new component V1 will never throw an incomplete error
	return s
}

func (s *ApiKeyV1) ToApiKeyComponentsV1() (*ApiKeyComponentsV1, error) {

	if err := s.Validate(); err != nil {
		return nil, err
	}

	parts := strings.Split(string(*s), separator)

	return &ApiKeyComponentsV1{
		Type:     ApiKeyType(parts[0]),
		Checksum: parts[1],
		Indent:   parts[2],
		Secret:   parts[3],
	}, nil
}

func (s *ApiKeyV1) Validate() error {
	parts := strings.Split(string(*s), separator)

	// There should be 4 parts in the API key
	if len(parts) != 4 {
		slog.Error("API Keys V1 must have 4 distinct parts")
		return errors.New("incomplete API Key v1")
	}

	// Validate checksum
	if parts[1] != fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s%s%s", parts[0], parts[2], parts[3])))) {
		slog.Error("API Key V1 checksum does not match!")
		return errors.New("invalid API Key V1 checksum")
	}

	// Ensure correct prefix
	if parts[0] != string(ApiV1TypePrefix) {
		return errors.New("mismatching version prefix")
	}

	return nil
}
