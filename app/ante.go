package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcante "github.com/cosmos/ibc-go/v3/modules/core/ante"
	"github.com/cosmos/ibc-go/v3/modules/core/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions

	IBCKeeper  			*keeper.Keeper
	WasmConfig        *wasmTypes.WasmConfig
	TXCounterStoreKey sdk.StoreKey
}

 type MinCommissionDecorator struct{}

 func NewMinCommissionDecorator() MinCommissionDecorator {
 	return MinCommissionDecorator{}
 }

 func (MinCommissionDecorator) AnteHandle(
 	ctx sdk.Context, tx sdk.Tx,
 	simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
 	msgs := tx.GetMsgs()
 	minCommissionRate := sdk.NewDecWithPrec(5, 2)
 	for _, m := range msgs {
 		switch msg := m.(type) {
 		case *stakingtypes.MsgCreateValidator:
 			// prevent new validators joining the set with
 			// commission set below 5%
 			c := msg.Commission
 			if c.Rate.LT(minCommissionRate) {
 				return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "commission can't be lower than 5%")
 			}
 		case *stakingtypes.MsgEditValidator:
 			// if commission rate is nil, it means only
 			// other fields are affected - skip
 			if msg.CommissionRate == nil {
 				continue
 			}
 			if msg.CommissionRate.LT(minCommissionRate) {
 				return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "commission can't be lower than 5%")
 			}
 		default:
 			continue
 		}
 	}
 	return next(ctx, tx, simulate)
 }

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.WasmConfig == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "wasm config is required for ante builder")
	}
	if options.TXCounterStoreKey == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "tx counter key is required for ante builder")
	}

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreKey),
		ante.NewRejectExtensionOptionsDecorator(),
		NewMinCommissionDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewAnteDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
