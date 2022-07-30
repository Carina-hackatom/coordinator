package base

import (
	"github.com/Carina-hackatom/coordinator/utils"
	ut "github.com/Carina-hackatom/coordinator/utils/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"io"
)

func LoadClientPubInfo(ctx client.Context, keyname string) keyring.Info {
	accInfo, err := ctx.Keyring.Key(keyname)
	utils.CheckErr(err, "Cannot load client with key", ut.EXIT)
	return accInfo
}

func MakeClientWithNewAcc(ctx client.Context, keyname string, mnemonic string, bip44path string, algo keyring.SignatureAlgo) keyring.Info {
	accInfo := CreateAcc(ctx, keyname, mnemonic, bip44path, algo)
	return accInfo
}

func CreateAcc(ctx client.Context, keyname string, mnemonic string, bip44path string, algo keyring.SignatureAlgo) keyring.Info {
	info, err := ctx.Keyring.NewAccount(keyname, mnemonic, "", bip44path, algo)
	utils.CheckErr(err, "Cannot create account with those arguments. Check it!", ut.EXIT)
	return info
}

func GetPrivKey(ctx client.Context, keyname string) cryptotypes.PrivKey {
	privArmor, _ := ctx.Keyring.ExportPrivKeyArmor(keyname, "")
	novaPrivRaw, _, _ := crypto.UnarmorDecryptPrivKey(privArmor, "")
	return novaPrivRaw
}

func MakeContext(mb module.BasicManager, from string, tmRPC string, chainID string, root string, backend string, userInput io.Reader, userOutput io.Writer, genOnly bool) client.Context {
	encCfg := MakeEncodingConfig(mb)
	initClientCtx := client.Context{}.
		WithSimulation(false).
		WithSkipConfirmation(true).
		WithSignModeStr(flags.SignModeDirect).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithFrom(from).
		WithNodeURI(tmRPC).
		WithChainID(chainID).
		WithHomeDir(root).
		WithKeyringDir(root).
		WithInput(userInput).
		WithOutput(userOutput).
		WithGenerateOnly(genOnly).
		WithCodec(encCfg.Marshaler).
		WithInterfaceRegistry(encCfg.InterfaceRegistry).
		WithTxConfig(encCfg.TxConfig).
		WithLegacyAmino(encCfg.Amino)

	kb := MakeKeyring(initClientCtx, backend)
	tmClient, err := client.NewClientFromNode(tmRPC)
	utils.CheckErr(err, "-> Cannot set node client", 0)

	return initClientCtx.
		WithKeyring(kb).
		WithClient(tmClient)

}

func AddMoreFromInfo(ctx client.Context) client.Context {
	fromAddr, fromName, _, err := client.GetFromFields(ctx.Keyring, ctx.From, ctx.GenerateOnly)
	utils.CheckErr(err, "cannot get info from keyring", 0)
	ctx = ctx.WithFromAddress(fromAddr).WithFromName(fromName)
	return ctx
}
