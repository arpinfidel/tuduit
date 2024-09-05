package rose

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	str2duration "github.com/xhit/go-str2duration/v2"
)

type Parser struct {
	ctx           *ctxx.Context
	textMsgPrefix string
}

func NewParser(ctx *ctxx.Context, textMsgPrefix string) *Parser {
	return &Parser{
		ctx:           ctx,
		textMsgPrefix: textMsgPrefix,
	}
}

type Rose struct {
	Valid  bool
	Errors []error
}

func (p *Parser) ParseArgs(args []string, flags map[string]string, target any) (Rose, error) {
	return p.parseArgs(args, flags, target)
}

func (p *Parser) ParseJSON(jsonBytes []byte, target any) (Rose, error) {
	return p.parseJSON(jsonBytes, target)
}

func (p *Parser) ParseTextMsg(text string, target any) (Rose, error) {
	return p.parseTextMsg(p.textMsgPrefix, text, target)
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

func (p *Parser) setTimezone(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), p.ctx.User.CreatedAt.Location())
}

func (p *Parser) castType(v string, t reflect.Type) (val any, err error) {
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
		pt := t
		if t.Kind() != reflect.Ptr {
			pt = reflect.PtrTo(t)
		}
		if pt.Implements(reflect.TypeOf(new(entity.Base36Parser)).Elem()) {
			i := reflect.New(t).Interface()
			err := i.(entity.Base36Parser).ParseBase36(v)
			if err == nil {
				return reflect.ValueOf(i).Elem().Interface(), nil
			}
		}

		switch t.String() {
		default:
			i, err := castJSON(v, t)
			if err != nil {
				return nil, err
			}
			return i, nil

		case "time.Time":
			errs := []string{}
			i, err := time.Parse("2006-01-02 15:04:05", v)
			if err == nil {
				return p.setTimezone(i), nil
			}
			errs = append(errs, err.Error())

			i, err = time.Parse("2006-01-02 15:04", v)
			if err == nil {
				return p.setTimezone(i), nil
			}
			errs = append(errs, err.Error())

			i, err = time.Parse("2006-01-02", v)
			if err == nil {
				return p.setTimezone(i), nil
			}
			errs = append(errs, err.Error())

			return nil, fmt.Errorf("argument is not a valid time: %s\n%s", v, strings.Join(errs, "\n\t"))
		}

	case reflect.Ptr:
		v, err := p.castType(v, t.Elem())
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

		split := strings.Split(v, ",")
		switch t.Elem().Kind() {
		default:
			s := reflect.MakeSlice(t, 0, 0)
			for _, v := range split {
				i, err := p.castType(v, t.Elem())
				if err != nil {
					return nil, err
				}
				s = reflect.Append(s, reflect.ValueOf(i))
			}
			return s.Interface(), nil
		case reflect.String:
			return split, nil
		case reflect.Int:
			i := []int{}
			s := split
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

func (p *Parser) parseArgs(args []string, flags map[string]string, target any) (rose Rose, err error) {
	if reflect.TypeOf(target).Kind() != reflect.Pointer {
		return Rose{}, fmt.Errorf("target must be a pointer")
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
				return Rose{}, fmt.Errorf("argument %s is not a struct", typ.Name)
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
			val, err := p.castType(args[i-flattened], typ.Type)
			if err != nil {
				rose.Valid = false
				rose.Errors = append(rose.Errors, err)
				continue
			}

			value.Set(reflect.ValueOf(val))
			set = true
		}

		names := []string{}
		for name := range field.Names {
			names = append(names, name)

			val, ok := flags[name]
			if !ok {
				continue
			}

			v, err := p.castType(val, typ.Type)
			if err != nil {
				rose.Valid = false
				rose.Errors = append(rose.Errors, err)
				return Rose{}, err
			}

			value.Set(reflect.ValueOf(v))
			set = true
			break
		}

		if !set && field.Default != "" {
			v, err := p.castType(field.Default, typ.Type)
			if err != nil {
				return Rose{}, err
			}
			value.Set(reflect.ValueOf(v))
			set = true
		}

		if !set && field.Required && values[i].IsZero() {
			rose.Valid = false
			rose.Errors = append(rose.Errors, fmt.Errorf("argument %s is required", names[0]))

			continue
		}

	}

	return rose, nil
}

func (p *Parser) parseJSON(jsonBytes []byte, target any) (rose Rose, err error) {
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

	return p.parseArgs(nil, flags, target)
}

// ParseHTTP parses an HTTP request into a struct
func (p *Parser) ParseHTTP(r *http.Request, target any) (rose Rose, err error) {
	// Parse query parameters
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	// Parse body
	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return rose, fmt.Errorf("failed to read request body: %w", err)
		}
		defer r.Body.Close()

		if len(body) > 0 {
			_, err := p.parseJSON(body, target)
			if err != nil {
				return rose, fmt.Errorf("failed to parse request body: %w", err)
			}
		}
	}

	// Merge query parameters and body parameters
	flags := make(map[string]string)
	for k, v := range queryParams {
		flags[k] = v
	}

	// Parse path parameters
	pathParams := make([]string, 0)
	if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
		for i, param := range routeContext.URLParams.Values {
			flags[routeContext.URLParams.Keys[i]] = param
		}
	}

	return p.parseArgs(pathParams, flags, target)
}

// parseTextMsg parses a text message into a struct
// first line is treated as args and inline flags
// each subsquent line that starts with flagPrefix is treated as a flag
// keep reading flag until the next line that starts with flagPrefix or end of message
func (p *Parser) parseTextMsg(flagPrefix string, text string, target any) (Rose, error) {
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
				part = ""
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

	rose, err := p.parseArgs(args, flags, target)
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

var timeType = reflect.TypeOf(time.Time{})

func ChangeTimezone(v any, loc *time.Location) (res any) {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("not ptr")
	}

	vv := changeTimezone(reflect.ValueOf(v).Elem(), loc)
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(vv))

	return vv
}

func changeTimezone(val reflect.Value, loc *time.Location) (res any) {

	if val.Interface() == nil {
		return nil
	}

	if val.Type() == timeType {
		t := val.Interface().(time.Time).In(loc)
		return t
	}

	switch val.Kind() {
	default:
		return val.Interface()

	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			t := changeTimezone(field, loc)
			if field.Kind() == reflect.Ptr && field.IsNil() {
				continue
			}
			field.Set(reflect.ValueOf(t))
		}
		return val.Interface()

	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}

		t := changeTimezone(val.Elem(), loc)
		tPtr := reflect.New(val.Type().Elem())
		tPtr.Elem().Set(reflect.ValueOf(t))
		val.Set(reflect.ValueOf(tPtr.Interface()))
		return val.Interface()

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			t := changeTimezone(val.Index(i), loc)
			val.Index(i).Set(reflect.ValueOf(t))
		}
		return val.Interface()
	}
}

func JSONMarshal(v any) ([]byte, error) {
	jsonc := jsoniter.Config{TagKey: "rose"}.Froze()
	return jsonc.Marshal(v)
}
