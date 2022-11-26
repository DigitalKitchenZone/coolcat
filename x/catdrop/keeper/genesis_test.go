package keeper_test

import (
	"github.com/DigitalKitchenLabs/coolcat/v1/x/catdrop/types"
)

func (s *KeeperTestSuite) TestExportGenesis() {
	app, ctx := s.app, s.ctx
	app.CatdropKeeper.InitGenesis(ctx, *types.DefaultGenesis())
	// app.ClaimKeeper.SetParams(ctx, types.DefaultParams())
	exported := app.CatdropKeeper.ExportGenesis(ctx)
	params := types.DefaultParams()
	params.AirdropStartTime = ctx.BlockTime()
	s.Require().Equal(params.AllowedClaimers, exported.Params.AllowedClaimers)
}
