package applepay

import (
	"crypto/tls"
	"os"
	"testing"
)

var (
	identifier                      = os.Getenv("ApplePay_MerchantIdentifier")
	defaultDisplayName              = os.Getenv("ApplePay_MerchantDisplayName")
	identityCertificateCertFile     = os.Getenv("ApplePay_MerchantIdentityCertificateCertFile")
	identityCertificateKeyFile      = os.Getenv("ApplePay_MerchantIdentityCertificateKeyFile")
	paymentSessionInitiative        = os.Getenv("ApplePay_PaymentSessionInitiative")
	paymentSessionInitiativeContext = os.Getenv("ApplePay_PaymentSessionInitiativeContext")
)

func TestPaymentSession(t *testing.T) {
	identityCertificate, err := tls.LoadX509KeyPair(identityCertificateCertFile, identityCertificateKeyFile)
	if err != nil {
		t.Fatal(err)
	}

	merchant := &Merchant{
		DefaultPaymentSessionRequest: PaymentSessionRequest{
			MerchantIdentifier: identifier,
			DisplayName:        defaultDisplayName,
			Initiative:         paymentSessionInitiative,
			InitiativeContext:  paymentSessionInitiativeContext,
		},
		IdentityCertificate: identityCertificate,
	}

	t.Log(merchant.DefaultPaymentSessionRequest)

	session, err := merchant.PaymentSession(
		"https://apple-pay-gateway-cert.apple.com/paymentservices/startSession",
		PaymentSessionRequest{},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(session))
}
