package bitcoin

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
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

func DisperseFunds(privateKey string, addresses []btcutil.Address, amount int64, client *rpcclient.Client) {
	utxos, err := client.ListUnspentMinMaxAddresses(0, math.MaxInt32, addresses)
	if err != nil {

	}

	for _, address := range addresses {
		var destinationScript []byte

		tx := wire.NewMsgTx(wire.TxVersion)
		destinationScript, err = txscript.PayToAddrScript(address)
		if err != nil {

		}

		txOut := wire.NewTxOut(amount, destinationScript)
		tx.AddTxOut(txOut)

		for i, utxo := range utxos {

			sig, err := txscript.SignTxOutput(&chaincfg.MainNetParams, tx, i, []byte(utxo.ScriptPubKey), txscript.SigHashAll, privateKey, txscript.KeyClosure(nil), nil)
			if err != nil {
				return err
			}

			tx.TxIn[i].SignatureScript = sig
		}

	}

	inputTxHash, err := wire.NewAlertFromPayload("input_transaction_hash")
	if err != nil {
		log.Fatalf("Error parsing input transaction hash: %v", err)
	}
	prevOutPoint := wire.NewOutPoint(inputTxHash, 0)
	txIn := wire.NewTxIn(prevOutPoint, nil, nil)
	tx.AddTxIn(txIn)
	// Add output(s) to the transaction, including the destination address.
	destinationAddress := "destination_bitcoin_address"
	destinationScript, err := txscript.PayToAddrScript(destinationAddress, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Error creating destination script: %v", err)
	}
	txOut := wire.NewTxOut(20000000, destinationScript) // 0.2 BTC in satoshis (20000000)
	tx.AddTxOut(txOut)

	// Sign the transaction inputs.
	// You will need to sign the input using the private key associated with the UTXO.
	// For simplicity, we omit the signing part here.
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatalf("Error creating a private key: %v", err)
	}

	// Sign the transaction input.
	// Replace 'tx' and 'inputIndex' with your transaction and the index of the input to sign.
	sig, err := txscript.SignTxOutput(&chaincfg.MainNetParams, tx, inputIndex, destinationScript, txscript.SigHashAll, privateKey, txscript.KeyClosure(nil), nil, nil)
	if err != nil {
		log.Fatalf("Error signing the transaction: %v", err)
	}

	tx.TxIn[inputIndex].SignatureScript = sig
	// Send the signed transaction to the Bitcoin network.
	txid, err := client.SendRawTransaction(tx, false)
	if err != nil {
		log.Fatalf("Error broadcasting the transaction: %v", err)
	}

}

func ConsolidateFunds(privateKeys []string, address common.Address, amount *big.Int, client *ethclient.Client) {
	for _, privateKey := range privateKeys {
		_ = privateKey
	}
}

func ConsolidateOrdinals(privateKeys []string, address btcutil.AddressTaproot, tokenID string, client *rpcclient.Client) {

}

func ConsolidateBRC20(privateKeys []string, address btcutil.AddressTaproot, BRC20TokenID string, client *rpcclient.Client) {
}
