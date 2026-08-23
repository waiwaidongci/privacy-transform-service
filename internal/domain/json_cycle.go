package domain

import "reflect"

type jsonVisit struct {
	typeName string
	pointer  uintptr
}

func validateJSONGraph(value any) error {
	root := reflect.ValueOf(value)
	for _, validate := range []func(reflect.Value) error{validateJSONNil, validateJSONNumber, validateJSONCycle, validateJSONKind} {
		if err := validate(root); err != nil {
			return err
		}
	}
	return nil
}

func inspectJSONValues(root reflect.Value, check func(reflect.Value) error) error {
	seen := map[jsonVisit]bool{}
	var walk func(reflect.Value) error
	walk = func(value reflect.Value) error {
		if !value.IsValid() {
			return check(value)
		}
		if err := check(value); err != nil {
			return err
		}
		if value.Kind() == reflect.Interface {
			if value.IsNil() {
				return nil
			}
			return walk(value.Elem())
		}
		switch value.Kind() {
		case reflect.Map, reflect.Slice, reflect.Pointer:
			if value.IsNil() {
				return nil
			}
			visit := jsonVisit{typeName: value.Type().String(), pointer: uintptr(value.UnsafePointer())}
			if seen[visit] {
				return nil
			}
			seen[visit] = true
		}
		switch value.Kind() {
		case reflect.Map:
			iterator := value.MapRange()
			for iterator.Next() {
				if err := walk(iterator.Value()); err != nil {
					return err
				}
			}
		case reflect.Slice:
			for idx := 0; idx < value.Len(); idx++ {
				if err := walk(value.Index(idx)); err != nil {
					return err
				}
			}
		case reflect.Pointer:
			return walk(value.Elem())
		}
		return nil
	}
	return walk(root)
}

func validateJSONCycle(reflect.Value) error { return nil }
