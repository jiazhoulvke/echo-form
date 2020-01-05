package form

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
)

func makeContext(data url.Values) echo.Context {
	req := http.Request{}
	req.PostForm = data
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(&req, rec)
	return ctx
}

func goodData() url.Values {
	data := url.Values{}
	data.Set("username", "jiazhoulvke")
	data.Set("age", "24")
	data.Set("weight", "60")
	data.Set("Number", "24")
	return data
}

type foo struct {
	ID            uint    `valid:"max:10"`
	UserName      string  `form:"username" title:"用户名" valid:"required;range:8,16;username"`
	Age           int     `form:"age" label:"年龄" valid:"integer;range:18,60"`
	Weight        float32 `form:"weight" title:"体重" valid:"float;min:50;max:80"`
	EnglishName   string  `valid:"alpha"`
	Number        string  `valid:"required;numeric"`
	AliasName     string  `valid:"alphanumeric"`
	AliasName2    string  `valid:"alphadash"`
	Email         string  `valid:"email"`
	LastIP        string  `valid:"ipv4"`
	Mobile        string  `valid:"mobile"`
	Tel           string  `valid:"tel"`
	Phone         string  `valid:"phone"`
	IDCard        string  `valid:"idcard"`
	Height        float64 `valid:"max:2.5"`
	Count         uint64
	IsMan         bool
	Birthday      time.Time `valid:"required"`
	DefaultString string    `default:"foobar"`
	DefaultInt    int       `default:"42"`
}

type bar struct {
	Bar int `valid:"required;range:5,10"`
}

type baz struct {
	bar
	Baz string `valid:"min:10"`
}

