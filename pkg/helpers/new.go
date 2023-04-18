package helpers

import (
	"time"
)

// time.Time isn't really a primitive, but it's only pointer is to Location which would be the same one after the copy
type primitive interface {
	~int | ~int64 | ~float32 | ~float64 | ~string | ~bool | time.Time
}

// New takes any type of primitive and returns a pointer to a copy that is allocated in the heap.
// This is needed to allow setting values (hard coded or struct members) by reference to the api structs.
func New[T primitive](value T) *T {
	return &value
}
