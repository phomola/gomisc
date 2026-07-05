package copier

import (
	"errors"
	"fmt"
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
		if fs.Tag.Get("copier") == "-" {
			continue
		}
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
		case fs.Type.Kind() == reflect.Float64 && fd.Type.Kind() == reflect.Float64:
			fieldCopiers = append(fieldCopiers, func(dst, src unsafe.Pointer) error {
				dst = unsafe.Add(dst, do)
				src = unsafe.Add(src, so)
				*(*float64)(dst) = *(*float64)(src)
				return nil
			})
		case fs.Type.Kind() == reflect.String && fd.Type.Kind() == reflect.String:
			fieldCopiers = append(fieldCopiers, func(dst, src unsafe.Pointer) error {
				dst = unsafe.Add(dst, do)
				src = unsafe.Add(src, so)
				*(*string)(dst) = *(*string)(src)
				return nil
			})
		case fs.Type.Kind() == reflect.Pointer && fd.Type.Kind() == reflect.Pointer && fs.Type.Elem() == fd.Type.Elem():
			fieldCopiers = append(fieldCopiers, func(dst, src unsafe.Pointer) error {
				dst = unsafe.Add(dst, do)
				src = unsafe.Add(src, so)
				*(*unsafe.Pointer)(dst) = *(*unsafe.Pointer)(src)
				return nil
			})
		case fs.Type.Kind() == reflect.Struct && fd.Type.Kind() == reflect.Struct:
			c, err := getCopier(fd.Type, fs.Type)
			if err != nil {
				return nil, err
			}
			fieldCopiers = append(fieldCopiers, func(dst, src unsafe.Pointer) error {
				dst = unsafe.Add(dst, do)
				src = unsafe.Add(src, so)
				return c(dst, src)
			})
		default:
			return nil, fmt.Errorf("can't create copier for %s -> %s", fs.Type, fd.Type)
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

// Copied ...
func Copied[D, S any](src *S) (*D, error) {
	var dst D
	if err := Copy(&dst, src); err != nil {
		return nil, err
	}
	return &dst, nil
}
