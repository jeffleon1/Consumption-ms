package infraestructure

import (
	"encoding/json"
	"fmt"

	"github.com/jeffleon1/consumption-ms/pkg/application"
)

func UnmarshalFilterConsumptionSerializer(data []byte) (FilterConsumptionSerializer, error) {
	var r FilterConsumptionSerializer
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FilterConsumptionSerializer) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FilterConsumptionSerializer struct {
	Period    []string    `json:"period"`
	DataGraph []DataGraph `json:"data_graph"`
}

type DataGraph struct {
	MeterID            int       `json:"meter_id"`
	Address            string    `json:"address"`
	Active             []float64 `json:"active"`
	ReactiveInductive  []float64 `json:"reactive_inductive"`
	ReactiveCapacitive []float64 `json:"reactive_capacitive"`
	Exported           []float64 `json:"exported"`
}

func (f *FilterConsumptionSerializer) ToFilterConsumptionSerializer(data []application.Serializer) {
	for _, values := range data {
		f.Period = values.Period
		f.DataGraph = append(f.DataGraph, DataGraph{
			MeterID:            values.MeterID,
			Address:            fmt.Sprintf("Mock address %d", values.MeterID),
			Active:             values.Active,
			ReactiveInductive:  values.ReactiveInductive,
			ReactiveCapacitive: values.ReactiveCapacitive,
			Exported:           values.Exported,
		})
	}
}
