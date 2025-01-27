package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/vars"
)

func StructToMap(in any, snakeKey bool) map[string]any {
	var newMap map[string]any
	data, _ := json.Marshal(in)
	json.Unmarshal(data, &newMap)

	for k, v := range newMap {
		if v == nil {
			delete(newMap, k)
		}
	}

	if snakeKey {
		for k, v := range newMap {
			newMap[strcase.ToSnake(k)] = v
			delete(newMap, k)
		}
	}

	return newMap
}

func InterfaceToMap(in any, snakeKey bool) map[string]any {

	newMap := make(map[string]any)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			val := v.MapIndex(k)
			if snakeKey {
				newMap[strcase.ToSnake(k.String())] = val.Interface()
			} else {
				newMap[k.String()] = val.Interface()
			}

		}
	}
	return newMap
}

func ExistingKeyInMap(in any, key string) bool {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			if key == k.String() {
				return true
			}
		}
	}
	return false
}

func ExistingInArray(in any, f any) bool {

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Slice {
		vals := reflect.ValueOf(in)
		for i := 0; i < vals.Len(); i++ {
			vf := reflect.ValueOf(f)
			if vals.Index(i).Interface() == vf.Interface() {
				return true
			}
		}
	}
	return false
}

func ExistingInArrayIndex(in any, f any) int {

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Slice {
		vals := reflect.ValueOf(in)
		for i := 0; i < vals.Len(); i++ {
			vf := reflect.ValueOf(f)
			if vals.Index(i).Interface() == vf.Interface() {
				return i
			}
		}
	}
	return -1
}

func Transcode(in, out any) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

func FetchBingImage() (img string, copyright string, err error) {
	randomInt := rand.Intn(5)

	url := fmt.Sprintf("https://www.bing.com/HPImageArchive.aspx?format=xml&idx=%d&n=1", randomInt)

	var r []byte
	r, err = HttpRequestWithBytesResponse(types.HttpRequestType{
		Method:  "GET",
		FullURL: url,
	})
	if err != nil {
		return
	}

	var bingResp types.BingImageAPIResponse
	err = xml.NewDecoder(bytes.NewBuffer(r)).Decode(&bingResp)
	if err != nil {
		return
	}

	image, ok := vars.BGImageCache.Get(bingResp.Image.URLBase)
	if !ok {
		imageUrl := fmt.Sprintf("https://bing.com%s_1920x1080.jpg", bingResp.Image.URLBase)

		var respImg []byte
		respImg, err = HttpRequestWithBytesResponse(types.HttpRequestType{
			Method:  "GET",
			FullURL: imageUrl,
		})
		if err != nil {
			return
		}
		vars.BGImageCache.SetWithTTL(bingResp.Image.URLBase, base64.StdEncoding.EncodeToString(respImg), 1, time.Hour*24)
		image = base64.StdEncoding.EncodeToString(respImg)
	}

	img = image.(string)
	copyright = bingResp.Image.Copyright

	return
}

func GetDifferenceSlice[T any](original, compare T) (added T, deleted T, err error) {
	vOriginal := reflect.ValueOf(original)

	vAdded := reflect.New(reflect.TypeOf(added)).Elem()
	vDeleted := reflect.New(reflect.TypeOf(deleted)).Elem()

	if vOriginal.Kind() != reflect.Slice {
		err = errors.New("original must be slice")
		return
	}

	vCompare := reflect.ValueOf(compare)
	if vCompare.Kind() != reflect.Slice {
		err = errors.New("compare must be slice")
		return
	}

	valCompare := reflect.ValueOf(compare)
	for i := 0; i < valCompare.Len(); i++ {
		e := valCompare.Index(i).Interface()

		if found, _ := findExistSliceElement(original, e); !found {
			vAdded = reflect.Append(vAdded, reflect.ValueOf(e))
		}
	}

	valOriginal := reflect.ValueOf(original)
	for i := 0; i < valOriginal.Len(); i++ {
		e := valOriginal.Index(i).Interface()
		if found, _ := findExistSliceElement(compare, e); !found {
			vDeleted = reflect.Append(vDeleted, reflect.ValueOf(e))
		}

	}

	added = vAdded.Interface().(T)
	deleted = vDeleted.Interface().(T)

	return
}

