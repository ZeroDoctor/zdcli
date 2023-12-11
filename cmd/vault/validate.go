package vault

import (
	"errors"
	"fmt"

	v "github.com/go-playground/validator/v10"
	"github.com/zerodoctor/zdcli/config"
)

type vPasswords struct {
	Pass        string `validate:"required,min=10,containsany=0123456789!@#$^*,eqfield=ConfirmPass"`
	ConfirmPass string
}

func validatePasswords(pass, confirmPass string) error {
	return v.New().Struct(&vPasswords{
		Pass:        pass,
		ConfirmPass: confirmPass,
	})
}

type vUsername struct {
	Username string `validate:"required,min=3,alpha"`
}

func validateUserName(userName string) error {
	return v.New().Struct(&vUsername{
		Username: userName,
	})
}

func formatErrors(errors []error, suffix string) string {
	str := ""

	for i := range errors {
		str += errors[i].Error() + suffix
	}

	return str
}

type VFlag int

func (b VFlag) Has(f VFlag) bool { return b&f != 0 }

const (
	VEndpoint VFlag = 1 << iota
	VToken
)

var ErrMissingVaultEndpoint error = errors.New("missing vault endpoint")
var ErrMissingVaultMainToken error = errors.New("missing vault main login token")

func validate(flag VFlag, cfg *config.Config) error {
	var errs []any

	if flag.Has(VEndpoint) && cfg.VaultEndpoint == "" {
		return ErrMissingVaultEndpoint
	}

	if flag.Has(VToken) {
		if _, ok := cfg.VaultTokens[cfg.VaultUser]; !ok {
			errs = append(errs, ErrMissingVaultMainToken)
		}
	}

	var format string
	for range errs {
		format += "[error=%w] "
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf(format, errs...)
	}

	return err
}
