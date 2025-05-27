package db

import "time"

func GetArticle(dateParam, slugParam string) ([]byte, error) {
	_, err := time.Parse("20060102", dateParam)
	if err != nil {
		return nil, err
	}

	switch slugParam {
	case "demo":
		return []byte("demo article text"), nil
	case "demo2":
		return []byte("second demo article text"), nil
	}

	return nil, nil
}
