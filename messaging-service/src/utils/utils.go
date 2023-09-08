package utils

import (
	"encoding/json"
	"errors"
)

func GetEventType(event string) (string, error) {
	e := map[string]interface{}{}
	err := json.Unmarshal([]byte(event), &e)
	if err != nil {
		return "", err
	}

	eType, ok := e["eventType"]
	if !ok {
		return "", errors.New("no event type present")
	}
	val, ok := eType.(string)
	if !ok {
		return "", errors.New("could not cast to event type")
	}
	return val, nil
}

func GetEventToken(event string) (string, error) {
	e := map[string]interface{}{}
	err := json.Unmarshal([]byte(event), &e)
	if err != nil {
		return "", err
	}

	eToken, ok := e["token"]
	if !ok {
		return "", errors.New("no event token present")
	}
	val, ok := eToken.(string)
	if !ok {
		return "", errors.New("could not cast to event token")
	}
	return val, nil
}


