package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
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

func InterfaceToMap(in interface{}, snakeKey bool) map[string]interface{} {

	newMap := make(map[string]interface{})
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

func ExistingKeyInMap(in interface{}, key string) bool {
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

func ExistingInArray(in interface{}, f interface{}) bool {

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

func Transcode(in, out interface{}) {
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
