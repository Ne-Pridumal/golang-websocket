package storage

import "errors"

var (
	entityNotFound      = errors.New("entity not found")
	enitityAlreadyExist = errors.New("entity already exist")
)
