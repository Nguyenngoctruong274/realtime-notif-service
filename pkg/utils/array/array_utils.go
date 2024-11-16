package array

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var (
	validReflectType = []reflect.Kind{reflect.Slice, reflect.Array}
)

// ExceptItems trả về mảng là các phần tử thuộc mảng "parents" nhưng không thuộc mảng "childs"
func ExceptItems(parents, childs interface{}) interface{} {
	pArr := reflect.ValueOf(parents)
	cArr := reflect.ValueOf(childs)

	if !ItemExists(validReflectType, pArr.Kind()) || !ItemExists(validReflectType, cArr.Kind()) {
		return nil
	}

	childMap := make(map[interface{}]bool)
	for i := 0; i < cArr.Len(); i++ {
		childMap[cArr.Index(i).Interface()] = true
	}

	result := make([]interface{}, 0)
	var val interface{}

	for i := 0; i < pArr.Len(); i++ {
		val = pArr.Index(i).Interface()
		if !childMap[val] {
			result = append(result, val)
		}
	}

	return result
}

func RemoveDuplicatesString(elements []string) []string {
	encountered := make(map[string]bool)
	result := make([]string, 0)
	// Create a map of all unique elements.
	for _, v := range elements {
		if _, ok := encountered[v]; !ok {
			result = append(result, v)
			encountered[v] = true
		}
	}
	return result
}

func GetItemsByLen(slice interface{}, len int) map[int]interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return map[int]interface{}{}
	}
	lenOfSlice := s.Len()
	if lenOfSlice < 1 {
		return map[int]interface{}{}
	}
	if len > lenOfSlice {
		len = lenOfSlice
	}
	result := make(map[int]interface{}, len)

	for i := 0; i < len; i++ {
		result[i] = s.Index(i).Interface()
	}

	return result
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// Intersection trả về xem liệu 2 mảng có giao nhau hay không, nếu có thì các phần tử nào giao nhau
func Intersection(firstArr, secondArr interface{}) (isConflict bool, values []interface{}) {
	fArr := reflect.ValueOf(firstArr)
	sArr := reflect.ValueOf(secondArr)

	if fArr.Kind() != reflect.Slice || sArr.Kind() != reflect.Slice {
		return false, nil
	}

	countMap := make(map[interface{}]bool)
	for i := 0; i < fArr.Len(); i++ {
		countMap[fArr.Index(i).Interface()] = true
	}

	for i := 0; i < sArr.Len(); i++ {
		if countMap[sArr.Index(i).Interface()] {
			isConflict = true
			values = append(values, sArr.Index(i).Interface())
		}
	}

	return
}

func ConvertArrayToString(list interface{}, sep string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(list), " ", sep, -1), "[]")
}

func IsExist(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}

	return false
}

func FindIndex(arr []string, value string) int {
	for index, item := range arr {
		if item == value {
			return index
		}
	}

	return -1
}

func ConvertArrayToStringDB(list interface{}) string {
	items, _ := json.Marshal(list)
	return strings.Trim(strings.Replace(string(items), "\"", "'", -1), "[]")
}

func MergeMapArrays[K comparable, V any](m1 map[K][]V, m2 map[K][]V) map[K][]V {
	merged := make(map[K][]V)
	for key, value := range m1 {
		merged[key] = append(merged[key], value...)
	}
	for key, value := range m2 {
		merged[key] = append(merged[key], value...)
	}
	return merged
}
