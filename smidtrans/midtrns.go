package smidtrans

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/veritrans/go-midtrans"
	"time"
)

type Midtrans struct {
	SK  string
	CK  string
	ENV Env
}

type Env string

const Sandbox Env = "SANDBOX"
const Production Env = "PRODUCTION"

type CoreGatewayMidtrans struct {
	Core midtrans.CoreGateway
}

type IBankVAPayments interface {
	ChargeReqPermataVirtualAccount(
		itemID, itemName, trxID string,
		custName, custPhone string,
		grossAmt int64,
	) (*midtrans.ChargeReq, time.Time, time.Time)
	ChargeReqMandiriBill(
		itemID, itemName, trxID string,
		custName, custPhone string,
		grossAmt int64,
	) (*midtrans.ChargeReq, time.Time, time.Time)
	ChargeReqBNIVirtualAccount(
		itemID, itemName, trxID string,
		custName, custPhone string,
		grossAmt int64,
	) (*midtrans.ChargeReq, time.Time, time.Time)
	ChargeReqBCAVirtualAccount(
		itemID, itemName, trxID string,
		custName, custPhone string,
		grossAmt int64,
	) (*midtrans.ChargeReq, time.Time, time.Time)
	ChargeReqBRIVirtualAccount(
		itemID, itemName, trxID string,
		custName, custPhone string,
		grossAmt int64,
	) (*midtrans.ChargeReq, time.Time, time.Time)
	RequestCharge(chargeReq *midtrans.ChargeReq) (*midtrans.Response, error)
}

func (m *Midtrans) InitializeMidtransClient() *CoreGatewayMidtrans {
	midclient := midtrans.NewClient()
	midclient.ServerKey = m.SK
	midclient.ClientKey = m.CK
	switch m.ENV {
	case Sandbox:
		midclient.APIEnvType = midtrans.Sandbox
	case Production:
		midclient.APIEnvType = midtrans.Production
	default:
		panic(errors.New("invalid environment type"))
	}
	return &CoreGatewayMidtrans{
		Core: midtrans.CoreGateway{
			Client: midclient,
		},
	}
}

type SignatureVerify struct {
	OrderID     string
	StatusCode  string
	GrossAmount string
	ServerKey   string
}

func (m *CoreGatewayMidtrans) VerifySignature(signature string, verify SignatureVerify) (bool, error) {
	formulas := verify.OrderID + verify.StatusCode + verify.GrossAmount + verify.ServerKey
	newSha512 := sha512.New()
	newSha512.Write([]byte(formulas))
	formatted := fmt.Sprintf("%x", newSha512.Sum(nil))
	if signature == formatted {
		return true, nil
	}
	return false, errors.New("not authentic")
}

type MidtransTransaction struct {
	Status     string
	Gross      string
	ExpireAt   time.Time
	TrxCode    string
	TrxType    int
	Message    string
	Bank       string
	VaNumber   string
	BillerCode string
	BillKey    string
}

func DetermineConsumableResponse(res *midtrans.Response, expire time.Time) (*MidtransTransaction, error) {
	if res.PermataVaNumber != "" {
		return &MidtransTransaction{
			ExpireAt:   expire,
			TrxType:    1,
			Status:     res.StatusCode,
			Gross:      res.GrossAmount,
			Message:    res.StatusMessage,
			Bank:       "permata",
			VaNumber:   res.PermataVaNumber,
			BillerCode: res.BillerCode,
			BillKey:    res.BillKey,
		}, nil
	} else {
		if res.BillKey != "" && res.BillerCode != "" {
			return &MidtransTransaction{
				ExpireAt:   expire,
				TrxType:    1,
				Status:     res.StatusCode,
				Gross:      res.GrossAmount,
				Message:    res.StatusMessage,
				Bank:       "mandiri",
				VaNumber:   "",
				BillerCode: res.BillerCode,
				BillKey:    res.BillKey,
			}, nil
		} else {
			if len(res.VANumbers) != 0 {
				return &MidtransTransaction{
					ExpireAt:   expire,
					TrxType:    1,
					Status:     res.StatusCode,
					Gross:      res.GrossAmount,
					Message:    res.StatusMessage,
					Bank:       res.VANumbers[0].Bank,
					VaNumber:   res.VANumbers[0].VANumber,
					BillerCode: res.BillerCode,
					BillKey:    res.BillKey,
				}, nil
			} else {
				return nil, errors.New("error while doing payment")
			}
		}
	}
}

