package converter

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertSorobanAuthorizationEntry(e xdr.SorobanAuthorizationEntry) (SorobanAuthorizationEntry, error) {
	var result SorobanAuthorizationEntry

	credentials, err := ConvertSorobanCredentials(e.Credentials)
	if err != nil {
		return result, err
	}

	rootInvocation, err := ConvertSorobanAuthorizedInvocation(e.RootInvocation)
	if err != nil {
		return result, err
	}

	result.Credentials = credentials
	result.RootInvocation = rootInvocation

	return result, nil
}

func ConvertSorobanCredentials(c xdr.SorobanCredentials) (SorobanCredentials, error) {
	var result SorobanCredentials
	switch c.Type {
	case xdr.SorobanCredentialsTypeSorobanCredentialsSourceAccount:
		// void
		return result, nil
	case xdr.SorobanCredentialsTypeSorobanCredentialsAddress:
		address, err := ConvertSorobanAddressCredentials(*c.Address)
		if err != nil {
			return result, err
		}

		result.Address = &address
		return result, nil
	}

	return result, errors.Errorf("Invalid ConvertSorobanCredentials type %v\n", c.Type)
}

func ConvertSorobanAddressCredentials(c xdr.SorobanAddressCredentials) (SorobanAddressCredentials, error) {
	var result SorobanAddressCredentials

	address, err := ConvertScAddress(c.Address)
	if err != nil {
		return result, err
	}

	signature, err := ConvertScVal(c.Signature)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Nonce = int64(c.Nonce)
	result.SignatureExpirationLedger = uint32(c.SignatureExpirationLedger)
	result.Signature = signature

	return result, nil
}

func ConvertSorobanAuthorizedInvocation(i xdr.SorobanAuthorizedInvocation) (SorobanAuthorizedInvocation, error) {
	var result SorobanAuthorizedInvocation
	function, err := ConvertSorobanAuthorizedFunction(i.Function)
	if err != nil {
		return result, err
	}
	result.Function = function

	var subs []SorobanAuthorizedInvocation
	for _, xdrSub := range i.SubInvocations {
		sub, err := ConvertSorobanAuthorizedInvocation(xdrSub)
		if err != nil {
			return result, err
		}

		subs = append(subs, sub)
	}
	result.SubInvocations = subs

	return result, nil
}

func ConvertSorobanAuthorizedFunction(f xdr.SorobanAuthorizedFunction) (SorobanAuthorizedFunction, error) {
	var result SorobanAuthorizedFunction
	switch f.Type {
	case xdr.SorobanAuthorizedFunctionTypeSorobanAuthorizedFunctionTypeContractFn:
		contractFn, err := ConvertInvokeContractArgs(*f.ContractFn)
		if err != nil {
			return result, err
		}
		result.ContractFn = &contractFn

		return result, nil
	case xdr.SorobanAuthorizedFunctionTypeSorobanAuthorizedFunctionTypeCreateContractHostFn:
		createContract, err := ConvertCreateContractArgs(*f.CreateContractHostFn)
		if err != nil {
			return result, err
		}
		result.CreateContractHostFn = &createContract

		return result, nil
	}

	return result, errors.Errorf("Invalid SorobanAuthorizedFunction type %v", f.Type)
}

func ConvertSorobanTransactionData(d xdr.SorobanTransactionData) (SorobanTransactionData, error) {
	var result SorobanTransactionData

	resources, err := ConvertSorobanResources(d.Resources)
	if err != nil {
		return result, err
	}

	result.Ext = ConvertExtensionPoint(d.Ext)
	result.Resources = resources
	result.ResourceFee = int64(d.ResourceFee)

	return result, nil
}

func ConvertSorobanResources(r xdr.SorobanResources) (SorobanResources, error) {
	var result SorobanResources

	footPrint, err := ConvertLedgerFootprint(r.Footprint)
	if err != nil {
		return result, err
	}

	result.Footprint = footPrint
	result.Instructions = uint32(r.Instructions)
	result.ReadBytes = uint32(r.ReadBytes)
	result.WriteBytes = uint32(r.WriteBytes)

	return result, nil
}

func ConvertHostFunction(f xdr.HostFunction) (HostFunction, error) {
	var result HostFunction
	switch f.Type {
	case xdr.HostFunctionTypeHostFunctionTypeInvokeContract:
		invokeContract, err := ConvertInvokeContractArgs(*f.InvokeContract)
		if err != nil {
			return result, err
		}
		result.InvokeContract = &invokeContract

		return result, nil
	case xdr.HostFunctionTypeHostFunctionTypeCreateContract:
		createContract, err := ConvertCreateContractArgs(*f.CreateContract)
		if err != nil {
			return result, err
		}
		result.CreateContract = &createContract

		return result, nil
	case xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm:
		wasm := *f.Wasm
		result.Wasm = &wasm

		return result, nil
	}

	return result, errors.Errorf("Invalid host function type %v", f.Type)
}

