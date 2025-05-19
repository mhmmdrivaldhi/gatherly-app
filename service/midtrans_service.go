package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gatherly-app/models/dto"

	"github.com/go-resty/resty/v2"
)

type MidtransService interface {
	Pay(payload dto.MidtransSnapReq) (dto.MidtransSnapResp, error)
	CancelTransaction(orderID string) error
}

type midtransService struct {
	client    *resty.Client
	serverKey string
	urlPay    string
	urlResp   string
	url       string
}

func NewMidtransService(client *resty.Client, serverKey string) MidtransService {
	return &midtransService{
		client:    client,
		serverKey: serverKey,
		urlPay:    "https://app.sandbox.midtrans.com/snap/v1/transactions",
		urlResp:   "https://app.sandbox.midtrans.com/snap/v2/vtweb",
		url:       "https://api.sandbox.midtrans.com/v2",
	}
}

func (m *midtransService) Pay(payload dto.MidtransSnapReq) (dto.MidtransSnapResp, error) {
	encodedKey := base64.StdEncoding.EncodeToString([]byte(m.serverKey))

	resp, err := m.client.R().
		SetHeader("Authorization", "Basic "+encodedKey).
		SetBody(payload).
		Post(m.urlPay)

	if err != nil {
		return dto.MidtransSnapResp{}, err
	}

	var snapResp dto.MidtransSnapResp
	err = json.Unmarshal(resp.Body(), &snapResp)
	if err != nil {
		return dto.MidtransSnapResp{}, err
	}

	snapResp.RedirectURL = fmt.Sprintf("%s/%s", m.urlResp, snapResp.Token)

	return snapResp, nil
}

func (m *midtransService) CancelTransaction(orderID string) error {
	encodedKey := base64.StdEncoding.EncodeToString([]byte(m.serverKey))

	url := fmt.Sprintf("%s/%s/cancel", m.url, orderID)
	_, err := m.client.R().
		SetHeader("Authorization", "Basic "+encodedKey).
		Post(url)

	if err != nil {
		return err
	}

	return nil
}
