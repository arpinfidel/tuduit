package rose

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
)

type Parser struct {
	textMsgPrefix string
}

func NewParser(textMsgPrefix string) *Parser {
	return &Parser{
		textMsgPrefix: textMsgPrefix,
	}
}

type Rose struct {
	Valid  bool
	Errors []error
}

func (p *Parser) ParseArgs(args []string, flags map[string]string, target any) (Rose, error) {
	return parseArgs(args, flags, target)
}

func (p *Parser) ParseJSON(jsonBytes []byte, target any) (Rose, error) {
	return parseJSON(jsonBytes, target)
}

func (p *Parser) ParseTextMsg(text string, target any) (Rose, error) {
	return parseTextMsg(p.textMsgPrefix, text, target)
}

func Help(target any) (string, error) {
	return help(target)
}

type field struct {
	Names    map[string]struct{}
	Required bool
	Variadic bool
	Default  string
	Flatten  bool
}

func castJSON(v string, t reflect.Type) (i any, err error) {
	// create instace of type then unmarshal
	i = reflect.New(t).Interface()
	err = json.Unmarshal([]byte(v), i)
	if err != nil {
		return nil, fmt.Errorf("argument is not a valid type: %s", v)
	}

	// dereference pointer
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		i = reflect.ValueOf(i).Elem().Interface()
	}
	return i, nil
}

func castType(v string, t reflect.Type) (val any, err error) {
	switch t.Kind() {
	default:
		i, err := castJSON(v, t)
		if err != nil {
			return nil, err
		}
		return i, nil
	case reflect.String:
		return v, nil
	case reflect.Int:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i, nil
		}

		return nil, fmt.Errorf("argument is not a valid integer: %s", v)
	case reflect.Int32:
		i, err := strconv.ParseInt(v, 10, 32)
		if err == nil {
			return int32(i), nil
		}

		return nil, fmt.Errorf("argument is not a valid integer: %s", v)
	case reflect.Int64:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i, nil
		}
		d, err := str2duration.ParseDuration(v)
		if err == nil {
			return d, nil
		}
		return nil, fmt.Errorf("argument is not a valid integer or duration: %s", v)
	case reflect.Float32:
		i, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, fmt.Errorf("argument is not a valid float: %s", v)
		}
		return float32(i), nil
	case reflect.Float64:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("argument is not a valid float: %s", v)
		}
		return i, nil
	case reflect.Bool:
		if v == "" {
			return true, nil
		}

		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("argument is not a valid boolean: %s", v)
		}
		return b, nil
	// case time
	case reflect.Struct:
		switch t.String() {
		default:
			i, err := castJSON(v, t)
			if err != nil {
				return nil, err
			}
			return i, nil

		case "time.Time":
			i, err := time.Parse("2006-01-02 15:04:05", v)
			if err == nil {
				return i, nil
			}
			i, err = time.Parse("2006-01-02 15:04", v)
			if err == nil {
				return i, nil
			}
			i, err = time.Parse("2006-01-02", v)
			if err == nil {
				return i, nil
			}

			return nil, fmt.Errorf("argument is not a valid time: %s", v)
		}

	case reflect.Ptr:
		v, err := castType(v, t.Elem())
		if err != nil {
			return nil, err
		}

		// get pointer of underlying type of v
		p := reflect.New(t.Elem())
		p.Elem().Set(reflect.ValueOf(v))
		return p.Interface(), nil
	case reflect.Slice:
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			return castJSON(v, t)
		}
		switch t.Elem().Kind() {
		default:
			return nil, fmt.Errorf("argument is not a valid slice: %s", v)
		case reflect.String:
			return strings.Split(v, ","), nil
		case reflect.Int:
			i := []int{}
			s := strings.Split(v, ",")
			for _, v := range s {
				iv, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("argument is not a valid integer: %s", v)
				}
				i = append(i, iv)
			}
			return i, nil
		}
	}
}

func parseArgs(args []string, flags map[string]string, target any) (rose Rose, err error) {
	if reflect.TypeOf(target).Kind() != reflect.Pointer {
		return rose, fmt.Errorf("target must be a pointer")
	}

	rose = Rose{
		Valid: true,
	}

	fields := []reflect.StructField{}
	values := []reflect.Value{}
	for i := 0; i < reflect.TypeOf(target).Elem().NumField(); i++ {
		structField := reflect.TypeOf(target).Elem().Field(i)
		fields = append(fields, structField)
		values = append(values, reflect.ValueOf(target).Elem().Field(i))
	}

	flattened := 0

	for i := 0; i < len(fields); i++ {
		typ := fields[i]
		value := values[i]

		tag := typ.Tag.Get("rose")
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
			val := subParts[1]

			switch attr {
			case "required":
				field.Required = true
			case "variadic":
				field.Variadic = true // TODO: implement logic
			case "default":
				field.Default = val
			case "flatten":
				field.Flatten = true
			}
		}

		if field.Flatten {
			if typ.Type.Kind() != reflect.Struct {
				return rose, fmt.Errorf("argument %s is not a struct", typ.Name)
			}

			subFields := []reflect.StructField{}
			subValues := []reflect.Value{}
			for j := 0; j < typ.Type.NumField(); j++ {
				subFields = append(subFields, typ.Type.Field(j))
				subValues = append(subValues, value.Field(j))
			}

			values = append(values[:i+1], append(subValues, values[i+1:]...)...)
			fields = append(fields[:i+1], append(subFields, fields[i+1:]...)...)

			flattened++
			continue
		}

		set := false

		if i-flattened < len(args) {
			val, err := castType(args[i-flattened], typ.Type)
			if err != nil {
				rose.Valid = false
				rose.Errors = append(rose.Errors, err)
				continue
			}

			value.Set(reflect.ValueOf(val))
			set = true
		}

		for name := range field.Names {
			val, ok := flags[name]
			if !ok {
				continue
			}

			v, err := castType(val, typ.Type)
			if err != nil {
				rose.Valid = false
				rose.Errors = append(rose.Errors, err)
				return rose, err
			}

			value.Set(reflect.ValueOf(v))
			set = true
			break
		}

		if !set && field.Default != "" {
			v, err := castType(field.Default, typ.Type)
			if err != nil {

				return rose, err
			}
			value.Set(reflect.ValueOf(v))
			set = true
		}

		if !set && field.Required && values[i].IsZero() {
			rose.Valid = false
			rose.Errors = append(rose.Errors, fmt.Errorf("argument %s is required", typ.Name))

			continue
		}

	}

	return rose, nil
}

