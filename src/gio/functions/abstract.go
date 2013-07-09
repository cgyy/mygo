package functions

import (
)

type Function struct {
    // Run runs the function
    Func func(interface {}) interface {}
    // registered path
    Path string
}

var Functions = []*Function{}
