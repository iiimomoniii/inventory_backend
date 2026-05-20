package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ─── Strict JSON Keys ──────────────────────────────────────

// StrictJSONKeys — ตรวจสอบว่า JSON keys เป็น camelCase ตาม struct tag เท่านั้น
// ใช้งาน: utils.StrictJSONKeys(data, MyStruct{})
func StrictJSONKeys(data []byte, sample interface{}) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	valid := map[string]struct{}{}
	t := reflect.TypeOf(sample)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		valid[strings.Split(tag, ",")[0]] = struct{}{}
	}

	for key := range raw {
		if _, ok := valid[key]; !ok {
			return fmt.Errorf("invalid key %q — use camelCase only", key)
		}
	}
	return nil
}

// StrictUnmarshal — UnmarshalJSON + StrictJSONKeys รวมเป็น function เดียว
// ใช้งาน: utils.StrictUnmarshal(data, &req)
func StrictUnmarshal[T any](data []byte, out *T) error {
	if err := StrictJSONKeys(data, *out); err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}
