package cassandra

import (
	"errors"
	"fmt"
	"reflect"
	r "reflect"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
)

type Counter int

type SchemaConverter struct {
	Enums map[string]gocql.Type
}

func (s *SchemaConverter) stringTypeOf(i interface{}) (string, error) {
	fmt.Printf("Type: %v", reflect.ValueOf(i).Kind())
	_, isByteSlice := i.([]byte)
	if !isByteSlice {
		// Check if we found a higher kinded type
		switch reflect.ValueOf(i).Kind() {
		case reflect.Slice:
			elemVal := reflect.Indirect(reflect.New(reflect.TypeOf(i).Elem())).Interface()
			ct := s.cassaType(elemVal)
			//if isIn := funk.Contains(s.Enums, ct.String()); isIn {
			//	ct = s.Enums[ct.String()]
			//}
			if ct == gocql.TypeCustom {
				return "", fmt.Errorf("unsupported type %T", i)
			}
			return fmt.Sprintf("list<%v>", ct), nil
		case reflect.Map:
			keyVal := reflect.Indirect(reflect.New(reflect.TypeOf(i).Key())).Interface()
			elemVal := reflect.Indirect(reflect.New(reflect.TypeOf(i).Elem())).Interface()
			keyCt := s.cassaType(keyVal)
			elemCt := s.cassaType(elemVal)
			//if isIn := funk.Contains(s.Enums, keyCt.String()); isIn {
			//	keyCt = s.Enums[keyCt.String()]
			//}
			//if isIn := funk.Contains(s.Enums, elemCt.String()); isIn {
			//	elemCt = s.Enums[elemCt.String()]
			//}
			if keyCt == gocql.TypeCustom || elemCt == gocql.TypeCustom {
				return "", fmt.Errorf("unsupported map key or value type %T", i)
			}
			return fmt.Sprintf("map<%v, %v>", keyCt, elemCt), nil
		}
	}
	ct := s.cassaType(i)
	//if isIn := funk.Contains(s.Enums, ct.String()); isIn {
	//	ct = s.Enums[ct.String()]
	//}
	if ct == gocql.TypeCustom {
		return "", fmt.Errorf("unsupported type %T", i)
	}
	return cassaTypeToString(ct)
}

func (s *SchemaConverter) cassaType(i interface{}) gocql.Type {
	switch i.(type) {
	case int, int32:
		return gocql.TypeInt
	case int64:
		return gocql.TypeBigInt
	case string:
		return gocql.TypeText
	case float32:
		return gocql.TypeFloat
	case float64:
		return gocql.TypeDouble
	case bool:
		return gocql.TypeBoolean
	case time.Time:
		return gocql.TypeTimestamp
	case gocql.UUID:
		return gocql.TypeUUID
	case []byte:
		return gocql.TypeBlob
	case Counter:
		return gocql.TypeCounter
	}
	return gocql.TypeCustom
}

func cassaTypeToString(t gocql.Type) (string, error) {
	switch t {
	case gocql.TypeInt:
		return "int", nil
	case gocql.TypeBigInt:
		return "bigint", nil
	case gocql.TypeVarchar:
		return "varchar", nil
	case gocql.TypeFloat:
		return "float", nil
	case gocql.TypeDouble:
		return "double", nil
	case gocql.TypeBoolean:
		return "boolean", nil
	case gocql.TypeTimestamp:
		return "timestamp", nil
	case gocql.TypeUUID:
		return "uuid", nil
	case gocql.TypeBlob:
		return "blob", nil
	case gocql.TypeCounter:
		return "counter", nil
	default:
		return "", errors.New("unkown cassandra type")
	}
}

func (s *SchemaConverter) StructToMap(val interface{}) (map[string]interface{}, bool) {
	// indirect so function works with both structs and pointers to them
	structVal := r.Indirect(r.ValueOf(val))
	kind := structVal.Kind()
	if kind != r.Struct {
		return nil, false
	}
	sinfo := getStructInfo(structVal)
	mapVal := make(map[string]interface{}, len(sinfo.FieldsList))
	for _, field := range sinfo.FieldsList {
		if structVal.Field(field.Num).CanInterface() {
			mapVal[field.Key] = structVal.Field(field.Num).Interface()
		}
	}
	return mapVal, true
}

// MapToStruct converts a map to a struct. It is the inverse of the StructToMap
// function. For details see StructToMap.
func (s *SchemaConverter) MapToStruct(m map[string]interface{}, struc interface{}) error {
	val := r.Indirect(r.ValueOf(struc))
	sinfo := getStructInfo(val)
	for k, v := range m {
		if info, ok := sinfo.FieldsMap[k]; ok {
			structField := val.Field(info.Num)
			if structField.Type().Name() == r.TypeOf(v).Name() {
				structField.Set(r.ValueOf(v))
			}
		}
	}
	return nil
}

// FieldsAndValues returns a list field names and a corresponding list of values
// for the given struct. For details on how the field names are determined please
// see StructToMap.
func (s *SchemaConverter) FieldsAndValues(val interface{}) ([]string, []interface{}, bool) {
	// indirect so function works with both structs and pointers to them
	structVal := r.Indirect(r.ValueOf(val))
	kind := structVal.Kind()
	if kind != r.Struct {
		return nil, nil, false
	}
	sinfo := getStructInfo(structVal)
	fields := make([]string, len(sinfo.FieldsList))
	values := make([]interface{}, len(sinfo.FieldsList))
	for i, info := range sinfo.FieldsList {
		field := structVal.Field(info.Num)
		fields[i] = info.Key
		values[i] = field.Interface()
	}
	return fields, values, true
}

var structMapMutex sync.RWMutex
var structMap = make(map[r.Type]*structInfo)

type fieldInfo struct {
	Key string
	Num int
}

type structInfo struct {
	// FieldsMap is used to access fields by their key
	FieldsMap map[string]fieldInfo
	// FieldsList allows iteration over the fields in their struct order.
	FieldsList []fieldInfo
}

func getStructInfo(v r.Value) *structInfo {
	st := r.Indirect(v).Type()
	structMapMutex.RLock()
	sinfo, found := structMap[st]
	structMapMutex.RUnlock()
	if found {
		return sinfo
	}

	n := st.NumField()
	fieldsMap := make(map[string]fieldInfo, n)
	fieldsList := make([]fieldInfo, 0, n)
	for i := 0; i != n; i++ {
		field := st.Field(i)
		info := fieldInfo{Num: i}
		tag := field.Tag.Get("cql")
		// If there is no cql specific tag and there are no other tags
		// set the cql tag to the whole field tag
		if tag == "" && strings.Index(string(field.Tag), ":") < 0 {
			tag = string(field.Tag)
		}
		if tag != "" {
			info.Key = tag
		} else {
			info.Key = field.Name
		}

		if _, found = fieldsMap[info.Key]; found {
			msg := fmt.Sprintf("Duplicated key '%s' in struct %s", info.Key, st.String())
			panic(msg)
		}

		fieldsList = append(fieldsList, info)
		fieldsMap[info.Key] = info
	}
	sinfo = &structInfo{
		fieldsMap,
		fieldsList,
	}
	structMapMutex.Lock()
	structMap[st] = sinfo
	structMapMutex.Unlock()
	return sinfo
}
