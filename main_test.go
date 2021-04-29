package main

import (
	"testing"
)

// If sku = ADI123456WHTXXL && Member == Gold
// Then Discount 10%
func Test1(t *testing.T) {

	// promoBody := []byte(`
	// {
	// 	"promoName":"Season 1",
	// 	"validFrom":"20210501",
	// 	"validTo":"20210531",
	// 	"ruleSet":[
	// 	{
	// 	"ruleName":"One",
	//             "rule":"sku == xlsSkuCol1 && member == Gold" ,
	// 	"actions":[
	// 	"discount == 10/100 * price"
	// 	"finalPrice = price - discount"
	// 	]
	//             },
	// 	{
	// 	"ruleName":"Two",
	// 	"rule":"member == xlsMemberCol1" ,
	// 	"actions":[
	// 	"Points == 3",
	// 	"Flag == 2"
	// 	]
	//             }
	// 	],
	// 	"filters":[
	// 	"xlsSkuCol1-ADI123456WHTXXL",
	// 	"xlsSkuCol1-ADI123457WHTXXL",
	// 	"xlsMemberCol1-Gold"
	// 	]
	// 	}

	// `)

	// factBody := []byte(`{
	//      "sku" : "ADI123456WHTXXL",
	//      "member" : "Gold",
	//  	 "qty" : "3",
	//  	 "price" : "100"
	//  }`)

	//savePromo(promoBody)

	//facts := compute(factBody)

	// price := "90.00"
	// discount := "10.00"

	// if facts["price"] != price {
	// 	t.Errorf("Price Error, should be %s; got %s", price, facts["price"])
	// }

	// if facts["discount"] != discount {
	// 	t.Errorf("Discount  Error, should be %s; got %s", facts["discount"], discount)
	// }

}
