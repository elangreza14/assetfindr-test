package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type SuccessResponse struct {
	Data   any    `json:"data,omitempty"`
	Result string `json:"result"`
	Err    any    `json:"error,omitempty"`
}

type SuccessResponsePlain struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Result string `json:"result"`
	Err    any    `json:"error,omitempty"`
}

func NewBaseResponse(data any, err error) any {

	if message, ok := data.(string); ok && err == nil {
		return SuccessResponsePlain{
			Result: message,
		}
	}

	if data == nil && err == nil {
		return SuccessResponsePlain{
			Result: "ok",
		}
	}

	if err != nil {
		errRes := ErrorResponse{
			Result: "error",
		}

		errs := validateErrorStruct(err)
		if len(errs) > 0 {
			errRes.Result = "errors"
			errRes.Err = errs
			return errRes
		}

		errRes.Err = err.Error()
		return errRes
	}

	return SuccessResponse{
		Data:   data,
		Result: "ok",
	}
}

type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "gt":
		return "Should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func validateErrorStruct(err error) []ErrorField {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorField, len(ve))
		for i, fe := range ve {
			out[i] = ErrorField{fe.Field(), getErrorMsg(fe)}
		}
		return out
	}

	return []ErrorField{}
}
