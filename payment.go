package applepay

import (
	"encoding/json"
)

// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaysession/1778020-onpaymentauthorized
type ApplePayPayment struct {
	Token           ApplePayPaymentToken   `json:"token"`
	BillingContact  ApplePayPaymentContact `json:"billingContact"`
	ShippingContact ApplePayPaymentContact `json:"shippingContact"`
}

// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaypaymenttoken
type ApplePayPaymentToken struct {
	PaymentMethod         ApplePayPaymentMethod `json:"paymentMethod"`
	TransactionIdentifier string                `json:"transactionIdentifier"`
	PaymentData           json.RawMessage       `json:"paymentData"`
}

// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaypaymentmethod
type ApplePayPaymentMethod struct {
	DisplayName string                    `json:"displayName"`
	Network     string                    `json:"network"`
	Type        ApplePayPaymentMethodType `json:"type"`
	PaymentPass ApplePayPaymentPass       `json:"paymentPass"`
}

// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaypaymentmethodtype
type ApplePayPaymentMethodType string

const (
	PaymentMethodTypeDebit   ApplePayPaymentMethodType = "debit"
	PaymentMethodTypeCredit  ApplePayPaymentMethodType = "credit"
	PaymentMethodTypePrepaid ApplePayPaymentMethodType = "prepaid"
	PaymentMethodTypeStore   ApplePayPaymentMethodType = "store"
)

// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaypaymentpass
type ApplePayPaymentPass struct {
	PrimaryAccountIdentifier   string                             `json"primaryAccountIdentifier"`
	PrimaryAccountNumberSuffix string                             `json"primaryAccountNumberSuffix"`
	DeviceAccountIdentifier    string                             `json"deviceAccountIdentifier"`
	DeviceAccountNumberSuffix  string                             `json"deviceAccountNumberSuffix"`
	ActivationState            ApplePayPaymentPassActivationState `json:"activationState"`
}

// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaypaymentpassactivationstate
type ApplePayPaymentPassActivationState string

const (
	PaymentPassActivationStateactivated          ApplePayPaymentPassActivationState = "activated"
	PaymentPassActivationStaterequiresActivation ApplePayPaymentPassActivationState = "requiresActivation"
	PaymentPassActivationStateactivating         ApplePayPaymentPassActivationState = "activating"
	PaymentPassActivationStatesuspended          ApplePayPaymentPassActivationState = "suspended"
	PaymentPassActivationStatedeactivated        ApplePayPaymentPassActivationState = "deactivated"
)

// https://developer.apple.com/documentation/apple_pay_on_the_web/applepaypaymentcontact
type ApplePayPaymentContact struct {
	PhoneNumber           string   `json:"phoneNumber"`
	EmailAddress          string   `json:"emailAddress"`
	GivenName             string   `json:"givenName"`
	FamilyName            string   `json:"familyName"`
	PhoneticGivenName     string   `json:"phoneticGivenName"`
	PhoneticFamilyName    string   `json:"phoneticFamilyName"`
	AddressLines          []string `json:"addressLines"`
	SubLocality           string   `json:"subLocality"`
	Locality              string   `json:"locality"`
	PostalCode            string   `json:"postalCode"`
	SubAdministrativeArea string   `json:"subAdministrativeArea"`
	AdministrativeArea    string   `json:"administrativeArea"`
	Country               string   `json:"country"`
	CountryCode           string   `json:"countryCode"`
}
