package validator

import infra "github.com/ABA-Developer/go-backend-template/internal/infrastructure/validator"

type Validator = infra.Validator

func New(engine *infra.Validator) *Validator {
	return engine
}
