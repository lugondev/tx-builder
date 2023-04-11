package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"

	"github.com/lugondev/tx-builder/pkg/errors"
)

func UnmarshalBody(body io.Reader, req interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return errors.InvalidFormatError("failed to decode request body").AppendReason(err.Error())
	}

	err = GetValidator().Struct(req)
	if err != nil {
		if ves, ok := err.(validator.ValidationErrors); ok {
			var errMessage string
			for _, fe := range ves {
				errMessage += fmt.Sprintf("field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
			}

			return errors.InvalidParameterError("invalid body").AppendReason(errMessage)
		}

		return errors.InvalidFormatError("invalid body")
	}

	return nil
}

func UnmarshalInterface(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dst)
}
