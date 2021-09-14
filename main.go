package main

import (
	"bytes"
	"crypto/sha1"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"myprogs/reflectall/models"
	"reflect"
	"sort"
	"strings"
)

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
	//st := Request{}
	mp := map[string]models.String{}
	err := json.Unmarshal([]byte(js), &mp)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(stringToSSH1(key, mapToSortString(mp)))

	//fmt.Println(stringToSSH1(key, structToSortString(st)))
}

func structToSortString(input models.Request) (output string) {
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

func mapToSortString(input map[string]models.String) (output string) {
	//printMap(input)
	sl := []string{}
	mp := map[string]string{}
	for i, j := range input {
		if i == `signature` || len(j) == 0 {
			continue
		}
		i = strings.ToLower(i)
		sl = append(sl, i)
		mp[i] = b64.StdEncoding.EncodeToString([]byte(j))
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

func stringToSSH1(key, input string) (signature string) {
	sh := sha1.New()
	sh1 := sha1.New()
	sh1.Write([]byte(key + input))
	sh.Write([]byte(key + fmt.Sprintf("%x", sh1.Sum(nil))))
	signature = fmt.Sprintf("%x", sh.Sum(nil))
	return
}

func printMap(input map[string]models.String) {
	for i, j := range input {
		if len(j) == 0 {
			continue
		}
		fmt.Println(i, j)
	}
}
