package bitcoin

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/crypto"
	"math"
)

type Wallet struct {
	PrivateKey     *ecdsa.PrivateKey
	TaprootAddress *btcutil.AddressTaproot
	PublicKey      *btcutil.AddressPubKey
}

type Transaction struct {
	ChainHash *chainhash.Hash
}

type Client struct {
	Client *rpcclient.Client
}

/*
func NewClienct() {
	connCfg := &rpcclient.ConnConfig{
		Host:         "127.0.0.1:8332", // The Bitcoin Core server host and port.
		User:         "yourusername",   // Your RPC username.
		Pass:         "yourpassword",   // Your RPC password.
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Error connecting to the Bitcoin network: %v", err)
	}
}
*/

func InitBtcWallet(privateStrKey string) (*Wallet, error) {
	privateKey, err := crypto.HexToECDSA(privateStrKey)
	if err != nil {
		return nil, err
	}

	wifKey, err := btcutil.DecodeWIF(privateStrKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := btcutil.NewAddressPubKey(wifKey.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	taprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(wifKey.PrivKey.PubKey())), &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &Wallet{PrivateKey: privateKey, PublicKey: publicKey, TaprootAddress: taprootAddress}, nil
}

func GenerateBtcTaprootWallet() (*Wallet, error) {
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

	publicKey, err := btcutil.NewAddressPubKey(wifKey.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		PrivateKey:     wif.PrivKey.ToECDSA(),
		TaprootAddress: taprootAddress,
		PublicKey:      publicKey,
	}, nil
}

func DisperseFunds(privateKey string, addresses []btcutil.Address, amount int64, client *rpcclient.Client) ([]*chainhash.Hash, error) {
	wallet, err := InitBtcWallet(privateKey)
	if err != nil {
		return nil, err
	}

	utxos, err := client.ListUnspentMinMaxAddresses(0, math.MaxInt32, []btcutil.Address{wallet.PublicKey})
	if err != nil {
		return nil, err
	}

	txns := make([]*chainhash.Hash, len(addresses))

	for _, address := range addresses {
		var tx = wire.NewMsgTx(wire.TxVersion)
		var totalAmount int64
		var pkScript []byte
		var txHash *chainhash.Hash
		var wif *btcutil.WIF
		var sigScript []byte

		for _, unspent := range utxos {
			txIn := wire.NewTxIn(&wire.OutPoint{Hash: chainhash.HashH([]byte(unspent.TxID)), Index: unspent.Vout}, nil, nil)
			tx.AddTxIn(txIn)
			totalAmount += int64(unspent.Amount)
		}

		pkScript, err = txscript.PayToAddrScript(address)
		if err != nil {
			return txns, err
		}

		tx.AddTxOut(wire.NewTxOut(amount, pkScript))

		if totalAmount > amount {
			change := totalAmount - amount
			pkScript, err = txscript.PayToAddrScript(wallet.PublicKey)
			if err != nil {
				return txns, err
			}
			tx.AddTxOut(wire.NewTxOut(change, pkScript))
		}

		wif, err = btcutil.DecodeWIF(fmt.Sprint(wallet.PrivateKey))
		if err != nil {
			return txns, err
		}

		for i, txIn := range tx.TxIn {
			sigScript, err = txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, wif.PrivKey, true)
			if err != nil {
				return txns, err
			}
			txIn.SignatureScript = sigScript
		}

		txHash, err = client.SendRawTransaction(tx, false)
		if err != nil {
			return txns, err
		}

		txns = append(txns, txHash)
	}

	return txns, nil
}

func ConsolidateFunds(privateKeys []string, address btcutil.Address, amount int64, client *rpcclient.Client) ([]*chainhash.Hash, error) {
	txns := make([]*chainhash.Hash, len(privateKeys))

	for _, privateKey := range privateKeys {
		var tx = wire.NewMsgTx(wire.TxVersion)
		var totalAmount int64
		var pkScript []byte
		var txHash *chainhash.Hash
		var wif *btcutil.WIF
		var sigScript []byte

		wallet, err := InitBtcWallet(privateKey)
		if err != nil {
			return nil, err
		}

		utxos, err := client.ListUnspentMinMaxAddresses(0, math.MaxInt32, []btcutil.Address{wallet.PublicKey})
		if err != nil {
			return nil, err
		}

		for _, unspent := range utxos {
			txIn := wire.NewTxIn(&wire.OutPoint{Hash: chainhash.HashH([]byte(unspent.TxID)), Index: unspent.Vout}, nil, nil)
			tx.AddTxIn(txIn)
			totalAmount += int64(unspent.Amount)
		}

		pkScript, err = txscript.PayToAddrScript(address)
		if err != nil {
			return txns, err
		}

		tx.AddTxOut(wire.NewTxOut(amount, pkScript))

		if totalAmount > amount {
			change := totalAmount - amount
			pkScript, err = txscript.PayToAddrScript(wallet.PublicKey)
			if err != nil {
				return txns, err
			}
			tx.AddTxOut(wire.NewTxOut(change, pkScript))
		}

		wif, err = btcutil.DecodeWIF(fmt.Sprint(wallet.PrivateKey))
		if err != nil {
			return txns, err
		}

		for i, txIn := range tx.TxIn {
			sigScript, err = txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, wif.PrivKey, true)
			if err != nil {
				return txns, err
			}
			txIn.SignatureScript = sigScript
		}

		txHash, err = client.SendRawTransaction(tx, false)
		if err != nil {
			return txns, err
		}

		txns = append(txns, txHash)
	}

	return txns, nil
}

func ConsolidateOrdinals(privateKeys []string, address btcutil.AddressTaproot, tokenID string, client *rpcclient.Client) {

}

func ConsolidateBRC20(privateKeys []string, address btcutil.AddressTaproot, BRC20TokenID string, client *rpcclient.Client) {

}
