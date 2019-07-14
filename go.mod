module github.com/hyoungsungkim/nameservice

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190616100639-18415eedaf25
	github.com/cosmos/sdk-application-tutorial v0.0.0-20190707093158-91a381865095
	github.com/gorilla/mux v1.7.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.5
	google.golang.org/grpc v1.19.1
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
