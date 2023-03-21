package debugfmt

import (
	"encoding/json"
	"fmt"
)

func JsonMarshalIndent(v any, fmtStr string, fmtV ...any) {
	if fmtStr != "" {
		fmt.Printf(fmtStr, fmtV)
	}
	by, _ := json.MarshalIndent(v, "", "    ")
	fmt.Println(string(by))
}
