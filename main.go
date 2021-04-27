package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type Promo struct {
	PromoName string    `json:"promoName"`
	ValidFrom string    `json:"validFrom"`
	ValidTo   string    `json:"validTo"`
	RuleSet   []ruleSet `json:"ruleSet"`
	Filters   []string  `json:"filters"`
}

type ruleSet struct {
	RuleName string   `json:"ruleName"`
	Rule     string   `json:"rule"`
	Actions  []string `json:"actions"`
}

// type rule struct {
// 	Name      string    `json:"name"`
// 	ValidFrom string    `json:"validFrom"`
// 	ValidTo   string    `json:"validTo"`
// 	RuleBox   []ruleBox `json:"rule"`
// 	Filters   []string  `json:"filters"`
// }

// type ruleBox struct {
// 	RuleCmd    string   `json:"ruleCmd"`
// 	RuleAction []string `json:"ruleAction"`
// }

func main() {

	promoBody := []byte(`
	{
		"promoName":"Season 1",
		"validFrom":"20210501",
		"validTo":"20210531",
		"ruleSet":[
		{
		"ruleName":"One",
	            "rule":"sku == xlsSkuCol1" ,
		"actions":[
		"Qty =  2",
		"Flag = 1"
		]
	            },
		{
		"ruleName":"Two",
		"rule":"member == xlsMemberCol1" ,
		"actions":[
		"Points = 2",
		"Flag = 2"
		]
	            }
		],
		"filters":[
		"xlsSkuCol1-ADI123456WHTXXL",
		"xlsSkuCol1-ADI123457WHTXXL",
		"xlsMemberCol1-Gold"
		]
		}

	`)

	// promoBody := []byte(`
	// {
	// 	"promoName":"Season 1",
	// 	"validFrom":"20210501",
	// 	"validTo":"20210531",
	// 	"ruleSet":[
	// 	{
	// 	"ruleName":"Two",
	//  	"rule":"member == xlsMemberCol1" ,
	// 	"actions":[
	// 	"points = points * 2",
	// 	"flag = 2"
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

	factBody := `{
        "sku" : "ADI123456WHTXXL",
        "member" : "Gold"
    }`

	var promo Promo
	var facts map[string]interface{}

	var empty struct{}
	filters := make(map[string]struct{})

	json.Unmarshal(promoBody, &promo)
	json.Unmarshal([]byte(factBody), &facts)

	for _, v := range promo.Filters {
		//	fmt.Printf("Filters: %s\n", v)
		filters[v] = empty
	}

	for _, v := range promo.RuleSet {

		fmt.Printf("RuleName: %s\n", v.RuleName)
		fmt.Printf("Rule: %s\n", v.Rule)

		for _, a := range v.Actions {
			fmt.Printf("Action: %s\n", a)
		}

		// strs := strings.Split(v.Rule, " ")
		// key := fmt.Sprintf("%s.%s", strs[2], facts[strs[0]])
		// fmt.Printf("key: %s\n", key)
		// _, ok := filters[key]

		// if ok {
		// 	fmt.Println("ok")
		// } else {
		// 	fmt.Println("nok")
		// }

		exe(v.Rule, v.Actions, &facts, &filters)

		fmt.Println()
	}

}

func exe(rule string, actions []string, facts *map[string]interface{}, filters *map[string]struct{}) {

	ruleToken, _ := parser.ParseExpr(rule)

	//actionTokens := make([]*ast.Expr, 0)

	if eval(ruleToken, facts, filters) == "true" {
		for _, v := range actions {
			actionToken, _ := parser.ParseExpr(v)
			eval(actionToken, facts, filters)
		}
		fmt.Println("Done")
	}

	for k, v := range facts {

	}
}

func eval(exp ast.Expr, facts *map[string]interface{}, filters *map[string]struct{}) string {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return evalBinaryExpr(exp, facts, filters)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			return exp.Value
		case token.STRING:
			return strings.ReplaceAll(exp.Value, "\"", "")
		}
	case *ast.ParenExpr:
		return eval(exp.X, facts, filters)
	case *ast.Ident:
		// v, exist := (*facts)[exp.Name]
		// if exist {
		// 	return fmt.Sprintf("%v", v)
		// } else {
		// 	return exp.Name
		// }
		return exp.Name
	}

	return ""
}

func evalBinaryExpr(exp *ast.BinaryExpr, facts *map[string]interface{}, filters *map[string]struct{}) string {
	left := eval(exp.X, facts, filters)
	right := eval(exp.Y, facts, filters)

	switch exp.Op {
	case token.ADD:
		l, _ := strconv.Atoi(left)
		r, _ := strconv.Atoi(right)
		return strconv.Itoa(l + r)
	case token.SUB:
		l, _ := strconv.Atoi(left)
		r, _ := strconv.Atoi(right)
		return strconv.Itoa(l - r)
	case token.MUL:
		l, _ := strconv.Atoi(left)
		r, _ := strconv.Atoi(right)
		return strconv.Itoa(l * r)
	case token.QUO:
		l, _ := strconv.Atoi(left)
		r, _ := strconv.Atoi(right)
		return strconv.Itoa(l / r)
	case token.LAND:
		if left == "true" && right == "true" {
			return "true"
		} else {
			return "false"
		}
	case token.LOR:
		if left == "true" || right == "true" {
			return "true"
		} else {
			return "false"
		}
	case token.EQL:
		// Assign
		if left[0:1] == strings.ToUpper(string(left[0:1])) {
			(*facts)[strings.ToLower(string(left))] = right
		} else if right[0:3] == "xls" {
			v, exist := (*facts)[left]
			if exist {
				key := fmt.Sprintf("%s-%v", right, v)
				_, x := (*filters)[key]
				if x {
					return "true"
				} else {
					return "false"
				}

			} else {
				return "false"
			}
		} else {
			isEql := (*facts)[left] == right
			if isEql {
				return "true"
			} else {
				return "false"
			}
		}

	}

	return ""
}
