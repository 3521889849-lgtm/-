package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/url"
	"testing"

	"example_shop/common/config"
)

func TestVerify_RSA2WithAliPublicKey(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("marshal pub: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
	config.Cfg.AliPay.AliPublicKey = string(pubPEM)

	values := url.Values{}
	values.Set("app_id", "2021000118630000")
	values.Set("out_trade_no", "P-123")
	values.Set("trade_status", "TRADE_SUCCESS")
	values.Set("total_amount", "10.00")

	signData := buildSignData(values)
	sum := sha256.Sum256([]byte(signData))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, sum[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	values.Set("sign", base64.StdEncoding.EncodeToString(sig))
	values.Set("sign_type", "RSA2")

	if err := Verify(values); err != nil {
		t.Fatalf("verify: %v", err)
	}
}