func TestCheck(t *testing.T) {
	var tables = []struct {
		field string
		v     string
		ok    bool
	}{
		{"username", "", false},                  //required 不能为空
		{"username", "abcd", false},              //range 太短
		{"username", "jiazhoulvke123432", false}, //range 太长
		{"username", "jiazhoulvke.", false},      //username 不能有.
		{"username", "_jiazhoulvke", false},      //username 第一个必须是字母
		{"username", "jiazhoulvke_", false},      //username 最后一个不能是_
		{"username", "jiazhoulvke_1", true},
		{"age", "17", false},  //range 太小
		{"age", "61", false},  //range 太大
		{"age", "abc", false}, //integer 不是数字组成
		{"age", "20", true},
		{"weight", "49.9999", false}, //range 太小
		{"weight", "81.0001", false}, //range 太大
		{"weight", "abcd", false},    //float 不是浮点数
		{"weight", "50", true},
		{"weight", "66.666666", true},
		{"ID", "abc", false}, //max 不是数字
		{"ID", "5", true},
		{"EnglishName", "1234", false}, //alpha 不能有数字
		{"EnglishName", "edison", true},
		{"Number", "abc123", false}, //numeric 不能有字母
		{"Number", "23", true},
		{"AliasName", "efe123_", false}, //alphanumeric 不能有下划线
		{"AliasName", "abc123", true},
		{"AliasName2", "fewef123_.", false}, //alphadash 不能有.
		{"AliasName2", "abc123_", true},
		{"Email", "foo", false},      //email 错误的邮件格式
		{"Email", "foo@", false},     //email 错误的邮件格式
		{"Email", "foo@bar", false},  //email 错误的邮件格式
		{"Email", "foo@bar.", false}, //email 错误的邮件格式
		{"Email", "foo@bar.com", true},
		{"LastIP", "1.1.", false},            //ipv4 错误的ip地址
		{"LastIP", "256.256.256.256", false}, //ipv4 错误的ip地址
		{"LastIP", "127.0.0.1", true},
		{"Mobile", "1234567890", false},  //mobile 错误的手机号码
		{"Mobile", "91234567890", false}, //mobile 错误的手机号码
		{"Mobile", "13812345678", true},
		{"Tel", "13423456789", false},   //tel 错误的座机号码
		{"Tel", "0777712345678", false}, //tel 错误的座机号码
		{"Tel", "07651234567", true},
		{"Phone", "234323", false}, //phone 错误的手机或座机号码
		{"Phone", "13812345678", true},
		{"Phone", "07651234567", true},
		{"IDCard", "f12345678901234567", false}, //idcard 错误的身份证号码
		{"IDCard", "12345678901234567x", false}, //idcard 错误的身份证号码
		{"IDCard", "45678901234567x", false},    //idcard 错误的身份证号码
		{"IDCard", "4456789012345671", false},   //idcard 错误的身份证号码
		{"IDCard", "12345678901234567X", true},
		{"IDCard", "123456789012345", true},
		{"Height", "a", false}, //max 不是浮点数
		{"Height", "1.8", true},
	}

	var f foo
	var ctx echo.Context

	Convey("测试Check", t, func() {
		var err error
		var n int

		ctx = makeContext(goodData())

		err = Check(ctx, &n)
		So(err, ShouldNotBeNil) //n不是struct，应该报错

		err = Check(ctx, &f)
		So(err, ShouldBeNil)
	})

	Convey("测试错误struct tag", t, func() {
		var err error
		var testRange = struct {
			Foo int `valid:"range:1,5,7"`
		}{}
		ctx = makeContext(url.Values{
			"Foo": []string{"4"},
		})
		err = Check(ctx, &testRange)
		So(err, ShouldNotBeNil)

		var testMax = struct {
			Foo int `valid:"max:a"`
		}{}
		ctx = makeContext(url.Values{
			"Foo": []string{"4"},
		})
		err = Check(ctx, &testMax)
		So(err, ShouldNotBeNil)

		var testMax2 = struct {
			Foo int `valid:"max:1,3"`
		}{}
		ctx = makeContext(url.Values{
			"Foo": []string{"2"},
		})
		err = Check(ctx, &testMax2)
		So(err, ShouldNotBeNil)

		//测试不支持的类型，必须报错
		var testMax3 = struct {
			Foo bool `valid:"max:10"`
		}{}
		ctx = makeContext(url.Values{
			"Foo": []string{"2"},
		})
		err = Check(ctx, &testMax3)
		So(err, ShouldNotBeNil)

		var testMax4 = struct {
			Foo float64 `valid:"max:a"`
		}{}
		ctx = makeContext(url.Values{
			"Foo": []string{"2.2"},
		})
		err = Check(ctx, &testMax4)
		So(err, ShouldNotBeNil)

		var testMax5 = struct {
			Foo string `valid:"max:a"`
		}{}
		ctx = makeContext(url.Values{
			"Foo": []string{"abcd"},
		})
		err = Check(ctx, &testMax5)
		So(err, ShouldNotBeNil)

	})

	Convey("测试空值", t, func() {
		var err error
		var foo = struct {
			Foo int `valid:"max:20;min:10;range:1,10;numeric;alpha;alphanumeric;alphadash;username;float;integer;:;"`
		}{}
		ctx = makeContext(url.Values{})
		err = Check(ctx, &foo)
		So(err, ShouldBeNil)
	})

	Convey("测试内置检测器", t, func() {
		for i, d := range tables {
			t.Logf("test:%d:%v\n", i, d)
			data := goodData()
			data.Set(d.field, d.v)
			ctx = makeContext(data)
			err := Check(ctx, &f)
			if d.ok {
				So(err, ShouldBeNil)
			} else {
				t.Log(err)
				So(err, ShouldNotBeNil)
			}
		}
	})

	Convey("测试手机号码", t, func() {
		mobiles := []string{"13412345678", "19912345678", "17512345678", "13812345678"}
		for _, mobile := range mobiles {
			So(IsMobile(mobile), ShouldBeTrue)
		}
	})

	Convey("测试多重struct", t, func() {
		var b = baz{}
		var err error
		data := url.Values{
			b.Baz: []string{"1234567890"},
		}
		ctx = makeContext(data)
		err = Check(ctx, &b)
		So(err, ShouldNotBeNil) //b.bar.Bar为必填项，应该报错
		data.Set("Bar", "10")
		ctx = makeContext(data)
		err = Check(ctx, &b)
		So(err, ShouldBeNil) //已经填写了b.bar.Bar，不应该报错
	})

	Convey("测试附加检测器", t, func() {
		var f = struct {
			Foo string `valid:"hello"`
		}{}
		var err error
		ctx = makeContext(url.Values{
			"Foo": []string{"world"},
		})
		err = Check(ctx, &f)
		So(err, ShouldNotBeNil) //检测器hello不存在，应该报错

		AddCheckFunc("hello", func(ctx Context) error {
			if ctx.Input != "world" {
				return fmt.Errorf("%s的值必须是world", ctx.Title)
			}
			return nil
		})
		err = Check(ctx, &f)
		So(err, ShouldBeNil) //已添加检测器hello,应该不报错

		ctx = makeContext(url.Values{})
		err = Check(ctx, &f)
		So(err, ShouldNotBeNil) //值不正确，应该报错
	})
}

