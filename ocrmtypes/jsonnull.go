package ocrmtypes

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"
)

type JsonNullString struct {
	sql.NullString
}

func NewJsonNullString(str string) *JsonNullString {
	return &JsonNullString{NullString: sql.NullString{
		String: str,
		Valid:  true,
	}}
}

func (j *JsonNullString) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.String)
	} else {
		return json.Marshal(nil)
	}
}

func (j *JsonNullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		j.Valid = true
		j.String = *x
	} else {
		j.Valid = false
	}
	return nil
}

type JsonNullBool struct {
	sql.NullBool
}

func (j *JsonNullBool) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.Bool)
	} else {
		return json.Marshal(nil)
	}
}

func (j *JsonNullBool) UnmarshalJSON(data []byte) error {
	var x *bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		j.Valid = true
		j.Bool = *x
	} else {
		j.Valid = false
	}
	return nil
}

func NewJsonNullBool(value bool) *JsonNullBool {
	return &JsonNullBool{NullBool: sql.NullBool{
		Bool:  value,
		Valid: true,
	}}
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (j *JsonNullInt64) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (j *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		j.Valid = true
		j.Int64 = *x
	} else {
		j.Valid = false
	}
	return nil
}

func NewJsonNullInt64(value int64) *JsonNullInt64 {
	return &JsonNullInt64{NullInt64: sql.NullInt64{
		Int64: value,
		Valid: true,
	}}
}

type JsonNullInt32 struct {
	sql.NullInt32
}

func (j *JsonNullInt32) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.Int32)
	} else {
		return json.Marshal(nil)
	}
}

func (j *JsonNullInt32) UnmarshalJSON(data []byte) error {
	var x *int32
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		j.Valid = true
		j.Int32 = *x
	} else {
		j.Valid = false
	}
	return nil
}

func NewJsonNullInt32(value int32) *JsonNullInt32 {
	return &JsonNullInt32{NullInt32: sql.NullInt32{
		Int32: value,
		Valid: true,
	}}
}

type JsonNullTime struct {
	sql.NullTime
}

func NewJsonNullTime(time time.Time) *JsonNullTime {
	return &JsonNullTime{
		NullTime: sql.NullTime{
			Time:  time,
			Valid: true,
		},
	}
}

func (j *JsonNullTime) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.Time.Unix() - 6*60*60)
	} else {
		return json.Marshal(nil)
	}
}

func (j *JsonNullTime) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		j.Valid = true
		j.Time = time.Unix(*x, 0)
	} else {
		j.Valid = false
	}
	return nil
}

func (j *JsonNullTime) GetDbTimeWithLocal() *JsonNullTime {
	if j.Valid {
		return NewJsonNullTime(time.Date(j.Time.Year(), j.Time.Month(), j.Time.Day(), j.Time.Hour(), j.Time.Minute(), j.Time.Second(), j.Time.Nanosecond(), time.Local))
	} else {
		return nil
	}
}

type JsonNullFloat64 struct {
	sql.NullFloat64
}

func NewJsonNullFloat64(value float64) *JsonNullFloat64 {
	return &JsonNullFloat64{
		NullFloat64: sql.NullFloat64{
			Float64: value,
			Valid:   true,
		},
	}
}

func (j *JsonNullFloat64) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.Float64)
	} else {
		return json.Marshal(nil)
	}
}

func (j *JsonNullFloat64) UnmarshalJSON(data []byte) error {
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		j.Valid = true
		j.Float64 = *x
	} else {
		j.Valid = false
	}
	return nil
}

func GetNotNilValue(nullable interface{}) interface{} {
	switch reflect.TypeOf(nullable) {
	case reflect.TypeOf(&JsonNullBool{}):
		if reflect.ValueOf(nullable).IsNil() {
			return false
		}
		value := nullable
		return value.(*JsonNullBool).Bool
	case reflect.TypeOf(&JsonNullFloat64{}):
		if reflect.ValueOf(nullable).IsNil() {
			return float64(0.0)
		}
		value := nullable
		return value.(*JsonNullFloat64).Float64
	case reflect.TypeOf(&JsonNullInt32{}):
		if reflect.ValueOf(nullable).IsNil() {
			return int32(0)
		}
		value := nullable
		return value.(*JsonNullInt32).Int32
	case reflect.TypeOf(&JsonNullInt64{}):
		if reflect.ValueOf(nullable).IsNil() {
			return int64(0)
		}
		value := nullable
		return value.(*JsonNullInt64).Int64
	case reflect.TypeOf(&JsonNullString{}):
		if reflect.ValueOf(nullable).IsNil() {
			return ""
		}
		value := nullable
		return value.(*JsonNullString).String
	case reflect.TypeOf(&JsonNullTime{}):
		if reflect.ValueOf(nullable).IsNil() {
			return time.Unix(0, 0)
		}
		value := nullable
		return value.(*JsonNullTime).Time
	}
	return 0
}
