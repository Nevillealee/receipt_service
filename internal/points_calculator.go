package internal

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CalculatePoints(receipt *Receipt) int {
	totalPoints := 0

	alphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`)
	retailerChars := alphanumeric.FindAllString(receipt.Retailer, -1)
	totalPoints += len(retailerChars)

	total, err := strconv.ParseFloat(receipt.Total, 64)

	if err == nil {
		totalCents := int(total * 100)
		if totalCents%100 == 0 {
			totalPoints += 50
		}

		if totalCents%25 == 0 {
			totalPoints += 25
		}
	}

	totalPoints += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)

		if len(desc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)

			if err == nil {
				points := math.Ceil(price * 0.2)
				totalPoints += int(points)
			}
		}
	}

	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)

	if err == nil {
		if date.Day()%2 == 1 {
			totalPoints += 6
		}
	}

	t, err := time.Parse("15:04", receipt.PurchaseTime)

	if err == nil {
		purchaseTime := t.Hour()*60 + t.Minute()

		if purchaseTime > 14*60 && purchaseTime < 16*60 {
			totalPoints += 10
		}
	}

	return totalPoints
}
