package utils

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/leekchan/accounting"
	"github.com/mitchellh/mapstructure"
)

var tableConvertCases = map[rune]rune{
	'ấ':  'a',
	'ầ':  'a',
	'ẩ':  'a',
	'ẫ':  'a',
	'ậ':  'a',
	'Ấ':  'a',
	'Ầ':  'a',
	'Ẩ':  'a',
	'Ẫ':  'a',
	'Ậ':  'a',
	'ắ':  'a',
	'ằ':  'a',
	'ẳ':  'a',
	'ẵ':  'a',
	'ặ':  'a',
	'Ắ':  'a',
	'Ằ':  'a',
	'Ẳ':  'a',
	'Ẵ':  'a',
	'Ặ':  'a',
	'á':  'a',
	'à':  'a',
	'ả':  'a',
	'ã':  'a',
	'ạ':  'a',
	'â':  'a',
	'ă':  'a',
	'Á':  'a',
	'À':  'a',
	'Ả':  'a',
	'Ã':  'a',
	'Ạ':  'a',
	'Â':  'a',
	'Ă':  'a',
	'ế':  'e',
	'ề':  'e',
	'ể':  'e',
	'ễ':  'e',
	'ệ':  'e',
	'Ế':  'e',
	'Ề':  'e',
	'Ể':  'e',
	'Ễ':  'e',
	'Ệ':  'e',
	'é':  'e',
	'è':  'e',
	'ẻ':  'e',
	'ẽ':  'e',
	'ẹ':  'e',
	'ê':  'e',
	'É':  'e',
	'È':  'e',
	'Ẻ':  'e',
	'Ẽ':  'e',
	'Ẹ':  'e',
	'Ê':  'e',
	'í':  'i',
	'ì':  'i',
	'ỉ':  'i',
	'ĩ':  'i',
	'ị':  'i',
	'Í':  'i',
	'Ì':  'i',
	'Ỉ':  'i',
	'Ĩ':  'i',
	'Ị':  'i',
	'ố':  'o',
	'ồ':  'o',
	'ổ':  'o',
	'ỗ':  'o',
	'ộ':  'o',
	'Ố':  'o',
	'Ồ':  'o',
	'Ổ':  'o',
	'Ộ':  'o',
	'ớ':  'o',
	'ờ':  'o',
	'ở':  'o',
	'ỡ':  'o',
	'ợ':  'o',
	'Ớ':  'o',
	'Ờ':  'o',
	'Ở':  'o',
	'Ỡ':  'o',
	'Ợ':  'o',
	'ó':  'o',
	'ò':  'o',
	'ỏ':  'o',
	'õ':  'o',
	'ọ':  'o',
	'ô':  'o',
	'ơ':  'o',
	'Ó':  'o',
	'Ò':  'o',
	'Ỏ':  'o',
	'Õ':  'o',
	'Ọ':  'o',
	'Ô':  'o',
	'Ơ':  'o',
	'ứ':  'u',
	'ừ':  'u',
	'ử':  'u',
	'ữ':  'u',
	'ự':  'u',
	'Ứ':  'u',
	'Ừ':  'u',
	'Ử':  'u',
	'Ữ':  'u',
	'Ự':  'u',
	'ú':  'u',
	'ù':  'u',
	'ủ':  'u',
	'ũ':  'u',
	'ụ':  'u',
	'ư':  'u',
	'Ú':  'u',
	'Ù':  'u',
	'Ủ':  'u',
	'Ũ':  'u',
	'Ụ':  'u',
	'Ư':  'u',
	'ý':  'y',
	'ỳ':  'y',
	'ỷ':  'y',
	'ỹ':  'y',
	'ỵ':  'y',
	'Ý':  'y',
	'Ỳ':  'y',
	'Ỷ':  'y',
	'Ỹ':  'y',
	'Ỵ':  'y',
	'đ':  'd',
	'Đ':  'd',
	'&':  ' ',
	'<':  ' ',
	'>':  ' ',
	'-':  ' ',
	'ç':  'c',
	'~':  ' ',
	'¥':  'y',
	'©':  'c',
	'®':  'r',
	'\t': ' ',
	'_':  ' ',
	'+':  ' ',
	',':  ' ',
	'*':  ' ',
	'\\': ' ',
	'^':  ' ',
	'$':  ' ',
	'|':  ' ',
	'?':  ' ',
	'[':  ' ',
	']':  ' ',
	'{':  ' ',
	'}':  ' ',
	'(':  ' ',
	')':  ' ',
}

func RemoveAccents(s string) string {
	s1 := ""
	s = strings.ToLower(s)
	for _, c := range s {
		if v, ok := tableConvertCases[c]; ok {
			s1 += string(v)
		} else {
			s1 += string(c)
		}
	}
	return strings.TrimSpace(RemoveDuplicatedWhiteSpace(s1))
}

func RemoveDuplicatedWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	before := utf8.MaxRune
	for _, c := range str {
		if before == ' ' && c == ' ' {
			continue
		}
		b.WriteRune(c)
		before = c
	}
	return b.String()
}

// StringToInt return integer from string instantly, be careful when using it (no err validation)
func StringToInt(str string) int64 {
	result, _ := strconv.Atoi(str)
	return int64(result)
}

func StripSpecial(str string) string {
	re := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+`)
	return re.ReplaceAllString(str, "")
}

// nolint: errchkjson, nestif
func Inspect(data interface{}) map[string]interface{} {
	mapData := map[string]interface{}{}
	valueOf := reflect.ValueOf(data)
	if valueOf.Kind() != reflect.Interface && valueOf.Kind() != reflect.Ptr {
		mapData[valueOf.Kind().String()] = data
		return mapData
	}
	vals := valueOf.Elem()
	if vals.Kind() == reflect.Slice {
		jsonData, _ := json.Marshal(data)
		mapData[vals.Kind().String()] = string(jsonData)
		return mapData
	}
	for i := 0; i < vals.NumField(); i++ {
		valueField := vals.Field(i)
		typeField := vals.Type().Field(i)
		f := valueField.Interface()
		val := reflect.ValueOf(f)
		fieldName := typeField.Name
		if jsonTag := typeField.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			fieldName = jsonTag
		}

		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			mapData[fieldName] = strconv.FormatInt(val.Int(), 10)
		case reflect.Float32, reflect.Float64:
			mapData[fieldName] = strconv.FormatFloat(val.Float(), 'f', -1, 64)
		case reflect.String:
			mapData[fieldName] = val.String()
		case reflect.Slice, reflect.Struct, reflect.Array:
			jsonData, _ := json.Marshal(f)
			mapData[fieldName] = string(jsonData)
		default:
			if f == nil || (val.Kind() == reflect.Ptr && val.IsNil()) {
			} else {
				if val.Kind() == reflect.Ptr {
					if !val.IsNil() {
						jsonData, _ := json.Marshal(&f)
						mapData[fieldName] = string(jsonData)
					}
				} else {
					mapData[fieldName] = f
				}
			}
		}
	}
	return mapData
}

// nolint: nestif, forcetypeassert
func CustomDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		var err error
		if f.Kind() == reflect.String &&
			(t.Kind() == reflect.Slice || t.Kind() == reflect.Struct || t.Kind() == reflect.Ptr) {
			if len(data.(string)) < 1 {
				return nil, nil
			}
			if t.Kind() == reflect.Ptr {
				temp := t.Elem()
				if temp == reflect.TypeOf(time.Time{}) {
					switch f.Kind() {
					case reflect.String:
						return time.Parse(time.RFC3339, data.(string))
					case reflect.Float64:
						return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
					case reflect.Int64:
						return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
					default:
						return data, nil
					}
				}
			}

			if t == reflect.TypeOf(time.Time{}) {
				switch f.Kind() {
				case reflect.String:
					return time.Parse(time.RFC3339, data.(string))
				case reflect.Float64:
					return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
				case reflect.Int64:
					return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
				default:
					return data, nil
				}
			}
			var jsonMap interface{}
			err = json.Unmarshal([]byte(data.(string)), &jsonMap)
			if err != nil {
				return nil, err
			}
			return jsonMap, err
		}
		if f.Kind() == reflect.Map {
			if t.Kind() == reflect.Slice {
				_str, strOk := data.(map[string]string)
				_struct, struOk := data.(string)

				if strOk {
					if val, ok := _str[t.Kind().String()]; ok {
						var jsonMap []interface{}
						err := json.Unmarshal([]byte(val), &jsonMap)
						return jsonMap, err
					}
				} else if struOk {
					var jsonMap interface{}
					err := json.Unmarshal([]byte(_struct), &jsonMap)
					return jsonMap, err
				}
			}
			if t.Kind() != reflect.Struct {
				_mapValue, strOk := data.(map[string]string)
				if strOk {
					if val, ok := _mapValue[t.Kind().String()]; ok {
						var jsonMap interface{}
						err := json.Unmarshal([]byte(val), &jsonMap)
						return jsonMap, err
					}
				}
			}
		}
		return data, nil
	}
}

func ToFormatMoneyString(value interface{}, symbol string, decimal string) string {
	ac := accounting.NewAccounting(symbol, 0, decimal, decimal, "%s%v", "%s(%v)", "%s0")
	return ac.FormatMoney(value)
}

func GenerateSenPaySignature(args ...string) string {
	var str string
	if len(args) == 1 {
		str = args[0]
	} else {
		str = strings.Join(args, "|")
	}

	return str
}

func Capital(s string) (res string) {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	res = string(r)
	return
}
