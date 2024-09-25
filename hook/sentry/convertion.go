package sentry

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/XiBao/logger/common"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

func (h Hook) convertEvent(e *zerolog.Event, level zerolog.Level, message string) (sentry.Event, error) {
	var record sentry.Event

	record.Level = sentry.Level(level.String())
	record.Message = message
	record.Timestamp = zerolog.TimestampFunc()
	fields := convertFields(e)
	record.Extra = make(map[string]interface{}, len(fields))
	var retErr error
	for k, v := range fields {
		switch k {
		case zerolog.ErrorFieldName:
			if err, ok := v.(error); ok {
				if retErr == nil {
					retErr = err
				} else {
					retErr = errors.Join(retErr, err)
				}
				record.SetException(err, -1)
			} else {
				record.Exception = append(record.Exception, sentry.Exception{
					Value:      convertValue(v),
					Stacktrace: common.Stacktrace(),
				})
			}
		default:
			record.Extra[k] = v
		}
	}
	return record, retErr
}

// convertFields extracts and converts zerolog event fields to OpenTelemetry key-value pairs.
//
// This function iterates over all fields present in a zerolog event, converting each field
// to an OpenTelemetry log.KeyValue structure. The conversion process is handled by the
// convertValue function, which adapts the field's value to the appropriate OpenTelemetry
// log.Value type based on the value's underlying type.
//
// Parameters:
// - e *zerolog.Event: The zerolog event containing the fields to be converted.
//
// Returns:
// - map[string]interface: A map of event fields key values representing the converted fields.
func convertFields(e *zerolog.Event) map[string]interface{} {
	ev := fmt.Sprintf("%s}", reflect.ValueOf(e).Elem().FieldByName("buf"))
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(ev), &data); err != nil {
		return nil
	}

	return data
}

func convertValue(v interface{}) string {
	switch v := v.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case []byte:
		return strings.ToValidUTF8(string(v), "")
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case uint:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	case uint64:
		return strconv.FormatUint(v, 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	}

	t := reflect.TypeOf(v)
	if t == nil {
		return ""
	}
	val := reflect.ValueOf(v)
	switch t.Kind() { //nolint:exhaustive // We only handle the types we care about.
	case reflect.Struct:
		return fmt.Sprintf("%+v", v)
	case reflect.Slice, reflect.Array:
		items := make([]string, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			items = append(items, convertValue(val.Index(i).Interface()))
		}
		bs, _ := json.Marshal(items)
		return string(bs)
	case reflect.Map:
		kvs := make(map[string]string, val.Len())
		for _, k := range val.MapKeys() {
			var key string
			if k.Kind() == reflect.Struct {
				key = fmt.Sprintf("%+v", k.Interface())
			} else {
				key = fmt.Sprintf("%v", k.Interface())
			}
			kvs[key] = convertValue(val.MapIndex(k).Interface())
		}
		bs, _ := json.Marshal(kvs)
		return string(bs)
	case reflect.Ptr, reflect.Interface:
		return convertValue(val.Elem().Interface())
	}

	return fmt.Sprintf("unhandled attribute type: (%s) %+v", t, v)
}
