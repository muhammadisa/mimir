package mimir

import (
	"crypto/rand"
	"fmt"
	"github.com/veritrans/go-midtrans"
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
	ChargeReqPermataVirtualAccount(itemID, itemName, custName, custPhone string, grossAmt int64) *midtrans.ChargeReq
	ChargeReqMandiriBill(itemID, itemName, custName, custPhone string, grossAmt int64) *midtrans.ChargeReq
	ChargeReqBNIVirtualAccount(itemID, itemName, custName, custPhone string, grossAmt int64) *midtrans.ChargeReq
	ChargeReqBCAVirtualAccount(itemID, itemName, custName, custPhone string, grossAmt int64) *midtrans.ChargeReq
	ChargeReqBRIVirtualAccount(itemID, itemName, custName, custPhone string, grossAmt int64) *midtrans.ChargeReq
	RequestCharge(chargeReq *midtrans.ChargeReq) (*midtrans.Response, error)
}

func trxCodeGenerator() string {
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

type ConsumableResponse struct {
	Status     string
	Gross      string
	Message    string
	Bank       string
	VaNumber   string
	BillerCode string
	BillKey    string
}

func DetermineConsumableResponse(res *midtrans.Response) *ConsumableResponse {
	if res.PermataVaNumber != "" {
		return &ConsumableResponse{
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
			return &ConsumableResponse{
				Status:     res.StatusCode,
				Gross:      res.GrossAmount,
				Message:    res.StatusMessage,
				Bank:       res.Bank,
				VaNumber:   "",
				BillerCode: res.BillerCode,
				BillKey:    res.BillKey,
			}
		} else {
			return &ConsumableResponse{
				Status:  res.StatusCode,
				Gross:   res.GrossAmount,
				Message: res.StatusMessage,
				Bank:       res.VANumbers[0].Bank,
				VaNumber:   res.VANumbers[0].VANumber,
				BillerCode: res.BillerCode,
				BillKey:    res.BillKey,
			}
		}
	}
}

func (c *CoreGatewayMidtrans) ChargeReqMandiriBill(
	itemID, itemName string,
	custName, custPhone string,
	grossAmt int64,
) *midtrans.ChargeReq {
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceEchannel,
		MandiriBillBankTransferDetail: &midtrans.MandiriBillBankTransferDetail{
			BillInfo1: "complete payment",
			BillInfo2: "dept",
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxCodeGenerator(),
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
	}
}

func (c *CoreGatewayMidtrans) ChargeReqPermataVirtualAccount(
	itemID, itemName string,
	custName, custPhone string,
	grossAmt int64,
) *midtrans.ChargeReq {
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankPermata,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxCodeGenerator(),
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
	}
}

func (c *CoreGatewayMidtrans) ChargeReqBNIVirtualAccount(
	itemID, itemName string,
	custName, custPhone string,
	grossAmt int64,
) *midtrans.ChargeReq {
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBni,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxCodeGenerator(),
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
	}
}

func (c *CoreGatewayMidtrans) ChargeReqBCAVirtualAccount(
	itemID, itemName string,
	custName, custPhone string,
	grossAmt int64,
) *midtrans.ChargeReq {
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBca,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxCodeGenerator(),
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
	}
}

func (c *CoreGatewayMidtrans) ChargeReqBRIVirtualAccount(
	itemID, itemName string,
	custName, custPhone string,
	grossAmt int64,
) *midtrans.ChargeReq {
	return &midtrans.ChargeReq{
		PaymentType: midtrans.SourceBankTransfer,
		BankTransfer: &midtrans.BankTransferDetail{
			Bank: midtrans.BankBri,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  trxCodeGenerator(),
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
	}
}

func (c *CoreGatewayMidtrans) RequestCharge(chargeReq *midtrans.ChargeReq) (*midtrans.Response, error) {
	response, err := c.Core.Charge(chargeReq)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
