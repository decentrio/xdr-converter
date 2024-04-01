package converter

import (
	"encoding/json"

	"github.com/stellar/go/xdr"
)

func MarshalJSONEnvelopeXdr(inp []byte) ([]byte, error) {
	var xdrTxEnvelope xdr.TransactionEnvelope

	err := xdrTxEnvelope.UnmarshalBinary(inp)
	if err != nil {
		return nil, err
	}

	envelope, err := ConvertTransactionEnvelope(xdrTxEnvelope)
	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func MarshalJSONResultXdr(inp []byte) ([]byte, error) {
	var xdrTxResultPair xdr.TransactionResultPair

	err := xdrTxResultPair.UnmarshalBinary(inp)
	if err != nil {
		return nil, err
	}

	resultPair, err := ConvertTransactionResultPair(xdrTxResultPair)
	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(resultPair)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func MarshalJSONResultMetaXdr(inp []byte) ([]byte, error) {
	var xdrTxResultMeta xdr.TransactionResultMeta

	err := xdrTxResultMeta.UnmarshalBinary(inp)
	if err != nil {
		return nil, err
	}

	resultMeta, err := ConvertTransactionResultMeta(xdrTxResultMeta)
	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(resultMeta)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func MarshalJSONContractEventXdr(inp []byte) ([]byte, error) {
	var xdrContractEvent xdr.ContractEvent

	err := xdrContractEvent.UnmarshalBinary(inp)
	if err != nil {
		return nil, err
	}

	event, err := ConvertContractEvent(xdrContractEvent)
	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