func findExistSliceElement(source, element any) (found bool, err error) {
	val := reflect.ValueOf(source)
	if val.Kind() != reflect.Slice {
		err = errors.New("source must be slice")
		return
	}

	valSource := reflect.ValueOf(source)
	for i := 0; i < valSource.Len(); i++ {
		e := valSource.Index(i).Interface()

		if e == element {
			found = true
			return
		}
	}

	return
}

func DeleteSliceElement[T any](source []T, element T) (result []T, err error) {
	index := ExistingInArrayIndex(source, element)
	if index == -1 {
		err = errors.New("not found element in slice")
		return
	}

	result = source[:index]
	result = append(result, source[index+1:]...)
	return
}

func CheckDifference[T any](a, b T) (bool, error) {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	switch v1.Kind() {
	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			if !reflect.DeepEqual(v1.Field(i).Interface(), v2.Field(i).Interface()) {
				return true, nil
			}
		}

	case reflect.Array, reflect.Slice:
		if v1.Len() != v2.Len() {
			return true, nil
		}
		for i := 0; i < v1.Len(); i++ {
			if !reflect.DeepEqual(v1.Index(i).Interface(), v2.Index(i).Interface()) {
				return true, nil
			}
		}

	case reflect.Map:
		if v1.Len() != v2.Len() {
			return true, nil
		}
		for _, key := range v1.MapKeys() {
			if !v2.MapIndex(key).IsValid() || !reflect.DeepEqual(v1.MapIndex(key).Interface(), v2.MapIndex(key).Interface()) {
				return true, nil
			}
		}

	default:
		if !reflect.DeepEqual(a, b) {
			return true, nil
		}
	}

	return false, nil
}

func GetDifferenceMapResult[T any](a, b T) (map[string]any, error) {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	if v1.Type() != v2.Type() {
		return nil, fmt.Errorf("types do not match: %v vs %v", v1.Type(), v2.Type())
	}

	differences := make(map[string]any)

	switch v1.Kind() {
	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			fieldName := v1.Type().Field(i).Name
			val1 := v1.Field(i).Interface()
			val2 := v2.Field(i).Interface()
			if !reflect.DeepEqual(val1, val2) {
				differences[fieldName] = map[string]any{
					"old": val1,
					"new": val2,
				}
			}
		}

	case reflect.Array, reflect.Slice:
		maxLen := max(v1.Len(), v2.Len())
		for i := 0; i < maxLen; i++ {
			var val1, val2 any
			if i < v1.Len() {
				val1 = v1.Index(i).Interface()
			}
			if i < v2.Len() {
				val2 = v2.Index(i).Interface()
			}
			if !reflect.DeepEqual(val1, val2) {
				differences[fmt.Sprintf("%d", i)] = map[string]any{
					"old": val1,
					"new": val2,
				}
			}
		}

	case reflect.Map:
		keys := make(map[any]struct{})
		for _, key := range v1.MapKeys() {
			keys[key.Interface()] = struct{}{}
		}
		for _, key := range v2.MapKeys() {
			keys[key.Interface()] = struct{}{}
		}

		for key := range keys {
			val1 := v1.MapIndex(reflect.ValueOf(key)).Interface()
			val2 := v2.MapIndex(reflect.ValueOf(key)).Interface()
			if !reflect.DeepEqual(val1, val2) {
				differences[fmt.Sprintf("%v", key)] = map[string]any{
					"old": val1,
					"new": val2,
				}
			}
		}

	default:
		if !reflect.DeepEqual(a, b) {
			differences["value"] = map[string]any{
				"old": a,
				"new": b,
			}
		}
	}

	return differences, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
