package service

import (
	"testing"
	"time"

	"github.com/goProjects/loan_app/app/model"
	"github.com/stretchr/testify/assert"
)

func TestCalculateInstallmentSchedule(t *testing.T) {
	tests := []struct {
		amount          int64
		term            int64
		frequencyInDays int64
		expected        []model.Installment
	}{
		{
			amount:          100,
			term:            3,
			frequencyInDays: 4,
			expected: []model.Installment{
				{AmountInCents: 33, DueDate: time.Now().AddDate(0, 0, 4), SerialNo: 1},
				{AmountInCents: 33, DueDate: time.Now().AddDate(0, 0, 8), SerialNo: 2},
				{AmountInCents: 34, DueDate: time.Now().AddDate(0, 0, 12), SerialNo: 3},
			},
		},
		{
			amount:          200,
			term:            4,
			frequencyInDays: 5,
			expected: []model.Installment{
				{AmountInCents: 50, DueDate: time.Now().AddDate(0, 0, 5), SerialNo: 1},
				{AmountInCents: 50, DueDate: time.Now().AddDate(0, 0, 10), SerialNo: 2},
				{AmountInCents: 50, DueDate: time.Now().AddDate(0, 0, 15), SerialNo: 3},
				{AmountInCents: 50, DueDate: time.Now().AddDate(0, 0, 20), SerialNo: 4},
			},
		},
		{
			amount:          105,
			term:            5,
			frequencyInDays: 3,
			expected: []model.Installment{
				{AmountInCents: 21, DueDate: time.Now().AddDate(0, 0, 3), SerialNo: 1},
				{AmountInCents: 21, DueDate: time.Now().AddDate(0, 0, 6), SerialNo: 2},
				{AmountInCents: 21, DueDate: time.Now().AddDate(0, 0, 9), SerialNo: 3},
				{AmountInCents: 21, DueDate: time.Now().AddDate(0, 0, 12), SerialNo: 4},
				{AmountInCents: 21, DueDate: time.Now().AddDate(0, 0, 15), SerialNo: 5},
			},
		},
		{
			amount:          10,
			term:            4,
			frequencyInDays: 4,
			expected: []model.Installment{
				{AmountInCents: 2, DueDate: time.Now().AddDate(0, 0, 4), SerialNo: 1},
				{AmountInCents: 2, DueDate: time.Now().AddDate(0, 0, 8), SerialNo: 2},
				{AmountInCents: 3, DueDate: time.Now().AddDate(0, 0, 12), SerialNo: 3},
				{AmountInCents: 3, DueDate: time.Now().AddDate(0, 0, 16), SerialNo: 4},
			},
		},
	}

	for _, test := range tests {
		actual := calculateInstallmentSchedule(test.amount, test.term, test.frequencyInDays)
		for i := range actual {
			assert.Equal(t, test.expected[i].AmountInCents, actual[i].AmountInCents)
			assert.Equal(t, test.expected[i].SerialNo, actual[i].SerialNo)
			assert.True(t, datesAreEqual(test.expected[i].DueDate, actual[i].DueDate))
		}
	}
}

func datesAreEqual(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() && date1.Month() == date2.Month() && date1.Day() == date2.Day()
}
