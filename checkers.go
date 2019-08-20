package form

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
)

var checkers map[string]CheckFunc

func init() {
	checkers = map[string]CheckFunc{
		"required":     Required,
		"min":          Min,
		"max":          Max,
		"range":        Range,
		"alpha":        Alpha,
		"numeric":      Numeric,
		"alphanumeric": AlphaNumeric,
		"alphadash":    AlphaDash,
		"username":     UserName,
		"float":        Float,
		"integer":      Integer,
		"email":        Email,
		"ipv4":         IPv4,
		"mobile":       Mobile,
		"mobile2":      Mobile2,
		"tel":          Tel,
		"phone":        Phone,
		"idcard":       IDCard,
	}
}

//AddCheckFunc add CheckFunc
func AddCheckFunc(name string, c CheckFunc) {
	checkers[name] = c
}

//Context context
type Context struct {
	//Input 表单值
	Input string
	//Title 显示的标题
	Title string
	//Field 字段Type
	Field reflect.StructField
	//Value 字段Value
	Value reflect.Value
	//Params 检测器参数
	Params []string
	//Ctx echo的context
	Ctx echo.Context
}

//CheckFunc 检测函数
type CheckFunc func(c Context) error

//Required required
func Required(c Context) error {
	if !IsRequired(c.Input) {
		return fmt.Errorf("%s不能为空", c.Title)
	}
	return nil
}

//MinOrMax min or max
func MinOrMax(ctx Context, t string) error {
	if ctx.Input == "" {
		return nil
	}
	if len(ctx.Params) != 1 {
		return fmt.Errorf("参数错误")
	}
	if IsIntType(ctx.Field) {
		n, err := strconv.ParseInt(ctx.Params[0], 10, 64)
		if err != nil {
			return fmt.Errorf("参数错误:%v", err)
		}
		v, err := strconv.ParseInt(ctx.Input, 10, 64)
		if err != nil {
			return fmt.Errorf("输入的值错误:%v", err)
		}
		if t == "min" {
			if v < n {
				return fmt.Errorf("%s不能小于%d", ctx.Title, n)
			}
		} else {
			if v > n {
				return fmt.Errorf("%s不能大于%d", ctx.Title, n)
			}
		}
	} else if IsFloatType(ctx.Field) {
		n, err := strconv.ParseFloat(ctx.Params[0], 64)
		if err != nil {
			return fmt.Errorf("参数错误:%v", err)
		}
		v, err := strconv.ParseFloat(ctx.Input, 64)
		if err != nil {
			return fmt.Errorf("输入的值错误:%v", err)
		}
		if t == "min" {
			if v < n {
				return fmt.Errorf("%s不能小于%f", ctx.Title, n)
			}
		} else {
			if v > n {
				return fmt.Errorf("%s不能大于%f", ctx.Title, n)
			}
		}
	} else if IsStringType(ctx.Field) {
		n, err := strconv.Atoi(ctx.Params[0])
		if err != nil {
			return fmt.Errorf("参数错误:%v", err)
		}
		if t == "min" {
			if len(ctx.Input) < n {
				return fmt.Errorf("%s的长度不能小于%d", ctx.Title, n)
			}
		} else {
			if len(ctx.Input) > n {
				return fmt.Errorf("%s的长度不能大于%d", ctx.Title, n)
			}
		}
	} else {
		return fmt.Errorf("未支持格式%v", ctx.Field.Type.Kind())
	}
	return nil
}

//Min min
func Min(ctx Context) error {
	return MinOrMax(ctx, "min")
}

//Max max
func Max(ctx Context) error {
	return MinOrMax(ctx, "max")
}

//Range range
func Range(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if len(ctx.Params) != 2 {
		return fmt.Errorf("参数错误")
	}
	params := ctx.Params
	ctx.Params = []string{params[0]}
	if err := Min(ctx); err != nil {
		return err
	}
	ctx.Params = []string{params[1]}
	return Max(ctx)
}

//Alpha alpha
func Alpha(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsAlpha(ctx.Input) {
		return fmt.Errorf("%s只允许包含字母", ctx.Title)
	}
	return nil
}

//Numeric numeric
func Numeric(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsNumeric(ctx.Input) {
		return fmt.Errorf("%s只允许包含数字", ctx.Title)
	}
	return nil
}

//AlphaNumeric alpha or numeric
func AlphaNumeric(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsAlphaNumeric(ctx.Input) {
		return fmt.Errorf("%s只允许包含数字或字母", ctx.Title)
	}
	return nil
}

//AlphaDash 只允许包含数字或者字母以及下划线
func AlphaDash(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsAlphaDash(ctx.Input) {
		return fmt.Errorf("%s只允许包含数字或字母以及下划线", ctx.Title)
	}
	return nil
}

//UserName 只允许包含数字或者字母以及下划线，并且第一个字符必须为字母
func UserName(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsAlphaDash(ctx.Input) {
		return fmt.Errorf("%s只允许包含数字或字母以及下划线", ctx.Title)
	}
	s := fmt.Sprintf("%c", ctx.Input[0])
	if !IsAlpha(s) {
		return fmt.Errorf("%s的第一个字符必须为字母", ctx.Title)
	}
	if ctx.Input[len(ctx.Input)-1] == '_' {
		return fmt.Errorf("%s的最后一个字符不能为_", ctx.Title)
	}
	return nil
}

//Float 必须为浮点数
func Float(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsFloat(ctx.Input) {
		return fmt.Errorf("%s必须为浮点数", ctx.Title)
	}
	return nil
}

//Integer 必须为整数
func Integer(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsInteger(ctx.Input) {
		return fmt.Errorf("%s必须为整数", ctx.Title)
	}
	return nil
}

//Email email
func Email(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsEmail(ctx.Input) {
		return fmt.Errorf("%s不是正确email格式", ctx.Title)
	}
	return nil
}

//IPv4 必须为IPv4格式
func IPv4(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsIPv4(ctx.Input) {
		return fmt.Errorf("%s必须为正确的IPv4格式", ctx.Title)
	}
	return nil
}

//Mobile 必须为手机号码
func Mobile(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsMobile(ctx.Input) {
		return fmt.Errorf("%s必须为正确的手机号码", ctx.Title)
	}
	return nil
}

//Mobile2 必须为手机号码
func Mobile2(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsMobile2(ctx.Input) {
		return fmt.Errorf("%s必须为正确的手机号码", ctx.Title)
	}
	return nil
}

//Tel 必须为座机
func Tel(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsTel(ctx.Input) {
		return fmt.Errorf("%s必须为正确的座机号码", ctx.Title)
	}
	return nil
}

//Phone 必须为座机或手机
func Phone(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsPhone(ctx.Input) {
		return fmt.Errorf("%s必须为正确的手机或座机号码", ctx.Title)
	}
	return nil
}

//IDCard 必须为身份证号码
func IDCard(ctx Context) error {
	if ctx.Input == "" {
		return nil
	}
	if !IsIDCard(ctx.Input) {
		return fmt.Errorf("%s必须为正确的身份证号码", ctx.Title)
	}
	return nil
}
