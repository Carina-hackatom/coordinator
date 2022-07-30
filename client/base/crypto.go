package base

import (
	"github.com/Carina-hackatom/coordinator/utils"
	"github.com/Carina-hackatom/coordinator/utils/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func MakeKeyring(ctx client.Context, backend string) keyring.Keyring {
	kb, err := newKeyringFromBackend(ctx, backend)
	utils.CheckErr(err, "Cannot generate keyring instance", types.EXIT)
	return kb
}

func newKeyringFromBackend(ctx client.Context, backend string) (keyring.Keyring, error) {
	if ctx.GenerateOnly || ctx.Simulate {
		return keyring.New(sdktypes.KeyringServiceName(), keyring.BackendMemory, ctx.KeyringDir, ctx.Input, ctx.KeyringOptions...)
	}

	return keyring.New(sdktypes.KeyringServiceName(), backend, ctx.KeyringDir, ctx.Input, ctx.KeyringOptions...)
}
