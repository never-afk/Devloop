package util

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

// Validator 验证器
type Validator struct {
	once     sync.Once
	validate *validator.Validate
}

var (
	_                   binding.StructValidator = &Validator{}
	universalTranslator *ut.UniversalTranslator
)

func Locales(l ...locales.Translator) {
	universalTranslator = ut.New(en.New(), l...)
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *Validator) ValidateStruct(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

// Engine 获取验证器
func (v *Validator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

// lazyInit 延迟初始化
func (v *Validator) lazyInit() {
	v.once.Do(func() {
		if universalTranslator == nil {
			universalTranslator = ut.New(en.New(), en.New(), zh.New())

		}
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// 获取form tag
		v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {

			name := field.Tag.Get("form")
			if name != "" {
				return name
			}

			name = field.Tag.Get("json")
			if name != "" {
				return name
			}
			return field.Name
		})

		for tag, languages := range DefaultTransMessage {
			registerTranslation(v.validate, tag, languages)
		}
	})
}

func registerTranslation(validate *validator.Validate, tag string, languages map[string]string) {
	translationFn := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Param(), fe.Field())
		return t
	}
	for k, v := range languages {
		if trans, ok := universalTranslator.GetTranslator(k); ok {
			_ = validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
				return ut.Add(tag, v, true)
			}, translationFn)
		}
	}
}

func GetErrors(err error) map[string]string {
	result := make(map[string]string)
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		for _, value := range errs {
			if trans, ok := universalTranslator.GetTranslator(`en`); ok {
				result[value.Field()] = value.Translate(trans)
			}
		}
	}
	return result
}
