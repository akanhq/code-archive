package main

import (
	"fmt"
	"testing"
)

//func maxProfit(prices []int) int {
//
//	var lowPrice, heightPrice int
//	var lowIndex int
//	if len(prices) > 0 {
//		lowPrice = prices[0]
//		for i, price := range prices {
//			if price < lowPrice {
//				lowPrice = price
//				lowIndex = i
//			}
//		}
//		heightPrice = prices[lowIndex]
//		for i := lowIndex; i < len(prices); i++ {
//			if prices[i] > heightPrice {
//				heightPrice = prices[i]
//			}
//		}
//	}
//
//	if lowPrice != 0 && heightPrice != 0 {
//		return heightPrice - lowPrice
//	}
//
//	return 0
//}

func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}

	var minPrice = prices[0]
	var maxProfit = 0
	for _, price := range prices {
		if price < minPrice {
			minPrice = price
		} else if price-minPrice > maxProfit {
			maxProfit = price - minPrice
		}
	}
	return maxProfit

}

func TestMaxProfit(t *testing.T) {

	prices := []int{2, 4, 1}
	fmt.Printf("计算价格：%d", maxProfit(prices))
}
