package main

import (
	"log"
	"testing"
)

// If sku = ADI123456WHTXXL && Member == Gold
// Then Discount 10%
// Net Price will be Price 0 Discount

// qty = 3
// flag = 2
// trace =  -> One -> Two
// discount = 10.00
// finalprice = 90.00
// points = 3
// sku = ADI123456WHTXXL
// member = Gold
// price = 100

func Test1(t *testing.T) {

	brePackage := []byte(`
	{
		"promoName":"Season 1",
		"validFrom":"20210501",
		"validTo":"20210531",
		"ruleSet":	[
						{
						"ruleName":"One",
						"rule":"sku == xlsSkuCol1" ,
						"actions":	[
									"_discount == 10/100 * price",
									"_netPrice == price - discount",
									"_flag == 0"
									]
	    				},
						{
						"ruleName":"Two",
						"rule":"member == Gold" ,
						"actions":	[
									"_points == 3",
									"_flag == 1"
									]
						},
						{
						"ruleName":"Three",
						"rule":"flag != 1" ,
						"actions":	[
									"_points == 2"
									]
						}

					],
					"filters":[
					"xlsSkuCol1-ADI123456WHTXXL",
					"xlsSkuCol1-ADI123457WHTXXL",
					"xlsMemberCol1-Gold"
					]
	}`)

	factBody := []byte(`{
	     "sku" : "ADI123456WHTXXL",
	     "member" : "Gold",
	 	 "qty" : "3",
	 	 "price" : "100"
	 }`)

	useBrePackage(brePackage)

	facts, err := process(factBody)
	if err != nil {
		log.Fatal(err)
	}
	price := "100"
	discount := "10.00"
	netPrice := "90.00"
	points := "3"

	if facts["price"] != price {
		t.Errorf("Price Error, should be %s; got %s", price, facts["price"])
	}

	if facts["discount"] != discount {
		t.Errorf("Discount  Error, should be %s; got %s", discount, facts["discount"])
	}

	if facts["netPrice"] != netPrice {
		t.Errorf("Net Price Error, should be %s; got %s", netPrice, facts["netPrice"])
	}

	if facts["points"] != points {
		t.Errorf("Net Price Error, should be %s; got %s", points, facts["points"])
	}

}
