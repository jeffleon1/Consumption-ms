package infraestructure

import "encoding/json"

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
	MeterID            int64   `json:"meter_id"`
	Address            string  `json:"address"`
	Active             []int64 `json:"active"`
	ReactiveInductive  []int64 `json:"reactive_inductive"`
	ReactiveCapacitive []int64 `json:"reactive_capacitive"`
	Exported           []int64 `json:"exported"`
}