func ConvertInvokeContractArgs(a xdr.InvokeContractArgs) (InvokeContractArgs, error) {
	var result InvokeContractArgs

	contractAddress, err := ConvertScAddress(a.ContractAddress)
	if err != nil {
		return result, err
	}

	funcName := ScSymbol(a.FunctionName)

	var args []ScVal
	for _, xdrArg := range a.Args {
		arg, err := ConvertScVal(xdrArg)
		if err != nil {
			return result, err
		}

		args = append(args, arg)
	}

	result.ContractAddress = contractAddress
	result.FunctionName = funcName
	result.Args = args

	return result, nil
}

func ConvertCreateContractArgs(a xdr.CreateContractArgs) (CreateContractArgs, error) {
	var result CreateContractArgs

	contractIdPreimage, err := ConvertContractIdPreimage(a.ContractIdPreimage)
	if err != nil {
		return result, err
	}

	executable, err := ConvertContractExecutable(a.Executable)
	if err != nil {
		return result, err
	}

	result.ContractIdPreimage = contractIdPreimage
	result.Executable = executable

	return result, nil
}

func ConvertContractExecutable(e xdr.ContractExecutable) (ContractExecutable, error) {
	var result ContractExecutable
	switch e.Type {
	case xdr.ContractExecutableTypeContractExecutableWasm:
		wasmHash := (*e.WasmHash).HexString()
		result.WasmHash = &wasmHash
		return result, nil
	case xdr.ContractExecutableTypeContractExecutableStellarAsset:
		return result, nil
	}

	return result, errors.Errorf("Invalid contract executable type %v", e.Type)
}

func ConvertContractIdPreimage(p xdr.ContractIdPreimage) (ContractIdPreimage, error) {
	var result ContractIdPreimage

	switch p.Type {
	case xdr.ContractIdPreimageTypeContractIdPreimageFromAddress:
		fromAddress, err := ConvertContractIdPreimageFromAddress(*p.FromAddress)
		if err != nil {
			return result, err
		}
		result.FromAddress = &fromAddress

		return result, nil
	case xdr.ContractIdPreimageTypeContractIdPreimageFromAsset:
		fromAsset, err := ConvertAsset(*p.FromAsset)
		if err != nil {
			return result, err
		}
		result.FromAsset = &fromAsset

		return result, nil
	}
	return result, errors.Errorf("Invalid contract id preimage type %v", p.Type)
}

func ConvertContractIdPreimageFromAddress(p xdr.ContractIdPreimageFromAddress) (ContractIdPreimageFromAddress, error) {
	var result ContractIdPreimageFromAddress
	address, err := ConvertScAddress(p.Address)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Salt = p.Salt.String()

	return result, nil
}

func ConvertContractCodeEntry(e xdr.ContractCodeEntry) ContractCodeEntry {
	return ContractCodeEntry{
		Ext:  ConvertContractCodeEntryExt(e.Ext),
		Hash: e.Hash.HexString(),
		Code: e.Code,
	}
}

func ConvertContractCodeEntryExt(e xdr.ContractCodeEntryExt) ContractCodeEntryExt {
	var result ContractCodeEntryExt
	switch e.V {
	case 0:
		result.V = 0
	case 1:
		result.V = 1
		v1Xdr := e.MustV1()
		v1 := ConvertContractCodeEntryV1(v1Xdr)
		result.V1 = &v1
	}

	return result
}

func ConvertContractCodeEntryV1(e xdr.ContractCodeEntryV1) ContractCodeEntryV1 {
	return ContractCodeEntryV1{
		Ext:        ConvertExtensionPoint(e.Ext),
		CostInputs: ConvertContractCodeCostInputs(e.CostInputs),
	}
}

func ConvertContractCodeCostInputs(i xdr.ContractCodeCostInputs) ContractCodeCostInputs {
	return ContractCodeCostInputs{
		Ext:               ConvertExtensionPoint(i.Ext),
		NInstructions:     uint32(i.NInstructions),
		NFunctions:        uint32(i.NFunctions),
		NGlobals:          uint32(i.NGlobals),
		NTableEntries:     uint32(i.NTableEntries),
		NTypes:            uint32(i.NTypes),
		NDataSegments:     uint32(i.NDataSegments),
		NElemSegments:     uint32(i.NElemSegments),
		NImports:          uint32(i.NImports),
		NExports:          uint32(i.NExports),
		NDataSegmentBytes: uint32(i.NDataSegmentBytes),
	}
}

