package models

import (
	"errors"
	"strings"
)

var (
	HTTPMethods = []string{"GET", "POST", "PUT", "DELETE"}
)

func parseHTTPMethod(fields []string) string {
	for _, field := range fields {
		name := strings.ToUpper(strings.TrimSpace(field))
		for _, hname := range HTTPMethods {
			if name == hname {
				return name
			}
		}
	}
	return "GET"
}

func parseURL(fields []string) (string, error) {
	for _, field := range fields {
		if strings.Contains(field, "/") {
			return strings.Trim(field, `" '`), nil
		}
	}
	return "", errors.New("url not found")
}
