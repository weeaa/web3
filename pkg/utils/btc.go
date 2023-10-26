package utils

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type BitcoinWallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *btcutil.AddressTaproot
}

func GenerateBtcTaprootWallet() (*BitcoinWallet, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	wif := btcutil.WIF{
		PrivKey: privateKey,
	}

	wifKey, err := btcutil.DecodeWIF(wif.String())
	if err != nil {
		return nil, err
	}

	taprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(wifKey.PrivKey.PubKey())), &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &BitcoinWallet{
		PrivateKey: wif.PrivKey.ToECDSA(),
		PublicKey:  taprootAddress,
	}, nil
}

func ConsolidateBtcFunds() {}

func DisperseBtcFunds() {}
