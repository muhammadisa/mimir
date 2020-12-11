package mimir

import (
	"crypto/rand"
	"fmt"
	"github.com/veritrans/go-midtrans"
	"time"
)

type Midtrans struct {
	SK  string
	CK  string
	ENV midtrans.EnvironmentType
}

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

func TrxCodeGenerator() string {
	RandomCrypto, _ := rand.Prime(rand.Reader, 18)
	return fmt.Sprintf("TRX%d", RandomCrypto)
}

func (m *Midtrans) InitializeMidtransClient() *CoreGatewayMidtrans {
	midclient := midtrans.NewClient()
	midclient.ServerKey = m.SK
	midclient.ClientKey = m.CK
	midclient.APIEnvType = m.ENV
	return &CoreGatewayMidtrans{
		Core: midtrans.CoreGateway{
			Client: midclient,
		},
	}
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

func DetermineConsumableResponse(res *midtrans.Response, expire time.Time) *MidtransTransaction {
	if res.PermataVaNumber != "" {
		return &MidtransTransaction{
			ExpireAt:   expire,
			TrxType:    1,
			Status:     res.StatusCode,
			Gross:      res.GrossAmount,
			Message:    res.StatusMessage,
			Bank:       res.Bank,
			VaNumber:   res.PermataVaNumber,
			BillerCode: res.BillerCode,
			BillKey:    res.BillKey,
		}
	} else {
		if res.BillKey != "" && res.BillerCode != "" {
			return &MidtransTransaction{
				ExpireAt:   expire,
				TrxType:    1,
				Status:     res.StatusCode,
				Gross:      res.GrossAmount,
				Message:    res.StatusMessage,
				Bank:       res.Bank,
				VaNumber:   "",
				BillerCode: res.BillerCode,
				BillKey:    res.BillKey,
			}
		} else {
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
	trxExpire := trxTime.Add(1 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceEchannel,
		MandiriBillBankTransferDetail: &midtrans.MandiriBillBankTransferDetail{
			BillInfo1: "complete payment",
			BillInfo2: "dept",
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 8,
			Unit:           "HOUR",
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
	trxExpire := trxTime.Add(1 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankPermata,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 8,
			Unit:           "HOUR",
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
	trxExpire := trxTime.Add(1 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBni,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 8,
			Unit:           "HOUR",
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
	trxExpire := trxTime.Add(1 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBca,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 8,
			Unit:           "HOUR",
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
	trxExpire := trxTime.Add(1 * time.Hour)
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBri,
		},
		CustomExpiry: &midtrans.CustomExpiry{
			OrderTime:      trxTime.Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 8,
			Unit:           "HOUR",
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