func ConvertContractDataEntry(e xdr.ContractDataEntry) (ContractDataEntry, error) {
	var result ContractDataEntry

	ext := ConvertExtensionPoint(e.Ext)

	contract, err := ConvertScAddress(e.Contract)
	if err != nil {
		return result, err
	}

	key, err := ConvertScVal(e.Key)
	if err != nil {
		return result, err
	}

	val, err := ConvertScVal(e.Val)
	if err != nil {
		return result, err
	}

	result.Ext = ext
	result.Contract = contract
	result.Key = key
	result.Durability = int32(e.Durability)
	result.Val = val

	return result, nil
}

func ConvertConfigSettingContractComputeV0(c xdr.ConfigSettingContractComputeV0) ConfigSettingContractComputeV0 {
	return ConfigSettingContractComputeV0{
		LedgerMaxInstructions:           int64(c.LedgerMaxInstructions),
		TxMaxInstructions:               int64(c.TxMaxInstructions),
		FeeRatePerInstructionsIncrement: int64(c.FeeRatePerInstructionsIncrement),
		TxMemoryLimit:                   uint32(c.TxMemoryLimit),
	}
}

func ConvertConfigSettingContractLedgerCostV0(c xdr.ConfigSettingContractLedgerCostV0) ConfigSettingContractLedgerCostV0 {
	return ConfigSettingContractLedgerCostV0{
		LedgerMaxReadLedgerEntries:     uint32(c.LedgerMaxReadLedgerEntries),
		LedgerMaxReadBytes:             uint32(c.LedgerMaxReadBytes),
		LedgerMaxWriteLedgerEntries:    uint32(c.LedgerMaxWriteLedgerEntries),
		LedgerMaxWriteBytes:            uint32(c.LedgerMaxWriteBytes),
		TxMaxReadLedgerEntries:         uint32(c.TxMaxReadLedgerEntries),
		TxMaxReadBytes:                 uint32(c.TxMaxReadBytes),
		TxMaxWriteLedgerEntries:        uint32(c.TxMaxWriteLedgerEntries),
		TxMaxWriteBytes:                uint32(c.TxMaxWriteBytes),
		FeeReadLedgerEntry:             int64(c.FeeReadLedgerEntry),
		FeeWriteLedgerEntry:            int64(c.FeeWriteLedgerEntry),
		FeeRead1Kb:                     int64(c.FeeRead1Kb),
		BucketListTargetSizeBytes:      int64(c.BucketListTargetSizeBytes),
		WriteFee1KbBucketListLow:       int64(c.WriteFee1KbBucketListLow),
		WriteFee1KbBucketListHigh:      int64(c.WriteFee1KbBucketListHigh),
		BucketListWriteFeeGrowthFactor: uint32(c.BucketListWriteFeeGrowthFactor),
	}
}

func ConvertConfigSettingContractHistoricalDataV0(c xdr.ConfigSettingContractHistoricalDataV0) ConfigSettingContractHistoricalDataV0 {
	return ConfigSettingContractHistoricalDataV0{FeeHistorical1Kb: int64(c.FeeHistorical1Kb)}
}

func ConvertConfigSettingContractEventsV0(c xdr.ConfigSettingContractEventsV0) ConfigSettingContractEventsV0 {
	return ConfigSettingContractEventsV0{
		TxMaxContractEventsSizeBytes: uint32(c.TxMaxContractEventsSizeBytes),
		FeeContractEvents1Kb:         int64(c.FeeContractEvents1Kb),
	}
}

func ConvertConfigSettingContractBandwidthV0(c xdr.ConfigSettingContractBandwidthV0) ConfigSettingContractBandwidthV0 {
	return ConfigSettingContractBandwidthV0{
		LedgerMaxTxsSizeBytes: uint32(c.LedgerMaxTxsSizeBytes),
		TxMaxSizeBytes:        uint32(c.TxMaxSizeBytes),
		FeeTxSize1Kb:          int64(c.FeeTxSize1Kb),
	}
}

func ConvertContractCostParams(c xdr.ContractCostParams) ContractCostParams {
	var result ContractCostParams
	for _, xdrEntry := range c {
		entry := ConvertContractCostParamEntry(xdrEntry)
		result = append(result, entry)
	}
	return result
}

