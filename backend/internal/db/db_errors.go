package db

import "errors"

var (
	DbGenericErr         = errors.New("unable to create db client")
	DbInvalidAddr        = errors.New("invalid host address")
	DbInvalidPort        = errors.New("invalid port")
	DbInvalidMaxConnSize = errors.New("invalid max conn size")
)
