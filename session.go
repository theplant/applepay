package applepay

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// https://developer.apple.com/documentation/apple_pay_on_the_web/configuring_your_environment
// https://developer.apple.com/documentation/apple_pay_on_the_web
type Merchant struct {
	// https://developer.apple.com/documentation/apple_pay_on_the_web/setting_up_your_server
	Sandbox                      bool
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
	if err := m.checkValidationURL(validationURL); err != nil {
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

var productionValidationURLs = map[string]bool{
	"apple-pay-gateway.apple.com":               true,
	"apple-pay-gateway-nc-pod1.apple.com":       true,
	"apple-pay-gateway-nc-pod2.apple.com":       true,
	"apple-pay-gateway-nc-pod3.apple.com":       true,
	"apple-pay-gateway-nc-pod4.apple.com":       true,
	"apple-pay-gateway-nc-pod5.apple.com":       true,
	"apple-pay-gateway-pr-pod1.apple.com":       true,
	"apple-pay-gateway-pr-pod2.apple.com":       true,
	"apple-pay-gateway-pr-pod3.apple.com":       true,
	"apple-pay-gateway-pr-pod4.apple.com":       true,
	"apple-pay-gateway-pr-pod5.apple.com":       true,
	"apple-pay-gateway-nc-pod1-dr.apple.com":    true,
	"apple-pay-gateway-nc-pod2-dr.apple.com":    true,
	"apple-pay-gateway-nc-pod3-dr.apple.com":    true,
	"apple-pay-gateway-nc-pod4-dr.apple.com":    true,
	"apple-pay-gateway-nc-pod5-dr.apple.com":    true,
	"apple-pay-gateway-pr-pod1-dr.apple.com":    true,
	"apple-pay-gateway-pr-pod2-dr.apple.com":    true,
	"apple-pay-gateway-pr-pod3-dr.apple.com":    true,
	"apple-pay-gateway-pr-pod4-dr.apple.com":    true,
	"apple-pay-gateway-pr-pod5-dr.apple.com":    true,
	"cn-apple-pay-gateway-sh-pod1.apple.com":    true,
	"cn-apple-pay-gateway-sh-pod1-dr.apple.com": true,
	"cn-apple-pay-gateway-sh-pod2.apple.com":    true,
	"cn-apple-pay-gateway-sh-pod2-dr.apple.com": true,
	"cn-apple-pay-gateway-sh-pod3.apple.com":    true,
	"cn-apple-pay-gateway-sh-pod3-dr.apple.com": true,
	"cn-apple-pay-gateway-tj-pod1.apple.com":    true,
	"cn-apple-pay-gateway-tj-pod1-dr.apple.com": true,
	"cn-apple-pay-gateway-tj-pod2.apple.com":    true,
	"cn-apple-pay-gateway-tj-pod2-dr.apple.com": true,
	"cn-apple-pay-gateway-tj-pod3.apple.com":    true,
	"cn-apple-pay-gateway-tj-pod3-dr.apple.com": true,
}

var sandboxValidationURLs = map[string]bool{
	"apple-pay-gateway-cert.apple.com":    true,
	"cn-apple-pay-gateway-cert.apple.com": true,
}

// https://developer.apple.com/documentation/apple_pay_on_the_web/setting_up_your_server
func (m *Merchant) checkValidationURL(validationURL string) error {
	u, err := url.Parse(validationURL)
	if err != nil {
		return err
	}
	if m.Sandbox {
		if !sandboxValidationURLs[u.Host] && !productionValidationURLs[u.Host] {
			return errors.New("validationURL is not exists in production/sandbox validation URLs")
		}
	} else {
		if !productionValidationURLs[u.Host] {
			return errors.New("validationURL is not exists in production validation URLs")
		}
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