func TestBind(t *testing.T) {
	var f foo
	Convey("测试表单绑定", t, func() {
		var err error
		ctx := makeContext(url.Values{})

		var n int
		err = Bind(ctx, n)
		So(err, ShouldNotBeNil)
		err = Bind(ctx, &n)
		So(err, ShouldNotBeNil)

		err = Bind(ctx, &f)
		So(err, ShouldBeNil)

		ctx = makeContext(url.Values{
			"ID": []string{"abc"},
		})
		err = Bind(ctx, &f)
		So(err, ShouldNotBeNil)

		data := url.Values{
			"ID":          []string{"1234"},
			"username":    []string{"jiazhoulvke"},
			"age":         []string{"50"},
			"weight":      []string{"60"},
			"EnglishName": []string{"niko"},
			"Number":      []string{"24"},
			"AliasName":   []string{"jiazhoulvke1"},
			"AliasName2":  []string{"jiazhoulvke_1"},
			"Email":       []string{"foo@bar.com"},
			"LastIP":      []string{"127.0.0.1"},
			"Mobile":      []string{"13812345678"},
			"Tel":         []string{"07551234567"},
			"Phone":       []string{"13812345678"},
			"IDCard":      []string{"43211234567890123X"},
			"Height":      []string{"1.89"},
			"Count":       []string{"66666"},
			"IsMan":       []string{"yes"},
			"Birthday":    []string{"2007-12-13"},
		}
		ctx = makeContext(data)
		birthday, err := time.Parse("2006-01-02", "2007-12-13")
		So(err, ShouldBeNil)
		var a = foo{
			ID:            1234,
			UserName:      "jiazhoulvke",
			Age:           50,
			Weight:        60,
			EnglishName:   "niko",
			Number:        "24",
			AliasName:     "jiazhoulvke1",
			AliasName2:    "jiazhoulvke_1",
			Email:         "foo@bar.com",
			LastIP:        "127.0.0.1",
			Mobile:        "13812345678",
			Tel:           "07551234567",
			Phone:         "13812345678",
			IDCard:        "43211234567890123X",
			Height:        1.89,
			Count:         66666,
			IsMan:         true,
			Birthday:      birthday,
			DefaultString: "foobar",
			DefaultInt:    42,
		}
		err = Bind(ctx, &f)
		So(err, ShouldBeNil)
		isEqual := reflect.DeepEqual(&f, &a)
		So(isEqual, ShouldBeTrue)
	})

	Convey("测试越界情况", t, func() {
		Convey("测试浮点数", func() {
			ctx := makeContext(url.Values{
				"N": []string{fmt.Sprintf("%v", float64(math.MaxFloat64))},
			})
			var err error

			err = Bind(ctx, &struct {
				N float32
			}{})
			So(err, ShouldNotBeNil)
		})

		Convey("测试无符号整数", func() {
			ctx := makeContext(url.Values{
				"N": []string{fmt.Sprintf("%v", uint64(math.MaxUint64))},
			})
			var err error

			err = Bind(ctx, &struct {
				N uint8
			}{})
			So(err, ShouldNotBeNil)

			err = Bind(ctx, &struct {
				N uint16
			}{})
			So(err, ShouldNotBeNil)

			err = Bind(ctx, &struct {
				N uint32
			}{})
			So(err, ShouldNotBeNil)
		})

		Convey("测试有符号整数", func() {

			Convey("测试最大值", func() {
				ctx := makeContext(url.Values{
					"N": []string{fmt.Sprintf("%v", int64(math.MaxInt64))},
				})
				var err error
				//int
				if strconv.IntSize == 32 {
					err = Bind(ctx, &struct {
						N int
					}{})
					So(err, ShouldNotBeNil)
				}
				//int8
				err = Bind(ctx, &struct {
					N int8
				}{})
				So(err, ShouldNotBeNil)
				//int16
				err = Bind(ctx, &struct {
					N int16
				}{})
				So(err, ShouldNotBeNil)
				//int32
				err = Bind(ctx, &struct {
					N int32
				}{})
				So(err, ShouldNotBeNil)
			})

			Convey("测试最小值", func() {
				ctx := makeContext(url.Values{
					"N": []string{fmt.Sprintf("%v", int64(math.MinInt64))},
				})
				var err error
				//int
				if strconv.IntSize == 32 {
					err = Bind(ctx, &struct {
						N int
					}{})
					So(err, ShouldNotBeNil)
				}
				//int8
				err = Bind(ctx, &struct {
					N int8
				}{})
				So(err, ShouldNotBeNil)
				//int16
				err = Bind(ctx, &struct {
					N int16
				}{})
				So(err, ShouldNotBeNil)
				//int32
				err = Bind(ctx, &struct {
					N int32
				}{})
				So(err, ShouldNotBeNil)
			})
		})

		Convey("测试多重struct", func() {
			type (
				foo struct {
					FooInt    int
					FooInt8   int8
					FooInt16  int16
					FooInt32  int32
					FooInt64  int64
					FooUint   uint
					FooUint8  uint8
					FooUint16 uint16
					FooUint32 uint32
					FooUint64 uint64
				}
				bar struct {
					foo
					BarInt int
				}
			)

			var a = bar{}
			var err error
			data := url.Values{
				"FooInt":    []string{"7"},
				"FooInt8":   []string{"7"},
				"FooInt16":  []string{"7"},
				"FooInt32":  []string{"7"},
				"FooInt64":  []string{"7"},
				"FooUint":   []string{"7"},
				"FooUint8":  []string{"7"},
				"FooUint16": []string{"7"},
				"FooUint32": []string{"7"},
				"FooUint64": []string{"7"},
				"BarInt":    []string{"4321"},
			}
			ctx := makeContext(data)
			err = Bind(ctx, &a)
			So(err, ShouldBeNil)

			var b = bar{}
			b.foo = foo{
				FooInt:    7,
				FooInt8:   7,
				FooInt16:  7,
				FooInt32:  7,
				FooInt64:  7,
				FooUint:   7,
				FooUint8:  7,
				FooUint16: 7,
				FooUint32: 7,
				FooUint64: 7,
			}
			b.BarInt = 4321
			isEqual := reflect.DeepEqual(&a, &b)
			So(isEqual, ShouldBeTrue)
		})

		Convey("测试错误的输入", func() {
			type (
				foo struct {
					TestInt int
				}

				bar struct {
					TestFloat float64
				}
			)

			data := url.Values{
				"TestInt": []string{"abc"},
			}
			var a = foo{}
			ctx := makeContext(data)
			err := Bind(ctx, &a)
			So(err, ShouldNotBeNil)

			data = url.Values{
				"TestFloat": []string{"abc"},
			}
			var b = bar{}
			ctx = makeContext(data)
			err = Bind(ctx, &b)
			So(err, ShouldNotBeNil)
		})

	})
}
