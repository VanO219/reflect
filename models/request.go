package models

type Request struct {
	Merchant       String `json:"merchant,omitempty"`
	Amount         String `json:"amount,omitempty"`
	OrderID        String `json:"order_id,omitempty"`
	Description    String `json:"description,omitempty"`
	SuccessUrl     String `json:"success_url,omitempty"`
	UnixTimestamp  String `json:"unix_timestamp,omitempty"`
	Salt           String `json:"salt,omitempty"`
	Testing        String `json:"testing,omitempty"`
	ClientPhone    String `json:"client_phone,omitempty"`
	ClientEmail    String `json:"client_email,omitempty"`
	ReceiptContact String `json:"receipt_contact,omitempty"`
	ReceiptItems   String `json:"receipt_items,omitempty"`
	CallbackUrl    String `json:"callback_url,omitempty"`
}
