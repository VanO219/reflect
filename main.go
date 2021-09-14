package main

import (
	"bytes"
	"crypto/sha1"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
)

type String string

func (s *String) UnmarshalJSON(bs []byte) (err error) {
	buf := string(bs)
	buf = strings.TrimSpace(buf)
	switch buf[0] {
	case '"':
		buf = buf[1 : len(buf)-1]
		*s = String(buf)
	default:
		*s = String(buf)
	}
	return
}

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

var js = `
{
    "merchant": "ad25ef06-1824-413f-8ef1-c08115b9b979",
    "amount": 973,
    "order_id": 14425840,
    "description": "Заказ №14425840",
    "success_url": "http://myawesomesite.com/payment_success",
    "unix_timestamp": 1573451160,
    "salt": "dPUTLtbMfcTGzkaBnGtseKlcQymCLrYI",
    "testing": 1,
    "client_phone": "+7 912 9876543",
    "client_email": "test@test.ru",
    "receipt_contact": "test@mail.com",
    "receipt_items": [{"discount_sum": 40, "name": "Товар 1", "payment_method": "full_prepayment", "payment_object": "commodity", "price": 48, "quantity": 10, "sno": "osn", "vat": "vat10"}, {"name": "Товар 2", "payment_method": "full_prepayment", "payment_object": "commodity", "price": 533, "quantity": 1, "sno": "osn", "vat": "vat10"}]
}
`
var key = "00112233445566778899aabbccddeeff"

func main() {
	st := Request{}

	err := json.Unmarshal([]byte(js), &st)
	if err != nil {
		log.Fatalln(err)
	}

	t := reflect.TypeOf(st)
	v := reflect.ValueOf(st)
	b := []string{}
	mp := map[string]string{}

	for i := 0; i < t.NumField(); i++ {
		if len(v.Field(i).String()) == 0 {
			continue
		}
		b = append(b, strings.Split(t.Field(i).Tag.Get("json"), ",")[0])
		mp[b[i]] = b64.StdEncoding.EncodeToString([]byte(v.Field(i).String()))
	}
	bb := bytes.Buffer{}
	sort.Strings(b)
	ln := len(b) - 1
	for i, j := range b {
		bb.WriteString(j)
		bb.WriteString("=")
		bb.WriteString(mp[j])
		if i != ln {
			bb.WriteString("&")
		}
	}
	sh := sha1.New()
	sh1 := sha1.New()

	sh1.Write([]byte(key + bb.String()))
	sh.Write([]byte(key + fmt.Sprintf("%x", sh1.Sum(nil))))
	fmt.Printf("%x", sh.Sum(nil))
}

func structToMapAndSlice(input Request) (output string) {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	var mp = map[string]string{}
	var sl = []string{}
	for i := 0; i < t.NumField(); i++ {
		if len(v.Field(i).String()) == 0 {
			continue
		}
		sl = append(sl, strings.Split(t.Field(i).Tag.Get("json"), ",")[0])
		mp[sl[i]] = b64.StdEncoding.EncodeToString([]byte(v.Field(i).String()))
	}
	sort.Strings(sl)
	bb := bytes.Buffer{}
	ln := len(sl) - 1
	for i, j := range sl {
		bb.WriteString(j)
		bb.WriteString("=")
		bb.WriteString(mp[j])
		if i != ln {
			bb.WriteString("&")
		}
	}
	output = bb.String()
	return
}
