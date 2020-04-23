package applepay

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"
)

// https://developer.apple.com/documentation/apple_pay_on_the_web/configuring_your_environment
// https://developer.apple.com/documentation/apple_pay_on_the_web
type Merchant struct {
	DefaultPaymentSessionRequest PaymentSessionRequest

	// tls.LoadX509KeyPair()
	IdentityCertificate tls.Certificate
	identityClientOnce  sync.Once
	identityClient      *http.Client
}

// https://developer.apple.com/documentation/apple_pay_on_the_web/apple_pay_js_api/requesting_an_apple_pay_payment_session
type PaymentSessionRequest struct {
	MerchantIdentifier string `json:"merchantIdentifier"`
	DisplayName        string `json:"displayName"`
	Initiative         string `json:"initiative"`
	InitiativeContext  string `json:"initiativeContext"`
}

// PaymentSession receives an opaque Apple Pay session object.
// https://developer.apple.com/documentation/apple_pay_on_the_web/apple_pay_js_api/requesting_an_apple_pay_payment_session
//
// validationURL:
// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaysession/1778021-onvalidatemerchant
func (m *Merchant) PaymentSession(validationURL string, req PaymentSessionRequest) (session []byte, err error) {
	if err := checkValidationURL(validationURL); err != nil {
		return nil, err
	}

	if req.MerchantIdentifier == "" {
		req.MerchantIdentifier = m.DefaultPaymentSessionRequest.MerchantIdentifier
	}
	if req.DisplayName == "" {
		req.DisplayName = m.DefaultPaymentSessionRequest.DisplayName
	}
	if req.Initiative == "" {
		req.Initiative = m.DefaultPaymentSessionRequest.Initiative
	}
	if req.InitiativeContext == "" {
		req.InitiativeContext = m.DefaultPaymentSessionRequest.InitiativeContext
	}

	cli := m.getIdentityClient()
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := cli.Post(validationURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// https://developer.apple.com/documentation/apple_pay_on_the_web/setting_up_your_server
func checkValidationURL(validationURL string) error {
	u, err := url.Parse(validationURL)
	if err != nil {
		return err
	}
	hostReg := regexp.MustCompile("^apple-pay-gateway(-.+)?.apple.com$")
	if !hostReg.MatchString(u.Host) {
		return errors.New("validationURL is not belongs to apple")
	}
	if u.Scheme != "https" {
		return errors.New("validationURL scheme is not https")
	}
	return nil
}

func (m *Merchant) getIdentityClient() *http.Client {
	m.identityClientOnce.Do(func() {
		m.identityClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{
						m.IdentityCertificate,
					},
				},
			},
			Timeout: 30 * time.Second,
		}
	})
	return m.identityClient
}
