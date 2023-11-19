package misc

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var phoneValidator validator.Func = func(fl validator.FieldLevel) bool {
	fieldValue := fl.Field()
	re := regexp.MustCompile(`^[0-9() +-]{10,}$`)
	if fieldValue.Kind().String() != "string" {
		log.Printf("field %q is not a string", fl.FieldName())
		return false
	}
	if !re.MatchString(fieldValue.String()) {
		log.Printf("field %q doesn't match regexp %q", fl.FieldName(), re.String())
		return false
	}
	return true
}

func RegisterValidators(validate *validator.Validate) {
	validators := map[string]validator.Func{
		"phone": phoneValidator,
	}
	for k, v := range validators {
		if err := validate.RegisterValidation(k, v); err != nil {
			panic(err)
		}
	}
}
