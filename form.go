package form

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var (
	//FormFields form key
	FormFields = []string{"form", "json"}
	//LabelFields title
	LabelFields = []string{"title", "label", "json"}
	//ValidField title
	ValidField = "valid"
)

//Check 检测表单值
func Check(ctx echo.Context, o interface{}) error {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("参数必须为struct")
	}
	return checkStruct(t, v, ctx)
}

//Bind 绑定表单值
func Bind(ctx echo.Context, o interface{}) error {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		t = t.Elem()
		v = v.Elem()
	} else {
		return fmt.Errorf("参数必须为struct的指针")
	}
	return bindStruct(t, v, ctx)
}

func bindStruct(t reflect.Type, v reflect.Value, ctx echo.Context) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Type.Kind() == reflect.Struct {
			if err := bindStruct(field.Type, value, ctx); err != nil {
				return err
			}
		} else {
			if err := bindField(field, value, ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func bindField(f reflect.StructField, v reflect.Value, ctx echo.Context) error {
	title := defaultField(f, LabelFields)
	input := ctx.FormValue(defaultField(f, FormFields))
	if input == "" {
		return nil
	}
	if IsIntType(f) {
		if IsUintType(f) {
			value, err := strconv.ParseUint(input, 10, 64)
			if err != nil {
				return fmt.Errorf("%s必须为整数", title)
			}
			switch f.Type.Kind() {
			case reflect.Uint:
				if strconv.IntSize == 32 {
					if value > math.MaxUint32 {
						return fmt.Errorf("%s的数值越界", title)
					}
				}
				v.Set(reflect.ValueOf(uint(value)))
			case reflect.Uint8:
				if value > math.MaxUint8 {
					return fmt.Errorf("%s的数值越界", title)
				}
				v.Set(reflect.ValueOf(uint8(value)))
			case reflect.Uint16:
				if value > math.MaxUint16 {
					return fmt.Errorf("%s的数值越界", title)
				}
				v.Set(reflect.ValueOf(uint16(value)))
			case reflect.Uint32:
				if value > math.MaxUint32 {
					return fmt.Errorf("%s的数值越界", title)
				}
				v.Set(reflect.ValueOf(uint32(value)))
			case reflect.Uint64:
				v.Set(reflect.ValueOf(value))
			}
		} else {
			value, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				return fmt.Errorf("%s必须为整数", title)
			}
			switch f.Type.Kind() {
			case reflect.Int:
				if strconv.IntSize == 32 {
					if value > math.MaxInt32 {
						return fmt.Errorf("%s的数值越界", title)
					}
					if value < math.MinInt32 {
						return fmt.Errorf("%s的数值越界", title)
					}
				}
				v.Set(reflect.ValueOf(int(value)))
			case reflect.Int8:
				if value > math.MaxInt8 {
					return fmt.Errorf("%s的数值越界", title)
				}
				if value < math.MinInt8 {
					return fmt.Errorf("%s的数值越界", title)
				}
				v.Set(reflect.ValueOf(int8(value)))
			case reflect.Int16:
				if value > math.MaxInt16 {
					return fmt.Errorf("%s的数值越界", title)
				}
				if value < math.MinInt16 {
					return fmt.Errorf("%s的数值越界", title)
				}
				v.Set(reflect.ValueOf(int16(value)))
			case reflect.Int32:
				if value > math.MaxInt32 {
					return fmt.Errorf("%s的数值越界", title)
				}
				if value < math.MinInt32 {
					return fmt.Errorf("%s的数值越界", title)
				}
				v.Set(reflect.ValueOf(int32(value)))
			case reflect.Int64:
				v.Set(reflect.ValueOf(value))
			}
		}
	} else if IsFloatType(f) {
		value, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("%s必须为浮点数", title)
		}
		switch f.Type.Kind() {
		case reflect.Float32:
			if value > math.MaxFloat32 {
				return fmt.Errorf("%s的数值越界", title)
			}
			v.Set(reflect.ValueOf(float32(value)))
		case reflect.Float64:
			v.Set(reflect.ValueOf(value))
		}
	} else if IsStringType(f) {
		v.SetString(input)
	} else if f.Type.Kind() == reflect.Bool {
		v.SetBool(true) //凡是有值的皆为真
	}
	return nil
}

func checkStruct(t reflect.Type, v reflect.Value, ctx echo.Context) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Type.Kind() == reflect.Struct {
			if err := checkStruct(field.Type, value, ctx); err != nil {
				return err
			}
		} else {
			if err := checkField(field, value, ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkField(t reflect.StructField, v reflect.Value, ctx echo.Context) error {
	title := defaultField(t, LabelFields)
	input := ctx.FormValue(defaultField(t, FormFields))
	rules := parseRules(t)
	for _, r := range rules {
		if r.Name == "" {
			continue
		}
		checkFunc, ok := checkers[r.Name]
		if !ok {
			return fmt.Errorf("检测器%s找不到", r.Name)
		}
		c := Context{
			Input:  input,
			Title:  title,
			Params: r.Params,
			Field:  t,
			Value:  v,
			Ctx:    ctx,
		}
		if err := checkFunc(c); err != nil {
			return err
		}
	}
	return nil
}

func defaultField(t reflect.StructField, fields []string) string {
	var field string
	for _, f := range fields {
		field = t.Tag.Get(f)
		if field != "" {
			break
		}
	}
	if field == "" {
		field = t.Name
	}
	return field
}

type rule struct {
	Name   string
	Params []string
}

func parseRules(t reflect.StructField) []rule {
	rules := make([]rule, 0, 4)
	valid := t.Tag.Get(ValidField)
	for _, r := range strings.Split(valid, ";") {
		if strings.Contains(r, ":") {
			rl := strings.SplitN(r, ":", 2)
			if len(rl[0]) == 0 {
				continue
			}
			rules = append(rules, rule{
				Name:   rl[0],
				Params: strings.Split(rl[1], ","),
			})
		} else {
			rules = append(rules, rule{
				Name: r,
			})
		}
	}
	return rules
}
