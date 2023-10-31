package util

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

func NewValidate() *Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("label")
	})
	uni := ut.New(en.New(), zh.New())
	trans, _ := uni.GetTranslator("zh")
	_ = zhtranslations.RegisterDefaultTranslations(validate, trans)

	return &Validate{validator: validate, trans: trans}
}

type CustomValidateFunc struct {
	Handle  validator.Func
	TagName string
	Message string
}

type Validate struct {
	trans     ut.Translator
	validator *validator.Validate
}

func (v *Validate) Struct(s interface{}) error {
	err := v.validator.Struct(s)
	errs := translateError(err, v.trans)
	if errs == nil {
		return nil
	}
	var msg []string
	for _, v := range errs {
		msg = append(msg, v.Error())
	}
	return fmt.Errorf("%s", strings.Join(msg, "\n"))
}

func (v *Validate) Var(field interface{}, tag, label string) error {
	err := v.validator.Var(field, tag)
	errs := translateError(err, v.trans)
	if errs == nil {
		return nil
	}
	return fmt.Errorf("%s%s", label, errs[0])
}

func (v *Validate) VarCtx(ctx context.Context, field interface{}, tag, label string) error {
	err := v.validator.VarCtx(ctx, field, tag)
	errs := translateError(err, v.trans)
	if errs == nil {
		return nil
	}
	return fmt.Errorf("%s%s", label, errs[0])
}

func (v *Validate) RegisterValidation(customValidates ...CustomValidateFunc) {
	for _, customValidate := range customValidates {
		// 注册自定义函数
		_ = v.validator.RegisterValidation(customValidate.TagName, customValidate.Handle)
		// 根据提供的标记注册翻译
		_ = v.validator.RegisterTranslation(customValidate.TagName, v.trans, func(ut ut.Translator) error {
			return ut.Add(customValidate.TagName, customValidate.Message, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(customValidate.TagName, fe.Field(), fe.Field())
			return t
		})
	}
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}
