## echo-form ##

echo框架的表单校验绑定库。


#### 支持的校验规则 ####

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
- name
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


### 示例 ###

```go
package main

import (
	"fmt"
	"log"

	"github.com/jiazhoulvke/echo-form"
	"github.com/labstack/echo"
)

type userInfo struct {
	UserName  string  `valid:"required;range:8,16;name"`
	Password  string  `valid:"required;range:32,512"`
	EnName    string  `valid:"alpha"`
	Number    string  `valid:"numeric"`
	Age       int     `valid:"min:18"`
	Weight    float64 `valid:"max:200"`
	StrAge    string  `valid:"integer"`
	StrWeight string  `valid:"float"`
	LastIP    string  `valid:"ipv4"`
	Email     string  `valid:"email"`
	Mobile    string  `valid:"mobile"`
	Tel       string  `valid:"tel"`
	Phone     string  `valid:"phone"`
	IDCard    string  `valid:"idcard"`
}

func main() {
	var user = userInfo{
		UserName:  "jiazhoulvke",
		Password:  "abcdefg1234567890abcdefg1234567890",
		EnName:    "foobar",
		Age:       24,
		Weight:    100.123,
		StrAge:    "24",
		StrWeight: "100.123",
		LastIP:    "127.0.0.1",
		Email:     "jiazhoulvke+1984@gmail.com",
		Mobile:    "13812345678",
		Tel:       "07012345678",
		Phone:     "13812345678",
		IDCard:    "43214321432143211X",
	}
	e := echo.New()
	e.POST("/", func(ctx echo.Context) error {
		if err := form.Check(ctx, &user); err != nil {
			return ctx.HTML(200, fmt.Sprintf("输入有误:%v", err))
		}
		if err := form.Bind(ctx, &user); err != nil {
			return ctx.HTML(200, fmt.Sprintf("输入有误:%v", err))
		}
		return ctx.HTML(200, fmt.Sprintf("你的用户名:%s", user.UserName))
	})

	log.Println("server started!")
	log.Fatal(e.Start(":8080"))
}
```
