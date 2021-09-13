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
	//if pos := strings.Index(buf, `"`); pos >= 0 {
	//	buf = buf[pos + 1 : len(buf)-1]
	//}
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

//type Request struct {
//	Merchant       string          `json:"merchant,omitempty"`
//	Amount         float64         `json:"amount,omitempty"`
//	OrderID        int64           `json:"order_id,omitempty"`
//	Description    string          `json:"description,omitempty"`
//	SuccessUrl     string          `json:"success_url,omitempty"`
//	UnixTimestamp  uint64          `json:"unix_timestamp,omitempty"`
//	Salt           string          `json:"salt,omitempty"`
//	Testing        int64           `json:"testing,omitempty"`
//	ClientPhone    string          `json:"client_phone,omitempty"`
//	ClientEmail    string          `json:"client_email,omitempty"`
//	ReceiptContact string          `json:"receipt_contact,omitempty"`
//	ReceiptItems   json.RawMessage `json:"receipt_items,omitempty"`
//	CallbackUrl    string          `json:"callback_url,omitempty"`
//}

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
    "receipt_items": [
        {
            "discount_sum": 40,
            "name": "Товар 1",
            "payment_method": "full_prepayment",
            "payment_object": "commodity",
            "price": 48,
            "quantity": 10,
            "sno": "osn",
            "vat": "vat10"
        },
        {
            "name": "Товар 2",
            "payment_method": "full_prepayment",
            "payment_object": "commodity",
            "price": 533,
            "quantity": 1,
            "sno": "osn",
            "vat": "vat10"
        }
    ]
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
	//
	for i := 0; i < t.NumField(); i++ {
		if len(v.Field(i).String()) == 0 {
			continue
		}
		b = append(b, strings.Split(t.Field(i).Tag.Get("json"), ",")[0])
		mp[b[i]] = b64.StdEncoding.EncodeToString([]byte(strings.ToLower(v.Field(i).String())))
		//fmt.Println("tag:", strings.Replace(t.Field(i).Tag.Get("json"), ",omitempty", "", 1), v.Field(i))
		//fmt.Println("tag:", strings.Replace(t.Field(i).Tag.Get("json"), ",omitempty", "", 1), t.Field(i).Name, v.Field(i))
	}
	bb := bytes.Buffer{}
	sort.Strings(b)
	for _, j := range b {
		bb.WriteString(j)
		bb.WriteString("=")
		bb.WriteString(mp[j])
		bb.WriteString("&")
	}
	//fmt.Println(bb.String())

	h := sha1.New()
	h1 := sha1.New()

	h1.Write([]byte(key))
	h1.Write(bb.Bytes())
	h.Write([]byte(key))
	h.Write(h1.Sum(nil))
	fmt.Printf("%x", h.Sum(nil))
	fmt.Println()
	tms := `amount=OTcz&client_email=dGVzdEB0ZXN0LnJ1&client_phone=KzcgOTEyIDk4NzY1NDM=&description=0JfQsNC60LDQtyDihJYxNDQyNTg0MA==&merchant=YWQyNWVmMDYtMTgyNC00MTNmLThlZjEtYzA4MTE1YjliOTc5&order_id=MTQ0MjU4NDA=&receipt_contact=dGVzdEBtYWlsLmNvbQ==&receipt_items=W3siZGlzY291bnRfc3VtIjogNDAsICJuYW1lIjogItCi0L7QstCw0YAgMSIsICJwYXltZW50X21ldGhvZCI6ICJmdWxsX3ByZXBheW1lbnQiLCAicGF5bWVudF9vYmplY3QiOiAiY29tbW9kaXR5IiwgInByaWNlIjogNDgsICJxdWFudGl0eSI6IDEwLCAic25vIjogIm9zbiIsICJ2YXQiOiAidmF0MTAifSwgeyJuYW1lIjogItCi0L7QstCw0YAgMiIsICJwYXltZW50X21ldGhvZCI6ICJmdWxsX3ByZXBheW1lbnQiLCAicGF5bWVudF9vYmplY3QiOiAiY29tbW9kaXR5IiwgInByaWNlIjogNTMzLCAicXVhbnRpdHkiOiAxLCAic25vIjogIm9zbiIsICJ2YXQiOiAidmF0MTAifV0=&salt=ZFBVVEx0Yk1mY1RHemthQm5HdHNlS2xjUXltQ0xyWUk=&success_url=aHR0cDovL215YXdlc29tZXNpdGUuY29tL3BheW1lbnRfc3VjY2Vzcw==&testing=MQ==&unix_timestamp=MTU3MzQ1MTE2MA==`
	hh := sha1.New()
	hh1 := sha1.New()
	hh1.Write([]byte(key))
	hh1.Write([]byte(tms))
	hh.Write([]byte(key))
	hh.Write(hh1.Sum(nil))
	fmt.Printf("%x", hh.Sum(nil))
}

//func b64encode(input string) (output string, err error) {
//	output = input))
//	return
//}
