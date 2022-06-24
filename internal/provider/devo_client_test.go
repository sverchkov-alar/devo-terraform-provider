package provider

import (
	"fmt"
	"testing"
)

func TestGetAlert(t *testing.T) {
	alert, err := GetAlert("", "", "")
	if err != nil {
		println(err)
	}
	fmt.Print("Query: ")

	fmt.Println(alert[0])
}

func TestCreateAlert(t *testing.T) {

	var alert = Alert{
		Name:        "TerraformTestAlert",
		Message:     "Alert definition message created with alerting API",
		Description: "",
		Subcategory: "Alert",
		AlertCorrelationContext: AlertCorrelationContext{
			QuerySourceCode: "from vpc.aws.flow\n  select *",
			Priority:        "5",
			CorrelationTrigger: CorrelationTrigger{
				Kind: "each",
			},
		},
	}

	a, err := CreateAlert(alert, "", "")
	if err != nil {
		println(err)
	}
	fmt.Println(a)
	fmt.Println(a.Id)
}

func TestDeleteAlert(t *testing.T) {
	err := DeleteAlert("", "", "")
	if err != nil {
		fmt.Println(err)
	}

}

func TestUpdateAlert(t *testing.T) {

	var alert = Alert{
		Id:          "183599",
		Name:        "TerraformTestAlert123",
		Message:     "123123123 Updated Alert definition message created with alerting API",
		Description: "",
		Subcategory: "Alert",
		AlertCorrelationContext: AlertCorrelationContext{
			QuerySourceCode: "from vpc.aws.flow\n  select *",
			Priority:        "5",
			CorrelationTrigger: CorrelationTrigger{
				Kind: "each",
			},
		},
	}

	a, err := UpdateAlert(alert, "", "")
	if err != nil {
		println(err)
	}
	fmt.Println(a)
	fmt.Println(a.Id)
}
