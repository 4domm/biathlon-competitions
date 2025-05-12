package pkg

import (
	"fmt"
	"io"
)

type OutputWrapper struct {
	OutStream io.Writer
}

func NewOutputWrapper(writer io.Writer) *OutputWrapper {
	return &OutputWrapper{writer}
}

func (o *OutputWrapper) WriteData(data string) {
	fmt.Fprintln(o.OutStream, data)
}
