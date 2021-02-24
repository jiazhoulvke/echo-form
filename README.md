## echo-form ##

echo框架的表单校验绑定库。


### 支持的校验规则 ###

- required 
  必填
- min 
  最小值。当字段类型为int或者float时，值不能小于最小值；当字段类型为string时，长度不能小于最小值
- max 
  最大值。当字段类型为int或者float时，值不能大于最大值；当字段类型为string时，长度不能大于最大值
- range 
  范围值。当字段类型为int或者float时，值不能小于最小值或大于最大值；当字段类型为string时，长度不能小于最小值或大于最大值
- alpha 
  只能有字母。
- numeric
  只能有数字
- alphanumeric
  只能有字母或数字
- alphadash
  只能有字母或数字及下划线
- username
  只能有字母或数字及下划线,且第一个字符必须为字母，最后一个字符不能为下划线
- float
  必须为能转为浮点数的字符串
- integer
  必须为能转为整数的字符串
- email
  必须为正确的email格式
- ipv4
  必须为正确的ipv4型字符串
- mobile
  必须为正确的手机号码
- mobile2
  随着手机号段增加得越来越频繁，已经很难及时更新正则了，mobile2只校验手机号码第一位是不是1并且长度是否为11位
- tel
  必须为正确的座机号码
- phone
  必须为正确的手机或座机号码
- idcard
  必须为正确的身份证号码



### 绑定 ###

目前支持以下类型：

- int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64
- float32,float64
- string
- bool
- time.Time


### 示例 ###

```go
package main

import (
	"fmt"
	"log"

	form "github.com/jiazhoulvke/echo-form"
	echo "github.com/labstack/echo/v4"
)

type userInfo struct {
	UserName  string  `validate:"required;range:8,16;username"`
	Password  string  `validate:"required;range:32,512"`
	EnName    string  `validate:"alpha"`
	Number    string  `validate:"numeric"`
	Age       int     `validate:"min:18"`
	Weight    float64 `validate:"max:200"`
	StrAge    string  `validate:"integer"`
	StrWeight string  `validate:"float"`
	LastIP    string  `validate:"ipv4"`
	Email     string  `validate:"email"`
	Mobile    string  `validate:"mobile"`
	Tel       string  `validate:"tel"`
	Phone     string  `validate:"phone"`
	IDCard    string  `validate:"idcard"`
}

func main() {
	var user = userInfo{}
	e := echo.New()
	e.Binder = form.New(func(f *form.Form) {
		//可以修改预定义的tag，比如把校验用的tag从valid改成了validate
		f.ValidField = "validate"
	})
	e.POST("/", func(c echo.Context) error {
		if err := c.Bind(&user); err != nil {
			return c.String(200, fmt.Sprintf("输入有误:%v", err))
		}
		//也可以不注册binder，直接用
		// f := form.New()
		// if err := f.Bind(&d, c); err != nil {
		// 	return c.String(200, "")
		// }
		// if err := f.Check(&d, c); err != nil {
		// 	return c.String(200, "")
		// }
		return c.HTML(200, fmt.Sprintf("你的用户名:%s", user.UserName))
	})
	log.Println("server started!")
	log.Fatal(e.Start(":8080"))
}
```

如果需要绑定+校验一步到位的话可以自定义一个struct：

```go
type customBinder struct {
	f *form.Form
}

func (cb customBinder) Bind(o interface{}, c echo.Context) error {
	if err := cb.f.Bind(o, c); err != nil {
		return err
	}
	if err := cb.f.Check(o, c); err != nil {
		return err
	}
	return nil
}
```

然后将其注册为binder:

```go
cb := customBinder{
	f: form.New(),
}
e.Binder = cb
```