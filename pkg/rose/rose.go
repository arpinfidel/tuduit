package rose

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type field struct {
	Names    map[string]struct{}
	Required bool
	Variadic bool
}

func castType(v string, field reflect.StructField) (val any, err error) {
	switch field.Type.Kind() {
	default:
		// create instace of type then unmarshal
		i := reflect.New(field.Type).Interface()
		err := json.Unmarshal([]byte(v), i)
		if err != nil {
			return nil, fmt.Errorf("argument is not a valid type")
		}
		// dereference pointer
		if reflect.TypeOf(i).Kind() == reflect.Ptr {
			i = reflect.ValueOf(i).Elem().Interface()
		}
		return i, nil
	case reflect.String:
		return v, nil
	case reflect.Int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("argument is not a valid integer")
		}
		return i, nil
	case reflect.Bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("argument is not a valid boolean")
		}
		return b, nil
	}
}

func parseArgs[T any](args []string, flags map[string]string) (T, error) {
	var target T
	var targetI any = target

	// if target is struct, get its pointer
	if reflect.TypeOf(target).Kind() == reflect.Struct {
		targetI = &target
	}

	for i := 0; i < reflect.TypeOf(targetI).Elem().NumField(); i++ {
		structField := reflect.TypeOf(targetI).Elem().Field(i)
		tag := structField.Tag.Get("rose")
		if tag == "" {
			continue
		}

		partsStr := strings.Split(tag, ",")
		field := field{
			Names: map[string]struct{}{},
		}
		for _, part := range partsStr {
			subParts := strings.SplitN(part, "=", 2)
			if len(subParts) == 1 {
				field.Names[subParts[0]] = struct{}{}
				continue
			}

			attr := subParts[0]
			// val := subParts[1]

			switch attr {
			case "required":
				field.Required = true
			case "variadic":
				field.Variadic = true // TODO: implement logic
			}
		}

		set := false

		if i < len(args) {
			val, err := castType(args[i], structField)
			if err != nil {
				return target, err
			}
			reflect.ValueOf(&target).Elem().Field(i).Set(reflect.ValueOf(val))
			set = true
		}

		for name := range field.Names {
			val, ok := flags[name]
			if !ok {
				continue
			}

			v, err := castType(val, structField)
			if err != nil {
				return target, err
			}

			reflect.ValueOf(&target).Elem().Field(i).Set(reflect.ValueOf(v))
			set = true
			break
		}

		if !set && field.Required && reflect.ValueOf(&target).Elem().Field(i).IsZero() {
			return target, fmt.Errorf("argument %s is required", structField.Name)
		}
	}

	return target, nil
}

func parseJSON[T any](jsonBytes []byte) (T, error) {
	var target T

	j := map[string]json.RawMessage{}

	err := json.Unmarshal(jsonBytes, &j)
	if err != nil {
		return target, err
	}

	flags := map[string]string{}
	for k, v := range j {
		v := string(v)
		if v[0] == '"' && v[len(v)-1] == '"' {
			v = v[1 : len(v)-1]
		}
		flags[k] = v
	}

	return parseArgs[T](nil, flags)
}
