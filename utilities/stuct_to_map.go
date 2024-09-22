package utilities

import "encoding/json"

func StructToMap(s interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	databytes, err := json.Marshal(s)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(databytes, &m)

	if err != nil {
		return nil, err
	}

	return m, nil
}