func parseJSON(jsonBytes []byte, target any) (rose Rose, err error) {
	j := map[string]json.RawMessage{}

	err = json.Unmarshal(jsonBytes, &j)
	if err != nil {
		return rose, err
	}

	flags := map[string]string{}
	for k, v := range j {
		v := string(v)
		if v[0] == '"' && v[len(v)-1] == '"' {
			v = v[1 : len(v)-1]
		}
		flags[k] = v
	}

	return parseArgs(nil, flags, target)
}

// parseTextMsg parses a text message into a struct
// first line is treated as args and inline flags
// each subsquent line that starts with flagPrefix is treated as a flag
// keep reading flag until the next line that starts with flagPrefix or end of message
func parseTextMsg(flagPrefix string, text string, target any) (Rose, error) {
	args := []string{}
	flags := map[string]string{}

	lines := strings.Split(text, "\n")

	firstLine := lines[0]
	flSplit := []string{}
	flag := ""
	part := ""
	isArgs := true
	if firstLine != "" {
		flSplit = strings.Fields(firstLine)
	}

	for _, str := range flSplit {
		if strings.HasPrefix(str, flagPrefix) {
			isArgs = false

			if flag != "" {
				flags[flag] = part
			}
			flag = strings.TrimPrefix(str, flagPrefix)
			continue
		}

		if isArgs {
			args = append(args, str)
			continue
		}

		if part != "" {
			part += " "
		}
		part += str
	}

	if flag != "" {
		flags[flag] = part
	}

	flag = ""
	part = ""
	for _, line := range lines[1:] {
		if strings.HasPrefix(line, flagPrefix) {
			// set previous flag
			flags[flag] = part
			parts := strings.SplitN(line, " ", 2)
			flag = strings.TrimPrefix(parts[0], flagPrefix)
			part = parts[1]
			continue
		}

		part += "\n" + line
	}
	flags[flag] = part

	rose, err := parseArgs(args, flags, target)
	if err != nil {

		return rose, err
	}

	return rose, nil
}

// help returns the help message for the given struct
func help(target any) (string, error) {
	res := ""

	if reflect.TypeOf(target).Kind() == reflect.Struct {
		target = reflect.New(reflect.TypeOf(target)).Interface()
	}

	fields := []reflect.StructField{}

	for i := 0; i < reflect.TypeOf(target).Elem().NumField(); i++ {
		structField := reflect.TypeOf(target).Elem().Field(i)
		fields = append(fields, structField)
	}

	for i := 0; i < len(fields); i++ {
		typ := fields[i]

		tag := typ.Tag.Get("rose")
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
			val := subParts[1]

			switch attr {
			case "required":
				field.Required = true
			case "variadic":
				field.Variadic = true // TODO: implement logic
			case "default":
				field.Default = val
			case "flatten":
				field.Flatten = true
			}
		}

		if field.Flatten {
			if typ.Type.Kind() != reflect.Struct {
				return res, fmt.Errorf("argument %s is not a struct", typ.Name)
			}

			subFields := []reflect.StructField{}
			for j := 0; j < typ.Type.NumField(); j++ {
				subFields = append(subFields, typ.Type.Field(j))
			}

			fields = append(fields[:i+1], append(subFields, fields[i+1:]...)...)

			continue
		}

		names := []string{}
		for name := range field.Names {
			names = append(names, name)
		}

		attrs := []string{}
		if field.Required {
			attrs = append(attrs, "required")
		} else {
			attrs = append(attrs, "optional")
		}
		if field.Flatten {
			attrs = append(attrs, "flatten")
		}
		if field.Default != "" {
			attrs = append(attrs, fmt.Sprintf("default=%s", field.Default))
		}
		if field.Variadic {
			attrs = append(attrs, "variadic")
		}

		res += fmt.Sprintf("%s: %s %s\n", strings.Join(names, " | "), typ.Type.Kind().String(), strings.Join(attrs, ", "))
	}

	return res, nil
}
