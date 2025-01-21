package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var DateFotmat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {

	var nextDate time.Time

	parsedDate, err := time.Parse(DateFotmat, date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("repeat rull is null")
	}
	// разбираем паттерн
	repeatPattern := strings.Fields(repeat)

	switch repeatPattern[0] {
	// раз в год
	case "y":
		nextDate = parsedDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	// дни
	case "d":
		// ошибка если отсутствуют дни для правила
		if len(repeatPattern) == 1 {
			return "", errors.New("repeate rule is invalid")

		}
		days, err := strconv.Atoi(repeatPattern[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("repeate rule is invalid")
		}
		nextDate = parsedDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	default:
		return "", errors.New("repeate rule is invalid")
	}
	return nextDate.Format(DateFotmat), nil
}
