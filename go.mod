module github.com/Carina-hackatom/coordinator

go 1.18

require (
	github.com/cosmos/cosmos-sdk v0.45.4
)

replace (
	github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
	github.com/cosmos/ibc-go/v3 => github.com/Carina-hackatom/ibc-go/v3 v3.0.0-novachain
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)