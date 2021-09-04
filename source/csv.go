package source

import (
	"encoding/csv"
	"io"
)

// CSV2Map : takes a reader and returns an array of dictionaries, using the header row as the keys
func CSV2Map(reader io.Reader) ([]map[string]string, err) {
	r := csv.NewReader(reader)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}
