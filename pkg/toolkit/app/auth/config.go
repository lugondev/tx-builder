package auth

import (
	"github.com/lugondev/tx-builder/pkg/toolkit/app/auth/jwt/jose"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/auth/key"
	"github.com/spf13/pflag"
)

func Flags(f *pflag.FlagSet) {
	key.Flags(f)
	jose.Flags(f)
}
