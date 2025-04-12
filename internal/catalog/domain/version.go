package domain

import (
	"errors"
)

var ErrVersionConflict = errors.New("version conflict: aggregate has been modified concurrently")

type Version int
