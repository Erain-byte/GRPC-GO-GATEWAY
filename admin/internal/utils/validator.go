package utils

import (
	"github.com/go-playground/validator/v10"
)

// Validate 通用校验函数，接收任意结构体并返回校验错误
func Validate(data interface{}) error {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		return err
	}
	return nil
}
