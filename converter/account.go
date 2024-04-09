package converter

import (
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

func ConvertAccountId(a xdr.AccountId) (AccountId, error) {
	var result AccountId
	accountId, err := a.GetAddress()
	if err != nil {
		return result, err
	}

	result.Address = accountId

	return result, nil
}

func ConvertAccountEntry(e xdr.AccountEntry) (AccountEntry, error) {
	var result AccountEntry

	accountId, err := ConvertAccountId(e.AccountId)
	if err != nil {
		return result, err
	}

	var inflationDest AccountId
	if e.InflationDest != nil {
		inflationDest, err = ConvertAccountId(*e.InflationDest)
		if err != nil {
			return result, err
		}
	}

	var signers []Signer
	for _, xdrSigner := range e.Signers {
		signer, err := ConvertSigner(xdrSigner)
		if err != nil {
			return result, err
		}

		signers = append(signers, signer)
	}

	ext := ConvertAccountEntryExt(e.Ext)

	result.AccountId = accountId
	result.Balance = int64(e.Balance)
	result.SeqNum = int64(e.SeqNum)
	result.NumSubEntries = uint32(e.NumSubEntries)
	result.InflationDest = &inflationDest
	result.Flags = uint32(e.Flags)
	result.HomeDomain = string(e.HomeDomain)
	result.Thresholds = e.Thresholds[:]
	result.Signers = signers
	result.Ext = ext

	return result, nil
}

func ConvertAccountEntryExt(e xdr.AccountEntryExt) AccountEntryExt {
	var v1 AccountEntryExtensionV1
	if e.V1 != nil {
		v1 = ConvertAccountEntryExtensionV1(*e.V1)
	}

	return AccountEntryExt{
		V:  e.V,
		V1: &v1,
	}
}

func ConvertAccountEntryExtensionV1(e xdr.AccountEntryExtensionV1) AccountEntryExtensionV1 {
	return AccountEntryExtensionV1{
		Liabilities: ConvertLiabilities(e.Liabilities),
		Ext:         ConvertAccountEntryExtensionV1Ext(e.Ext),
	}
}

func ConvertLiabilities(l xdr.Liabilities) Liabilities {
	return Liabilities{
		Buying:  int64(l.Buying),
		Selling: int64(l.Selling),
	}
}

func ConvertAccountEntryExtensionV1Ext(e xdr.AccountEntryExtensionV1Ext) AccountEntryExtensionV1Ext {
	var v2 AccountEntryExtensionV2
	if e.V2 != nil {
		v2, _ = ConvertAccountEntryExtensionV2(*e.V2)
	}

	return AccountEntryExtensionV1Ext{
		V:  e.V,
		V2: &v2,
	}
}

func ConvertAccountEntryExtensionV2(e xdr.AccountEntryExtensionV2) (AccountEntryExtensionV2, error) {
	var signerSponsoringIDs []AccountId
	if e.SignerSponsoringIDs != nil {
		for _, xdrSigner := range e.SignerSponsoringIDs {
			if xdrSigner != nil {
				signer, err := ConvertAccountId(*xdrSigner)
				if err != nil {
					return AccountEntryExtensionV2{}, err
				}

				signerSponsoringIDs = append(signerSponsoringIDs, signer)
			}

		}
	}

	ext := ConvertAccountEntryExtensionV2Ext(e.Ext)

	return AccountEntryExtensionV2{
		NumSponsored:        uint32(e.NumSponsored),
		NumSponsoring:       uint32(e.NumSponsoring),
		SignerSponsoringIDs: signerSponsoringIDs,
		Ext:                 ext,
	}, nil
}

func ConvertAccountEntryExtensionV2Ext(e xdr.AccountEntryExtensionV2Ext) AccountEntryExtensionV2Ext {
	var v3 AccountEntryExtensionV3
	if e.V3 != nil {
		v3 = ConvertAccountEntryExtensionV3(*e.V3)
	}

	return AccountEntryExtensionV2Ext{
		V:  e.V,
		V3: &v3,
	}
}

func ConvertAccountEntryExtensionV3(e xdr.AccountEntryExtensionV3) AccountEntryExtensionV3 {
	return AccountEntryExtensionV3{
		Ext:       ConvertExtensionPoint(e.Ext),
		SeqLedger: uint32(e.SeqLedger),
		SeqTime:   uint64(e.SeqTime),
	}
}

// TODO: testing
func ConvertSigner(s xdr.Signer) (Signer, error) {
	var result Signer
	signerKey, err := ConvertSignerKey(s.Key)
	if err != nil {
		return result, err
	}
	result.Key = signerKey
	result.Weight = uint32(s.Weight)

	return result, nil
}

// TODO: testing
func ConvertSignerKey(k xdr.SignerKey) (SignerKey, error) {
	var result SignerKey
	address, err := k.GetAddress()
	if err != nil {
		return result, err
	}
	result.Address = address

	return result, nil
}

// TODO: testing
func ConvertMuxedAccount(ma xdr.MuxedAccount) (MuxedAccount, error) {
	var result MuxedAccount

	address, err := ma.GetAddress()
	if err != nil {
		return result, err
	}
	result.Address = address

	return result, nil
}

// TODO :testing
func ConvertRevokeSponsorshipOpSigner(s xdr.RevokeSponsorshipOpSigner) (RevokeSponsorshipOpSigner, error) {
	var result RevokeSponsorshipOpSigner

	accountId, err := ConvertAccountId(s.AccountId)
	if err != nil {
		return result, err
	}

	signerKey, err := ConvertSignerKey(s.SignerKey)
	if err != nil {
		return result, err
	}

	result.AccountId = accountId
	result.SignerKey = signerKey

	return result, nil
}

func ConvertDecoratedSignature(s xdr.DecoratedSignature) DecoratedSignature {
	return DecoratedSignature{
		Hint:      s.Hint[:],
		Signature: s.Signature,
	}
}

func ConvertEd25519(inp xdr.Uint256) (string, error) {
	raw := make([]byte, 32)
	copy(raw, inp[:])
	return strkey.Encode(strkey.VersionByteAccountID, raw)
}
