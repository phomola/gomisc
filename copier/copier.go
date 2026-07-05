package copier

import (
	"errors"
	"reflect"
	"sync"
	"unsafe"
)

var (
	copiers sync.Map
	// ErrInvalidType ...
	ErrInvalidType = errors.New("invalid type")
	// ErrFieldMissing ...
	ErrFieldMissing = errors.New("field missing")
)

type typePair struct {
	dst reflect.Type
	src reflect.Type
}

func getCopier(dst, src reflect.Type) (func(unsafe.Pointer, unsafe.Pointer) error, error) {
	pair := typePair{dst, src}
	if f, ok := copiers.Load(pair); ok {
		return f.(func(unsafe.Pointer, unsafe.Pointer) error), nil
	}
	if dst.Kind() != reflect.Struct || src.Kind() != reflect.Struct {
		return nil, ErrInvalidType
	}
	fieldCopiers := make([]func(unsafe.Pointer, unsafe.Pointer) error, 0, src.NumField())
	for fs := range src.Fields() {
		fd, ok := dst.FieldByName(fs.Name)
		if !ok {
			return nil, ErrFieldMissing
		}
		so := fs.Offset
		do := fd.Offset
		switch {
		case fs.Type.Kind() == reflect.Int && fd.Type.Kind() == reflect.Int:
			fieldCopiers = append(fieldCopiers, func(dst, src unsafe.Pointer) error {
				dst = unsafe.Add(dst, do)
				src = unsafe.Add(src, so)
				*(*int)(dst) = *(*int)(src)
				return nil
			})
		case fs.Type.Kind() == reflect.String && fd.Type.Kind() == reflect.String:
			fieldCopiers = append(fieldCopiers, func(dst, src unsafe.Pointer) error {
				dst = unsafe.Add(dst, do)
				src = unsafe.Add(src, so)
				*(*string)(dst) = *(*string)(src)
				return nil
			})
		}
	}
	f := func(dst, src unsafe.Pointer) error {
		for _, c := range fieldCopiers {
			if err := c(dst, src); err != nil {
				return err
			}
		}
		return nil
	}
	copiers.Store(pair, f)
	return f, nil
}

// Copy ...
func Copy[D, S any](dst *D, src *S) error {
	f, err := getCopier(reflect.TypeFor[D](), reflect.TypeFor[S]())
	if err != nil {
		return err
	}
	return f(unsafe.Pointer(dst), unsafe.Pointer(src))
}