func ConvertContractCostParamEntry(c xdr.ContractCostParamEntry) ContractCostParamEntry {
	return ContractCostParamEntry{
		Ext:        ConvertExtensionPoint(c.Ext),
		ConstTerm:  int64(c.ConstTerm),
		LinearTerm: int64(c.LinearTerm),
	}
}

func ConvertConfigSettingContractExecutionLanesV0(c xdr.ConfigSettingContractExecutionLanesV0) ConfigSettingContractExecutionLanesV0 {
	return ConfigSettingContractExecutionLanesV0{LedgerMaxTxCount: uint32(c.LedgerMaxTxCount)}
}

func ConvertEvictionIterator(i xdr.EvictionIterator) EvictionIterator {
	return EvictionIterator{
		BucketListLevel:  uint32(i.BucketListLevel),
		IsCurrBucket:     i.IsCurrBucket,
		BucketFileOffset: uint64(i.BucketFileOffset),
	}
}

const (
	XDR_BOOL       = "bool"
	XDR_U32        = "u32"
	XDR_I32        = "i32"
	XDR_U64        = "u64"
	XDR_I64        = "i64"
	XDR_TIME_POINT = "time"
	XDR_DURATION   = "duration"
	XDR_U128       = "u128"
	XDR_I128       = "i128"
	XDR_U256       = "u256"
	XDR_I256       = "i256"
	XDR_BYTES      = "bytes"
	XDR_STRING     = "string"
	XDR_SYM        = "sym"
	XDR_NONCE      = "nonce"
	XDR_VEC        = "vec"
	XDR_MAP        = "map"
	XDR_ADDRESS    = "address"
)

func ConvertToData(keyType string, keyValue string) (xdr.ScVal, error) {
	switch keyType {
	case XDR_BOOL:
		return convertToDataBool(keyValue)
	case XDR_U32:
		return convertToDataUint32(keyValue)
	case XDR_I32:
		return convertToDataInt32(keyValue)
	case XDR_U64:
		return convertToDataUint64(keyValue)
	case XDR_I64:
		return convertToDataInt64(keyValue)
	case XDR_TIME_POINT:
		return convertToDataTimePoint(keyValue)
	case XDR_DURATION:
		return convertToDataDuration(keyValue)
	case XDR_U128:
		return convertToDataUInt128Parts(keyValue)
	case XDR_I128:
		return convertToDataInt128Parts(keyValue)
	case XDR_U256:
		return convertToDataUInt256Parts(keyValue)
	case XDR_I256:
		return convertToDataInt256Parts(keyValue)
	case XDR_BYTES:
		return convertToDataScBytes(keyValue)
	case XDR_STRING:
		return convertToDataScString(keyValue)
	case XDR_SYM:
		return convertToDataScSymbol(keyValue)
	case XDR_NONCE:
		return convertToDataScNonceKey(keyValue)
	case XDR_ADDRESS:
		return convertToDataScAddress(keyValue)
	case XDR_VEC:
		return convertToDataScVec(keyValue)
	default:
		return xdr.ScVal{}, errors.New("not found type")
	}
}

