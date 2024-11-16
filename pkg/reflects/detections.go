package reflects

import (
	"reflect"
)

var DefaultIgnoreFiledChange = map[string]bool{
	"CreatedAt": true,
	"UpdatedAt": true,
	"DeletedAt": true,
	"CreatedBy": true,
	"UpdatedBy": true,
	"DeletedBy": true,
}

// the usecase when use structural
func StructChangeDetection(
	beforeStruct interface{},
	afterStruct interface{},
	bypassFiled map[string]bool,
) (changedField map[string][2]interface{}) {
	if afterStruct == nil {
		return StructFields(beforeStruct, bypassFiled)
	}
	beforeType := reflect.TypeOf(beforeStruct).Elem()
	beforeValue := reflect.ValueOf(beforeStruct).Elem()
	afterType := reflect.TypeOf(afterStruct).Elem()
	afterValue := reflect.ValueOf(afterStruct).Elem()
	changedField = make(map[string][2]interface{})
	for i := 0; i < beforeType.NumField(); i++ {
		if bypassFiled[afterType.Field(i).Name] {
			continue
		}
		compareValueBefore := beforeValue.Field(i)
		compareValueAfter := afterValue.Field(i)
		switch beforeValue.Field(i).Kind() {
		case reflect.Struct:
			bfValue := compareValueBefore.Addr()
			afValue := compareValueBefore.Addr()
			// recurse if have struct nested on the same field name
			StructChangeDetection(bfValue.Interface(), afValue.Interface(), bypassFiled)
			if !reflect.DeepEqual(compareValueBefore.Interface(), compareValueAfter.Interface()) {
				changedField[afterType.Field(i).Name] = [2]interface{}{
					compareValueBefore.Interface(), compareValueAfter.Interface(),
				}
			}
			// if slice changed then return slice changed
		case reflect.Slice:
			if !reflect.DeepEqual(compareValueBefore.Interface(), compareValueAfter.Interface()) {
				changedField[afterType.Field(i).Name] = [2]interface{}{
					compareValueBefore.Interface(), compareValueAfter.Interface(),
				}
			}
			// if array changed detection then return array changed
		case reflect.Array:
			if !reflect.DeepEqual(compareValueBefore.Interface(), compareValueAfter.Interface()) {
				changedField[afterType.Field(i).Name] = [2]interface{}{
					compareValueBefore.Interface(), compareValueAfter.Interface(),
				}
			}
			// default case field of the struct have a type default by golang
		default:
			if compareValueBefore.Interface() != compareValueAfter.Interface() {
				changedField[afterType.Field(i).Name] = [2]interface{}{
					compareValueBefore.Interface(), compareValueAfter.Interface(),
				}
			}
		}
	}
	return
}

// the usecase when use structural
func StructFields(
	input interface{},
	bypassFiled map[string]bool,
) (changedField map[string][2]interface{}) {
	inputType := reflect.TypeOf(input).Elem()
	inputValue := reflect.ValueOf(input).Elem()
	changedField = make(map[string][2]interface{})
	for i := 0; i < inputType.NumField(); i++ {
		if bypassFiled[inputType.Field(i).Name] {
			continue
		}
		field := inputValue.Field(i)
		switch inputValue.Field(i).Kind() {
		case reflect.Struct:
			addressField := field.Addr()
			// recurse if have struct nested on the same field name
			StructFields(addressField.Interface(), bypassFiled)
			// if slice changed then return slice changed
		case reflect.Slice:
			// TODO: inprove this package
			changedField[inputType.Field(i).Name] = [2]interface{}{
				nil, field.Interface(),
			}
			// if array changed detection then return array changed
		case reflect.Array:
			// TODO: inprove this package
			changedField[inputType.Field(i).Name] = [2]interface{}{
				nil, field.Interface(),
			} // default case field of the struct have a type default by golang
		default:
			changedField[inputType.Field(i).Name] = [2]interface{}{
				nil, field.Interface(),
			}
		}
	}
	return
}
