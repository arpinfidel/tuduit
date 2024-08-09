package handler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/urfave/cli/v2"
)

func (h *Handler) getArgs(ctx *ctxx.Context, c *cli.Context, v any) (res any, err error) {
	return h.validateArgsRaw(ctx, c.Args().Slice(), v)
}

func (h *Handler) validateArgsRaw(ctx *ctxx.Context, runArgs []string, v any) (res any, err error) {
	type arg struct {
		Name     string
		Required bool
		Variadic bool
	}

	args := []arg{}
	requiredCount := 0
	hasOptional := false
	hasVariadic := false
	for i := 0; i < reflect.TypeOf(v).Elem().NumField(); i++ {
		field := reflect.TypeOf(v).Elem().Field(i)
		tag := field.Tag.Get("arg")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")

		name := parts[0]

		required := false
		if len(parts) > 1 && parts[1] == "required" {
			required = true
		}
		if required {
			if hasOptional {
				return nil, fmt.Errorf("optional argument cannot be followed by required arguments")
			}
			requiredCount++
		} else {
			hasOptional = true
		}

		if field.Type.Kind() == reflect.Slice {
			if i != reflect.TypeOf(v).Elem().NumField()-1 {
				return nil, fmt.Errorf("variadic argument must be the last argument")
			}
			hasVariadic = true
		}

		args = append(args, arg{
			Name:     name,
			Required: required,
			Variadic: hasVariadic,
		})
	}

	if len(runArgs) < requiredCount {
		return nil, fmt.Errorf("missing required arguments")
	}

	if len(runArgs) > len(args) && !hasVariadic {
		return nil, fmt.Errorf("too many arguments")
	}

	castPrimitive := func(arg arg, v string, field reflect.StructField) (val any, err error) {
		switch field.Type.Kind() {
		case reflect.String:
			return v, nil
		case reflect.Int:
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("argument %s is not a valid integer", arg.Name)
			}
			return i, nil
		case reflect.Bool:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, fmt.Errorf("argument %s is not a valid boolean", arg.Name)
			}
			return b, nil
		}
		return nil, fmt.Errorf("argument %s is not a valid type", arg.Name)
	}

	// assign the arguments to the struct
	for i, arg := range args {
		if arg.Variadic {
			break
		}

		ra := runArgs[i]
		field := reflect.TypeOf(v).Elem().Field(i)

		if ra == "" {
			if arg.Required {
				return nil, fmt.Errorf("argument %s is required", arg.Name)
			}
			continue
		}

		val, err := castPrimitive(arg, ra, field)
		if err != nil {
			return nil, err
		}

		reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(val))
	}

	if hasVariadic && len(runArgs) >= len(args) {
		// arg := args[len(args)-1]
		field := reflect.TypeOf(v).Elem().Field(len(args) - 1)
		ra := []string{}
		for i := len(args) - 1; i < len(runArgs); i++ {
			ra = append(ra, runArgs[i])
		}

		// make slice of appropriate type
		slice := reflect.MakeSlice(field.Type, len(ra), len(ra))
		for i, v := range ra {
			slice.Index(i).Set(reflect.ValueOf(v))
		}
		reflect.ValueOf(v).Elem().Field(len(args) - 1).Set(slice)
	}

	return v, nil
}

func (h *Handler) makeFlags(v any) (res []cli.Flag, err error) {
	flags := []cli.Flag{}
	for i := 0; i < reflect.TypeOf(v).Elem().NumField(); i++ {
		field := reflect.TypeOf(v).Elem().Field(i)
		tag := field.Tag.Get("flag")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")

		if parts[0] == "" {
			return nil, fmt.Errorf("invalid flag tag: %s", tag)
		}
		name := parts[0]

		required := false
		aliases := []string{}
		for _, part := range parts[1:] {
			switch {
			case part == "required":
				required = true
			case strings.HasPrefix(part, "alias="):
				aliases = strings.Split(part[len("alias="):], ",")
			}
		}

		switch field.Type.Kind() {
		case reflect.String:
			flag := cli.StringFlag{
				Name:     name,
				Aliases:  aliases,
				Required: required,
			}
			flags = append(flags, &flag)
		case reflect.Bool:
			flag := cli.BoolFlag{
				Name:     name,
				Aliases:  aliases,
				Required: required,
			}
			flags = append(flags, &flag)
		case reflect.Int:
			flag := cli.IntFlag{
				Name:     name,
				Aliases:  aliases,
				Required: required,
			}
			flags = append(flags, &flag)
		case reflect.Int64:
			flag := cli.Int64Flag{
				Name:     name,
				Aliases:  aliases,
				Required: required,
			}
			flags = append(flags, &flag)
		case reflect.Slice:
			flag := cli.StringSliceFlag{
				Name:     name,
				Aliases:  aliases,
				Required: required,
			}
			flags = append(flags, &flag)
		default:
			return nil, fmt.Errorf("unsupported flag type: %s", field.Type.Kind())
		}
	}

	return flags, nil
}

func (h *Handler) getFlags(ctx *ctxx.Context, c *cli.Context, v any) (err error) {
	for i := 0; i < reflect.TypeOf(v).Elem().NumField(); i++ {
		field := reflect.TypeOf(v).Elem().Field(i)
		tag := field.Tag.Get("flag")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]

		fieldValue := reflect.ValueOf(v).Elem().Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(c.String(name))
		case reflect.Bool:
			fieldValue.SetBool(c.Bool(name))
		case reflect.Int:
			fieldValue.SetInt(int64(c.Int(name)))
		case reflect.Int64:
			fieldValue.SetInt(c.Int64(name))
		case reflect.Slice:
			fieldValue.Set(reflect.ValueOf(c.StringSlice(name)))
		default:
			return fmt.Errorf("unsupported flag type: %s", field.Type.Kind())
		}
	}

	return nil
}
