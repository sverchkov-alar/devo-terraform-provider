package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type CorrelationTrigger struct {
	Kind string `json:"kind"`
}

type AlertCorrelationContext struct {
	QuerySourceCode    string             `json:"querySourceCode,omitempty"`
	Priority           string             `json:"priority,omitempty"`
	CorrelationTrigger CorrelationTrigger `json:"correlationTrigger,omitempty"`
	ExternalOffset     string             `json:"externalOffset,omitempty"`
	InternalPeriod     string             `json:"internalPeriod,omitempty"`
	InternalOffset     string             `json:"internalOffset,omitempty"`
	Period             string             `json:"period,omitempty"`
	Threshold          string             `json:"threshold,omitempty"`
	BackPeriod         string             `json:"backPeriod,omitempty"`
	Absolute           string             `json:"absolute,omitempty"`
	AggregationColumn  string             `json:"aggregationColumn,omitempty"`
}

type Alert struct {
	Id                      string                  `json:"id,omitempty"`
	Name                    string                  `json:"name,omitempty"`
	Message                 string                  `json:"message,omitempty"`
	Description             string                  `json:"description,omitempty"`
	Subcategory             string                  `json:"subcategory,omitempty"`
	AlertCorrelationContext AlertCorrelationContext `json:"alertCorrelationContext,omitempty"`
}

func GetAlert(token string, endpoint string, id string) ([]Alert, error) {
	var alert []Alert
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint+"/alerts/v1/alertDefinitions", nil)
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["standAloneToken"] = []string{token}
	query := req.URL.Query()
	query.Add("size", "1000")
	query.Add("page", "0")
	query.Add("idFilter", id)
	req.URL.RawQuery = query.Encode()
	res, err := client.Do(req)
	if err != nil {
		return alert, err
	}
	defer res.Body.Close()
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return alert, err
	}
	json.Unmarshal(respBody, &alert)
	return alert, nil
}

func DeleteAlert(token string, endpoint string, id string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", endpoint+"/alerts/v1/alertDefinitions", nil)
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["standAloneToken"] = []string{token}
	query := req.URL.Query()
	query.Add("alertIds", id)
	req.URL.RawQuery = query.Encode()
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("request: failed to remove the an alert")
	}
	return nil
}

func CreateAlert(alert Alert, token string, endpoint string) (Alert, error) {
	body, _ := json.Marshal(alert)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", endpoint+"/alerts/v1/alertDefinitions", bytes.NewReader(body))
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["standAloneToken"] = []string{token}
	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return alert, errors.New("error for creating an alert")
	}
	if err != nil {
		return alert, err
	}
	defer res.Body.Close()
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return alert, err
	}
	json.Unmarshal(respBody, &alert)
	return alert, nil

}

func UpdateAlert(alert Alert, token string, endpoint string) (Alert, error) {
	body, _ := json.Marshal(alert)

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", endpoint+"/alerts/v1/alertDefinitions", bytes.NewReader(body))
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["standAloneToken"] = []string{token}
	query := req.URL.Query()
	query.Add("id", alert.Id)
	req.URL.RawQuery = query.Encode()
	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return alert, errors.New("error for updating an alert")
	}
	if err != nil {
		return alert, err
	}
	defer res.Body.Close()
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return alert, err
	}
	json.Unmarshal(respBody, &alert)
	return alert, nil

}
