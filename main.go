package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
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

type node struct {
	name       string
	expr       ast.Expr
	actionExpr []ast.Expr
}

var nodes map[string]*node
var promo Promo
var filters map[string]struct{}

func savePromo() {

	promoBody := []byte(`
	{
		"promoName":"Season 1",
		"validFrom":"20210501",
		"validTo":"20210531",
		"ruleSet":[
		{
		"ruleName":"One",
	            "rule":"sku == xlsSkuCol1 && member == Gold" ,
		"actions":[
		"Discount == 10/100 * price",
	 	"FinalPrice == price - discount"
		]
	            },
		{
		"ruleName":"Two",
		"rule":"member == xlsMemberCol1" ,
		"actions":[
		"Points == 3",
		"Flag == 2"
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

	var empty struct{}

	filters = make(map[string]struct{})

	// Unmarshall
	json.Unmarshal(promoBody, &promo)

	// Seup the dimensions
	for _, v := range promo.Filters {
		filters[v] = empty
	}

	// Create AST nodes
	compile()

}

func main() {

	savePromo()

	for {
		// fmt.Println("Press any key")
		// scanner := bufio.NewScanner(os.Stdin)
		// scanner.Scan()
		compute()
	}
}

func compile() {

	nodes = make(map[string]*node)

	for _, rule := range promo.RuleSet {

		ruleExpr, ruleErr := parser.ParseExpr(rule.Rule)
		if ruleErr != nil {
			log.Fatal(ruleErr)
		}

		spew.Dump(ruleExpr)

		nodes[rule.RuleName] = &node{name: rule.RuleName, expr: ruleExpr}

		for _, action := range rule.Actions {
			actionExpr, actionErr := parser.ParseExpr(action)

			spew.Dump(actionExpr)

			if actionErr != nil {
				log.Fatal(actionErr)
			}

			nodes[rule.RuleName].actionExpr = append(nodes[rule.RuleName].actionExpr, actionExpr)
		}

	}

}

func compute() {

	factBody := `{
        "sku" : "ADI123456WHTXXL",
        "member" : "Gold",
		"qty" : "3",
		"price" : "100"
    }`

	var facts map[string]string

	// Unmarshall
	json.Unmarshal([]byte(factBody), &facts)

	// Seyp trace
	facts["trace"] = ""

	// Traverse through all rules in the ruleset
	for _, v := range promo.RuleSet {

		fmt.Println("Before")

		for k, v := range facts {
			fmt.Printf("%s = %s\n", k, v)
		}

		exeNode(v.RuleName, v.Actions, &facts, &filters)

		fmt.Println()

		fmt.Println("After")

		for k, v := range facts {
			fmt.Printf("%s = %s\n", k, v)
		}

		fmt.Println()
	}

}

func exeNode(ruleName string, actions []string, facts *map[string]string, filters *map[string]struct{}) {

	nodex := nodes[ruleName]
	rule := nodex.expr

	if eval(rule, true, true, 0, facts, filters) == "true" {

		(*facts)["trace"] = (*facts)["trace"] + " -> " + ruleName

		for _, action := range nodex.actionExpr {
			eval(action, false, true, 0, facts, filters)
		}

		fmt.Println("Done")
	}

}

func eval(exp ast.Expr, isRule bool, isLeft bool, cnt int, facts *map[string]string, filters *map[string]struct{}) string {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return evalBinaryExpr(exp, isRule, isLeft, cnt, facts, filters)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			return exp.Value
		case token.STRING:
			return exp.Value
		}
	case *ast.ParenExpr:
		return eval(exp.X, isRule, isLeft, cnt, facts, filters)
	case *ast.Ident:

		// Assignment
		if isRule {
			return exp.Name
		} else {
			if exp.Name[0:1] == strings.ToUpper(string(exp.Name[0:1])) {
				return strings.ToLower(exp.Name)
			} else {
				v, exist := (*facts)[exp.Name]
				if exist {
					return fmt.Sprintf("%v", v)
				} else {
					return exp.Name
				}
			}
			// 	if !isRule {
			// 		if cnt == 1 {
			// 			return exp.Name
			// 		} else {
			// 			v, exist := (*facts)[exp.Name]
			// 			if exist {
			// 				return fmt.Sprintf("%v", v)
			// 			} else {
			// 				return exp.Name
			// 			}
			// 		}
			// 	} else {
			// 		v, exist := (*facts)[exp.Name]
			// 		if exist {
			// 			return fmt.Sprintf("%v", v)
			// 		} else {
			// 			return exp.Name
			// 		}
			// 	}
			// }
		}
	}
	return ""
}

func evalBinaryExpr(exp *ast.BinaryExpr, isRule bool, isLeft bool, cnt int, facts *map[string]string, filters *map[string]struct{}) string {

	left := eval(exp.X, isRule, true, cnt+1, facts, filters)
	right := eval(exp.Y, isRule, false, cnt+1, facts, filters)

	switch exp.Op {
	case token.ADD:
		leftFloat := strToFloat64(left)
		rightFloat := strToFloat64(right)

		ans := leftFloat + rightFloat

		return fmt.Sprintf("%.2f", ans)
	case token.SUB:
		leftFloat := strToFloat64(left)
		rightFloat := strToFloat64(right)

		ans := leftFloat - rightFloat

		return fmt.Sprintf("%.2f", ans)
	case token.MUL:
		leftFloat := strToFloat64(left)
		rightFloat := strToFloat64(right)

		ans := leftFloat * rightFloat

		return fmt.Sprintf("%.2f", ans)
	case token.QUO:
		leftFloat := strToFloat64(left)
		rightFloat := strToFloat64(right)

		ans := leftFloat / rightFloat

		return fmt.Sprintf("%.2f", ans)
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
		// Rule or Action
		if isRule {
			// Check Dimension
			if strings.HasPrefix(right, "xls") {
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

		} else {
			// Assignment
			(*facts)[strings.ToLower(string(left))] = right
		}

	}

	return ""
}

func strToFloat64(value string) float64 {
	floatNbr, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatalf("Unable to convert %v to float", value)
	}

	return floatNbr
}
