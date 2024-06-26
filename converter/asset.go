package converter

import (
	"time"

	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertTrustLineEntry(e xdr.TrustLineEntry) (TrustLineEntry, error) {
	var result TrustLineEntry
	accountId, err := ConvertAccountId(e.AccountId)
	if err != nil {
		return result, err
	}

	asset, err := ConvertTrustLineAsset(e.Asset)
	if err != nil {
		return result, err
	}

	ext := ConvertTrustLineEntryExt(e.Ext)

	result.AccountId = accountId
	result.Asset = asset
	result.Balance = int64(e.Balance)
	result.Limit = int64(e.Limit)
	result.Flags = uint32(e.Flags)
	result.Ext = ext

	return result, nil
}

func ConvertTrustLineEntryExt(e xdr.TrustLineEntryExt) TrustLineEntryExt {
	var v1 TrustLineEntryV1
	if e.V1 != nil {
		v1 = ConvertTrustLineEntryV1(*e.V1)
	}

	return TrustLineEntryExt{
		V:  e.V,
		V1: &v1,
	}
}

func ConvertTrustLineEntryV1(e xdr.TrustLineEntryV1) TrustLineEntryV1 {
	return TrustLineEntryV1{
		Liabilities: ConvertLiabilities(e.Liabilities),
		Ext:         ConvertTrustLineEntryV1Ext(e.Ext),
	}
}

func ConvertTrustLineEntryV1Ext(e xdr.TrustLineEntryV1Ext) TrustLineEntryV1Ext {
	var v2 TrustLineEntryExtensionV2
	if e.V2 != nil {
		v2 = ConvertTrustLineEntryExtensionV2(*e.V2)
	}

	return TrustLineEntryV1Ext{
		V:  e.V,
		V2: &v2,
	}
}

func ConvertTrustLineEntryExtensionV2(e xdr.TrustLineEntryExtensionV2) TrustLineEntryExtensionV2 {
	return TrustLineEntryExtensionV2{
		LiquidityPoolUseCount: int32(e.LiquidityPoolUseCount),
		Ext:                   ConvertTrustLineEntryExtensionV2Ext(e.Ext),
	}
}

func ConvertTrustLineEntryExtensionV2Ext(e xdr.TrustLineEntryExtensionV2Ext) TrustLineEntryExtensionV2Ext {
	return TrustLineEntryExtensionV2Ext{V: e.V}
}

// TODO: testing
func ConvertAsset(as xdr.Asset) (Asset, error) {
	var result Asset
	switch as.Type {
	case xdr.AssetTypeAssetTypeNative:
		result.AssetType = "native"
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		result.AssetType = "alphanum4"
		result.AssetCode = as.AlphaNum4.AssetCode[:]

		issuer, err := ConvertAccountId(as.AlphaNum4.Issuer)
		if err != nil {
			return result, err
		}
		result.Issuer = issuer

		return result, nil
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		result.AssetType = "alphanum12"
		result.AssetCode = as.AlphaNum12.AssetCode[:]

		issuer, err := ConvertAccountId(as.AlphaNum12.Issuer)
		if err != nil {
			return result, err
		}
		result.Issuer = issuer

		return result, nil
	case xdr.AssetTypeAssetTypePoolShare:
		result.AssetType = "poolshare"
	default:
		return result, errors.Errorf("unsupported asset type %v", as.Type)
	}
	return result, nil
}

// TODO: testing
func ConvertTrustLineAsset(a xdr.TrustLineAsset) (TrustLineAsset, error) {
	var result TrustLineAsset
	asset, err := ConvertAsset(a.ToAsset())
	if err != nil {
		return result, err
	}
	result.Asset = &asset

	if a.LiquidityPoolId != nil {
		xdrLpId := xdr.Hash(*a.LiquidityPoolId)
		lpId := PoolId(xdrLpId[:])
		result.LiquidityPoolId = &lpId
	}

	return result, nil
}

// TODO: testing
func ConvertLiquidityPoolConstantProductParameters(
	lpcpp xdr.LiquidityPoolConstantProductParameters,
) (LiquidityPoolConstantProductParameters, error) {
	var result LiquidityPoolConstantProductParameters

	assetA, err := ConvertAsset(lpcpp.AssetA)
	if err != nil {
		return result, err
	}

	assetB, err := ConvertAsset(lpcpp.AssetB)
	if err != nil {
		return result, err
	}

	result.AssetA = assetA
	result.AssetB = assetB
	result.Fee = int32(lpcpp.Fee)

	return result, nil
}

// TODO: testing
func ConvertLiquidityPoolParameters(lpp xdr.LiquidityPoolParameters) (LiquidityPoolParameters, error) {
	var result LiquidityPoolParameters

	switch lpp.Type {
	case xdr.LiquidityPoolTypeLiquidityPoolConstantProduct:
		lpcpp, err := ConvertLiquidityPoolConstantProductParameters(*lpp.ConstantProduct)
		if err != nil {
			return result, err
		}

		result.ConstantProduct = &lpcpp

		return result, nil
	}

	return result, errors.Errorf("invalid liquidity pool parameters type %v", lpp.Type)
}

// TODO: testing
func ConvertChangeTrustAsset(ta xdr.ChangeTrustAsset) (ChangeTrustAsset, error) {
	var result ChangeTrustAsset

	asset, err := ConvertAsset(ta.ToAsset())
	if err != nil {
		return result, err
	}
	result.Asset = &asset

	var liquidityPool LiquidityPoolParameters
	if ta.LiquidityPool != nil {
		var err error
		liquidityPool, err = ConvertLiquidityPoolParameters(*ta.LiquidityPool)
		if err != nil {
			return result, err
		}
	}
	result.LiquidityPool = &liquidityPool

	return result, nil
}

func ConvertClaimPredicates(inp []xdr.ClaimPredicate) ([]ClaimPredicate, error) {
	parts := make([]ClaimPredicate, len(inp))
	for i, pred := range inp {
		converted, err := ConvertClaimPredicate(pred)
		if err != nil {
			return parts, err
		}
		parts[i] = converted
	}
	return parts, nil
}

// TODO: testing
func ConvertClaimPredicate(cp xdr.ClaimPredicate) (ClaimPredicate, error) {
	var result ClaimPredicate

	switch cp.Type {
	case xdr.ClaimPredicateTypeClaimPredicateUnconditional:
		// void
		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateAnd:
		andPredicates, err := ConvertClaimPredicates(*cp.AndPredicates)
		if err != nil {
			return result, err
		}
		result.AndPredicates = &andPredicates

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateOr:
		orPredicates, err := ConvertClaimPredicates(*cp.OrPredicates)
		if err != nil {
			return result, err
		}
		result.OrPredicates = &orPredicates

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateNot:
		xdrNotPredicate, ok := cp.GetNotPredicate()
		if !ok {
			return result, errors.Errorf("invalid type ClaimPredicateTypeClaimPredicateNot")
		}

		notPredicate, err := ConvertClaimPredicate(*xdrNotPredicate)
		if err != nil {
			return result, err
		}
		result.NotPredicate = &notPredicate

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime:
		absBeforeEpoch := int64(*cp.AbsBefore)
		absBefore := time.Unix(absBeforeEpoch, 0).UTC()

		result.AbsBefore = &absBefore
		result.AbsBeforeEpoch = &absBeforeEpoch

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateBeforeRelativeTime:
		relBefore := int64(*cp.RelBefore)
		result.RelBefore = &relBefore

		return result, nil
	}

	return result, errors.Errorf("invalid ClaimPredicate type %v", cp.Type)
}

// TODO: testing
func ConvertClaimant(c xdr.Claimant) (Claimant, error) {
	var result Claimant

	switch c.Type {
	case xdr.ClaimantTypeClaimantTypeV0:
		xdrV0 := c.V0

		destination, err := ConvertAccountId(xdrV0.Destination)
		if err != nil {
			return result, err
		}

		predicate, err := ConvertClaimPredicate(c.V0.Predicate)
		if err != nil {
			return result, err
		}

		v0 := &ClaimantV0{
			Destination: destination,
			Predicate:   predicate,
		}
		result.V0 = v0

		return result, nil
	}

	return result, errors.Errorf("invalid claimant type %v", c.Type)
}

func ConvertConvertClaimableBalanceEntry(e xdr.ClaimableBalanceEntry) (ClaimableBalanceEntry, error) {
	var result ClaimableBalanceEntry

	balanceId, err := ConvertClaimableBalanceId(e.BalanceId)
	if err != nil {
		return result, err
	}

	var claimants []Claimant
	for _, xdrClaimant := range e.Claimants {
		claimant, err := ConvertClaimant(xdrClaimant)
		if err != nil {
			return result, err
		}

		claimants = append(claimants, claimant)
	}

	asset, err := ConvertAsset(e.Asset)
	if err != nil {
		return result, err
	}

	ext := ConvertClaimableBalanceEntryExt(e.Ext)

	result.BalanceId = balanceId
	result.Claimants = claimants
	result.Asset = asset
	result.Amount = int64(e.Amount)
	result.Ext = ext

	return result, nil
}

func ConvertClaimableBalanceEntryExt(e xdr.ClaimableBalanceEntryExt) ClaimableBalanceEntryExt {
	var v1 ClaimableBalanceEntryExtensionV1
	if e.V1 != nil {
		v1 = ConvertClaimableBalanceEntryExtensionV1(*e.V1)
	}

	return ClaimableBalanceEntryExt{
		V:  e.V,
		V1: &v1,
	}
}

func ConvertClaimableBalanceEntryExtensionV1(e xdr.ClaimableBalanceEntryExtensionV1) ClaimableBalanceEntryExtensionV1 {
	return ClaimableBalanceEntryExtensionV1{
		Flags: uint32(e.Flags),
		Ext:   ConvertClaimableBalanceEntryExtensionV1Ext(e.Ext),
	}
}

func ConvertClaimableBalanceEntryExtensionV1Ext(e xdr.ClaimableBalanceEntryExtensionV1Ext) ClaimableBalanceEntryExtensionV1Ext {
	return ClaimableBalanceEntryExtensionV1Ext{V: e.V}
}

// TODO: testing
func ConvertClaimableBalanceId(id xdr.ClaimableBalanceId) (ClaimableBalanceId, error) {
	var result ClaimableBalanceId

	switch id.Type {
	case xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0:
		v0 := (*id.V0).HexString()
		result.V0 = &v0
		return result, nil
	}

	return result, errors.Errorf("invalid ClaimableBalanceId type %v", id.Type)
}

// TODO: testing
func ConvertPrice(p xdr.Price) Price {
	return Price{
		N: int32(p.N),
		D: int32(p.D),
	}
}

func ConvertPathPaymentStrictReceiveResultSuccess(r xdr.PathPaymentStrictReceiveResultSuccess) (PathPaymentStrictReceiveResultSuccess, error) {
	var result PathPaymentStrictReceiveResultSuccess

	var offers []ClaimAtom
	for _, xdrOffer := range r.Offers {
		offer, err := ConvertClaimAtom(xdrOffer)
		if err != nil {
			return result, err
		}

		offers = append(offers, offer)
	}

	last, err := ConvertSimplePaymentResult(r.Last)
	if err != nil {
		return result, err
	}

	result.Offers = offers
	result.Last = last

	return result, nil
}

func ConvertPathPaymentStrictSendResultSuccess(r xdr.PathPaymentStrictSendResultSuccess) (PathPaymentStrictSendResultSuccess, error) {
	var result PathPaymentStrictSendResultSuccess

	var offers []ClaimAtom
	for _, xdrOffer := range r.Offers {
		offer, err := ConvertClaimAtom(xdrOffer)
		if err != nil {
			return result, err
		}

		offers = append(offers, offer)
	}

	last, err := ConvertSimplePaymentResult(r.Last)
	if err != nil {
		return result, err
	}

	result.Offers = offers
	result.Last = last

	return result, nil
}

func ConvertClaimAtom(c xdr.ClaimAtom) (ClaimAtom, error) {
	var result ClaimAtom

	switch c.Type {
	case xdr.ClaimAtomTypeClaimAtomTypeV0:
		v0, err := ConvertClaimOfferAtomV0(*c.V0)
		if err != nil {
			return result, err
		}

		result.V0 = &v0

		return result, nil
	case xdr.ClaimAtomTypeClaimAtomTypeOrderBook:
		orderBook, err := ConvertClaimOfferAtom(*c.OrderBook)
		if err != nil {
			return result, err
		}

		result.OrderBook = &orderBook

		return result, nil
	case xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool:
		lp, err := ConvertClaimLiquidityAtom(*c.LiquidityPool)
		if err != nil {
			return result, err
		}

		result.LiquidityPool = &lp

		return result, nil
	}

	return result, errors.Errorf("invalid ConvertClaimAtom type %v", c.Type)
}

func ConvertClaimOfferAtomV0(c xdr.ClaimOfferAtomV0) (ClaimOfferAtomV0, error) {
	var result ClaimOfferAtomV0

	sellerEd25519, err := ConvertEd25519(c.SellerEd25519)
	if err != nil {
		return result, err
	}

	assetSold, err := ConvertAsset(c.AssetSold)
	if err != nil {
		return result, err
	}

	assetBought, err := ConvertAsset(c.AssetBought)
	if err != nil {
		return result, err
	}

	result.SellerEd25519 = sellerEd25519
	result.OfferId = int64(c.OfferId)
	result.AssetSold = assetSold
	result.AmountSold = int64(c.AmountSold)
	result.AssetBought = assetBought
	result.AmountBought = int64(c.AmountBought)

	return result, nil
}

func ConvertClaimOfferAtom(c xdr.ClaimOfferAtom) (ClaimOfferAtom, error) {
	var result ClaimOfferAtom

	sellerId, err := ConvertAccountId(c.SellerId)
	if err != nil {
		return result, err
	}

	assetSold, err := ConvertAsset(c.AssetSold)
	if err != nil {
		return result, err
	}

	assetBought, err := ConvertAsset(c.AssetBought)
	if err != nil {
		return result, err
	}

	result.SellerId = sellerId
	result.OfferId = int64(c.OfferId)
	result.AssetSold = assetSold
	result.AmountSold = int64(c.AmountSold)
	result.AssetBought = assetBought
	result.AmountBought = int64(c.AmountBought)

	return result, nil
}

func ConvertClaimLiquidityAtom(c xdr.ClaimLiquidityAtom) (ClaimLiquidityAtom, error) {
	var result ClaimLiquidityAtom

	xdrPoolId := xdr.Hash(c.LiquidityPoolId)
	poolId := PoolId(xdrPoolId[:])

	assetSold, err := ConvertAsset(c.AssetSold)
	if err != nil {
		return result, err
	}

	assetBought, err := ConvertAsset(c.AssetBought)
	if err != nil {
		return result, err
	}

	result.LiquidityPoolId = poolId
	result.AssetSold = assetSold
	result.AmountSold = int64(c.AmountSold)
	result.AssetBought = assetBought
	result.AmountBought = int64(c.AmountBought)

	return result, nil
}

func ConvertSimplePaymentResult(r xdr.SimplePaymentResult) (SimplePaymentResult, error) {
	var result SimplePaymentResult

	destination, err := ConvertAccountId(r.Destination)
	if err != nil {
		return result, err
	}

	asset, err := ConvertAsset(r.Asset)
	if err != nil {
		return result, err
	}

	result.Destination = destination
	result.Asset = asset
	result.Amount = int64(r.Amount)

	return result, nil
}

func ConvertManageOfferSuccessResult(r xdr.ManageOfferSuccessResult) (ManageOfferSuccessResult, error) {
	var result ManageOfferSuccessResult

	var offersClaimed []ClaimAtom
	for _, xdrOffer := range r.OffersClaimed {
		offer, err := ConvertClaimAtom(xdrOffer)
		if err != nil {
			return result, err
		}

		offersClaimed = append(offersClaimed, offer)
	}

	offer, err := ConvertManageOfferSuccessResultOffer(r.Offer)
	if err != nil {
		return result, err
	}
	result.Offer = offer
	result.OffersClaimed = offersClaimed

	return result, nil
}

func ConvertManageOfferSuccessResultOffer(r xdr.ManageOfferSuccessResultOffer) (ManageOfferSuccessResultOffer, error) {
	var result ManageOfferSuccessResultOffer

	result.Effect = int32(r.Effect)

	var offer OfferEntry
	var err error
	if r.Offer != nil {
		offer, err = ConvertOfferEntry(*r.Offer)
		if err != nil {
			return result, err
		}
	}
	result.Offer = &offer

	return result, nil
}

func ConvertOfferEntry(e xdr.OfferEntry) (OfferEntry, error) {
	var result OfferEntry

	sellerId, err := ConvertAccountId(e.SellerId)
	if err != nil {
		return result, err
	}

	selling, err := ConvertAsset(e.Selling)
	if err != nil {
		return result, err
	}

	buying, err := ConvertAsset(e.Buying)
	if err != nil {
		return result, err
	}

	price := ConvertPrice(e.Price)

	result.SellerId = sellerId
	result.OfferId = int64(e.OfferId)
	result.Selling = selling
	result.Buying = buying
	result.Amount = int64(e.Amount)
	result.Price = price
	result.Flags = uint32(e.Flags)
	result.Ext = ConvertOfferEntryExt(e.Ext)

	return result, nil
}

func ConvertOfferEntryExt(e xdr.OfferEntryExt) OfferEntryExt {
	return OfferEntryExt{V: e.V}
}

func ConvertInflationPayout(i xdr.InflationPayout) InflationPayout {
	destination, _ := ConvertAccountId(i.Destination)
	return InflationPayout{
		Destination: destination,
		Amount:      int64(i.Amount),
	}
}

func ConvertLiquidityPoolEntry(e xdr.LiquidityPoolEntry) (LiquidityPoolEntry, error) {
	var result LiquidityPoolEntry
	body, err := ConvertLiquidityPoolEntryBody(e.Body)
	if err != nil {
		return result, err
	}

	result.LiquidityPoolId = PoolId(e.LiquidityPoolId[:])
	result.Body = body

	return result, nil
}

func ConvertLiquidityPoolEntryBody(b xdr.LiquidityPoolEntryBody) (LiquidityPoolEntryBody, error) {
	var result LiquidityPoolEntryBody

	switch b.Type {
	case xdr.LiquidityPoolTypeLiquidityPoolConstantProduct:
		constProduct, err := ConvertLiquidityPoolEntryConstantProduct(*b.ConstantProduct)
		if err != nil {
			return result, err
		}

		result.ConstantProduct = &constProduct

		return result, nil
	}
	return result, errors.Errorf("invalid LiquidityPoolEntryBody type %v", b.Type)
}

func ConvertLiquidityPoolEntryConstantProduct(p xdr.LiquidityPoolEntryConstantProduct) (LiquidityPoolEntryConstantProduct, error) {
	var result LiquidityPoolEntryConstantProduct

	params, err := ConvertLiquidityPoolConstantProductParameters(p.Params)
	if err != nil {
		return result, err
	}

	result.Params = params
	result.ReserveA = int64(p.ReserveA)
	result.ReserveB = int64(p.ReserveB)
	result.TotalPoolShares = int64(p.TotalPoolShares)
	result.PoolSharesTrustLineCount = int64(p.PoolSharesTrustLineCount)

	return result, nil
}