func (c *CoreGatewayMidtrans) ChargeReqMandiriBill(
	itemID, itemName, trxID string,
	custName, custPhone string,
	grossAmt int64,
) (*midtrans.ChargeReq, time.Time, time.Time) {
	trxTime := time.Now().Local()
	trxExpire := trxTime.Add(2 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceEchannel,
		MandiriBillBankTransferDetail: &midtrans.MandiriBillBankTransferDetail{
			BillInfo1: "complete payment",
			BillInfo2: "dept",
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 120,
			Unit:           "MINUTE",
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxID,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &midtrans.CustDetail{
			Phone: custPhone,
			FName: custName,
		},
		Items: &[]midtrans.ItemDetail{
			{
				ID:    itemID,
				Price: grossAmt,
				Qty:   1,
				Name:  itemName,
			},
		},
	}, trxTime, trxExpire
}

func (c *CoreGatewayMidtrans) ChargeReqPermataVirtualAccount(
	itemID, itemName, trxID string,
	custName, custPhone string,
	grossAmt int64,
) (*midtrans.ChargeReq, time.Time, time.Time) {
	trxTime := time.Now().Local()
	trxExpire := trxTime.Add(2 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankPermata,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 120,
			Unit:           "MINUTE",
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxID,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &midtrans.CustDetail{
			Phone: custPhone,
			FName: custName,
		},
		Items: &[]midtrans.ItemDetail{
			{
				ID:    itemID,
				Price: grossAmt,
				Qty:   1,
				Name:  itemName,
			},
		},
	}, trxTime, trxExpire
}

func (c *CoreGatewayMidtrans) ChargeReqBNIVirtualAccount(
	itemID, itemName, trxID string,
	custName, custPhone string,
	grossAmt int64,
) (*midtrans.ChargeReq, time.Time, time.Time) {
	trxTime := time.Now().Local()
	trxExpire := trxTime.Add(2 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBni,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 120,
			Unit:           "MINUTE",
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxID,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &midtrans.CustDetail{
			Phone: custPhone,
			FName: custName,
		},
		Items: &[]midtrans.ItemDetail{
			{
				ID:    itemID,
				Price: grossAmt,
				Qty:   1,
				Name:  itemName,
			},
		},
	}, trxTime, trxExpire
}

func (c *CoreGatewayMidtrans) ChargeReqBCAVirtualAccount(
	itemID, itemName, trxID string,
	custName, custPhone string,
	grossAmt int64,
) (*midtrans.ChargeReq, time.Time, time.Time) {
	trxTime := time.Now().Local()
	trxExpire := trxTime.Add(2 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBca,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 120,
			Unit:           "MINUTE",
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxID,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &midtrans.CustDetail{
			Phone: custPhone,
			FName: custName,
		},
		Items: &[]midtrans.ItemDetail{
			{
				ID:    itemID,
				Price: grossAmt,
				Qty:   1,
				Name:  itemName,
			},
		},
	}, trxTime, trxExpire
}

func (c *CoreGatewayMidtrans) ChargeReqBRIVirtualAccount(
	itemID, itemName, trxID string,
	custName, custPhone string,
	grossAmt int64,
) (*midtrans.ChargeReq, time.Time, time.Time) {
	trxTime := time.Now().Local()
	trxExpire := trxTime.Add(2 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBri,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 120,
			Unit:           "MINUTE",
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxID,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &midtrans.CustDetail{
			Phone: custPhone,
			FName: custName,
		},
		Items: &[]midtrans.ItemDetail{
			{
				ID:    itemID,
				Price: grossAmt,
				Qty:   1,
				Name:  itemName,
			},
		},
	}, trxTime, trxExpire
}

func (c *CoreGatewayMidtrans) RequestCharge(chargeReq *midtrans.ChargeReq) (*midtrans.Response, error) {
	response, err := c.Core.Charge(chargeReq)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
