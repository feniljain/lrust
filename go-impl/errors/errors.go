package errors

import "errors"

var (
	NoNodeWithGivenData = errors.New("No node found with given data")
)