func convertToDataBool(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseBool(value)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvBool, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUint32(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvU32, xdr.Uint32(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt32(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvI32, xdr.Int32(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUint64(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt64(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvI64, xdr.Int64(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataTimePoint(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvTimepoint, xdr.TimePoint(xdr.Uint64(data)))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataDuration(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvDuration, xdr.Duration(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUInt128Parts(value string) (xdr.ScVal, error) {
	data, err := parseU128String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvU128, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt128Parts(value string) (xdr.ScVal, error) {
	data, err := parseI128String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvI128, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataUInt256Parts(value string) (xdr.ScVal, error) {
	data, err := parseU256String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvU256, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataInt256Parts(value string) (xdr.ScVal, error) {
	data, err := parseI256String(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvI256, data)
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScBytes(value string) (xdr.ScVal, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes(data))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScString(value string) (xdr.ScVal, error) {
	res, err := xdr.NewScVal(xdr.ScValTypeScvString, xdr.ScString(value))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScSymbol(value string) (xdr.ScVal, error) {
	res, err := xdr.NewScVal(xdr.ScValTypeScvSymbol, xdr.ScSymbol(value))
	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScNonceKey(value string) (xdr.ScVal, error) {
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return xdr.ScVal{}, err
	}

	res, err := xdr.NewScVal(xdr.ScValTypeScvLedgerKeyNonce, xdr.ScNonceKey{
		Nonce: xdr.Int64(data),
	})

	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScAddress(value string) (xdr.ScVal, error) {
	aid := xdr.MustAddress(value)
	res, err := xdr.NewScVal(xdr.ScValTypeScvAddress, xdr.ScAddress{
		AccountId: &aid,
	})

	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func convertToDataScVec(value string) (xdr.ScVal, error) {
	keys := strings.Split(value, ",")
	scVec := xdr.ScVec{}
	for _, key := range keys {
		keySplit := strings.Split(key, "@")
		if len(keySplit) != 2 {
			return xdr.ScVal{}, fmt.Errorf("not enough len")
		}

		item, err := ConvertToData(keySplit[0], keySplit[1])
		if err != nil {
			return xdr.ScVal{}, err
		}
		scVec = append(scVec, item)
	}
	res, err := xdr.NewScVal(xdr.ScValTypeScvVec, &scVec)

	if err != nil {
		return xdr.ScVal{}, err
	}

	return res, err
}

func parseU128String(s string) (xdr.UInt128Parts, error) {
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10)
	if !ok {
		return xdr.UInt128Parts{}, fmt.Errorf("invalid number format")
	}

	// Mask for lower 64 bits
	lowMask := new(big.Int).SetUint64(^uint64(0))

	// Extract the lower 64 bits
	low := new(big.Int).And(bigInt, lowMask).Uint64()

	// Extract the higher 64 bits by shifting right 64 bits
	high := new(big.Int).Rsh(bigInt, 64).Uint64()

	return xdr.UInt128Parts{
		Hi: xdr.Uint64(high),
		Lo: xdr.Uint64(low),
	}, nil
}

func parseI128String(s string) (xdr.Int128Parts, error) {
	// Parse the string as a big integer
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10) // Assuming the string is base 10
	if !ok {
		return xdr.Int128Parts{}, fmt.Errorf("invalid number format")
	}

	// Handle negative numbers for signed 128-bit integers
	negative := bigInt.Sign() < 0
	if negative {
		bigInt = bigInt.Abs(bigInt) // Convert to positive for bitwise operations
	}

	// Mask for lower 64 bits
	lowMask := new(big.Int).SetUint64(^uint64(0))

	// Extract the lower 64 bits
	low := new(big.Int).And(bigInt, lowMask).Uint64()

	// Extract the higher 64 bits and cast to int64 for signed interpretation
	high := new(big.Int).Rsh(bigInt, 64).Int64()
	if negative {
		high = -high
	}

	return xdr.Int128Parts{
		Hi: xdr.Int64(high),
		Lo: xdr.Uint64(low),
	}, nil
}

func parseU256String(s string) (xdr.UInt256Parts, error) {
	// Parse the string as a big integer
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10) // Assuming the string is base 10
	if !ok {
		return xdr.UInt256Parts{}, fmt.Errorf("invalid number format")
	}

	// Mask for 64 bits
	mask64 := new(big.Int).SetUint64(^uint64(0))

	// Extract the four 64-bit parts
	lowLow := new(big.Int).And(bigInt, mask64).Uint64()
	lowHigh := new(big.Int).Rsh(bigInt, 64).And(bigInt, mask64).Uint64()
	highLow := new(big.Int).Rsh(bigInt, 128).And(bigInt, mask64).Uint64()
	highHigh := new(big.Int).Rsh(bigInt, 192).Uint64()

	return xdr.UInt256Parts{
		HiHi: xdr.Uint64(highHigh),
		HiLo: xdr.Uint64(highLow),
		LoHi: xdr.Uint64(lowHigh),
		LoLo: xdr.Uint64(lowLow),
	}, nil
}

func parseI256String(s string) (xdr.Int256Parts, error) {
	// Parse the string as a big integer
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(s, 10) // Assuming the string is base 10
	if !ok {
		return xdr.Int256Parts{}, fmt.Errorf("invalid number format")
	}

	// Handle negative numbers
	negative := bigInt.Sign() < 0
	if negative {
		bigInt = bigInt.Abs(bigInt)
	}

	// Mask for 64 bits
	mask64 := new(big.Int).SetUint64(^uint64(0))

	// Extract the four 64-bit parts
	lowLow := new(big.Int).And(bigInt, mask64).Uint64()
	lowHigh := new(big.Int).Rsh(bigInt, 64).And(bigInt, mask64).Uint64()
	highLow := new(big.Int).Rsh(bigInt, 128).And(bigInt, mask64).Uint64()
	highHigh := new(big.Int).Rsh(bigInt, 192).Int64()

	if negative {
		highHigh = -highHigh
	}

	return xdr.Int256Parts{
		HiHi: xdr.Int64(highHigh),
		HiLo: xdr.Uint64(highLow),
		LoHi: xdr.Uint64(lowHigh),
		LoLo: xdr.Uint64(lowLow),
	}, nil
}
