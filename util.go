package form

import (
	"reflect"
	"regexp"
	"strconv"
)

//IsRequired 是否有值
func IsRequired(str string) bool {
	return len(str) > 0
}

//IsNumeric 是否只包含数字
func IsNumeric(str string) bool {
	for _, v := range str {
		if '9' < v || v < '0' {
			return false
		}
	}
	return true
}

//IsAlphaNumeric 是否只包含数字或字母
func IsAlphaNumeric(str string) bool {
	for _, v := range str {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
			return false
		}
	}
	return true
}

var alphaDashPattern = regexp.MustCompile(`^\w+$`)

//IsAlphaDash 是否只包含数字或字母以及下划线
func IsAlphaDash(str string) bool {
	return alphaDashPattern.MatchString(str)
}

//IsAlpha 是否只包含字母
func IsAlpha(str string) bool {
	for _, v := range str {
		if v < 'A' || (v > 'Z' && v < 'a') || v > 'z' {
			return false
		}
	}
	return true
}

//IsIntType is int
func IsIntType(f reflect.StructField) bool {
	switch f.Type.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

//IsUintType is uint
func IsUintType(f reflect.StructField) bool {
	switch f.Type.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

//IsFloatType is float
func IsFloatType(f reflect.StructField) bool {
	return f.Type.Kind() == reflect.Float32 || f.Type.Kind() == reflect.Float64
}

//IsStringType is string
func IsStringType(f reflect.StructField) bool {
	return f.Type.Kind() == reflect.String
}

//IsFloat 是否为浮点数
func IsFloat(str string) bool {
	if _, err := strconv.ParseFloat(str, 64); err != nil {
		return false
	}
	return true
}

//IsInteger 是否为整数
func IsInteger(str string) bool {
	if _, err := strconv.ParseInt(str, 10, 64); err != nil {
		return false
	}
	return true
}

var emailPattern = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)

//IsEmail 是否为email
func IsEmail(str string) bool {
	return emailPattern.MatchString(str)
}

var ipv4Pattern = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)

//IsIPv4 是否为ipv4格式
func IsIPv4(str string) bool {
	return ipv4Pattern.MatchString(str)
}

var mobilePattern = regexp.MustCompile(`^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][056789]|[4][579]|99))\d{8}$`)

//IsMobile 是否为手机号码
func IsMobile(str string) bool {
	return mobilePattern.MatchString(str)
}

var telPattern = regexp.MustCompile(`^(0\d{2,3}(\-)?)?\d{7,8}$`)

//IsTel 是否为座机号码
func IsTel(str string) bool {
	return telPattern.MatchString(str)
}

//IsPhone 是否为手机或座机号码
func IsPhone(str string) bool {
	return IsMobile(str) || IsTel(str)
}

var idcartPattern = regexp.MustCompile(`(^\d{15}$)|(^\d{17}([0-9X])$)`)

//IsIDCard 是否为身份证号码
func IsIDCard(str string) bool {
	return idcartPattern.MatchString(str)
}
