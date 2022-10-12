package core

import (
	"errors"
)

type File struct {
	data *[]byte
}

func NewFile(data *[]byte) (File, error) {
	if data == nil || *data == nil {
		return File{}, errors.New("creating a File struct with nil data is not allowed.")
	}
	return File{data: data}, nil
}

func (f File) IsSet() bool {
	return f.data != nil
}

func (f File) Data() ([]byte, error) {
  err := errors.New("trying to read data from a zero-valued File.")
  if f.data == nil {
    return nil, err
  }
  d := *f.data
  if d == nil {
    return nil, err
  }
  return d, nil
}
