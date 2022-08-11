package util

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
	"github.com/haltingstate/secp256k1-go"
)

type NamePicker interface {
	TagName(st reflect.StructTag) string
}

type JSONTagPicker struct {
}

// TagName return tag name of json, if not confirm return empty string
func (jp *JSONTagPicker) TagName(st reflect.StructTag) string {
	jsonTag := st.Get("json")
	tagOpts := strings.Split(jsonTag, ",")
	if len(tagOpts) < 1 {
		return ""
	}
	jsonFieldName := strings.Trim(tagOpts[0], " ")
	// filter
	switch jsonFieldName {
	case "", "-":
		return ""
	}
	return jsonFieldName
}

func paramsMap(i interface{}) map[string]string {
	ret := make(map[string]string)
	picker := JSONTagPicker{}

	paramTypes := reflect.TypeOf(i)
	paramValues := reflect.ValueOf(i)

	switch paramTypes.Kind() {
	case reflect.Ptr:
		paramTypes = paramTypes.Elem()
		paramValues = paramValues.Elem()
	}

	for i := 0; i < paramTypes.NumField(); i++ {
		field := paramTypes.Field(i)
		if !field.IsExported() {
			continue
		}
		//pick tag name
		fieldName := picker.TagName(field.Tag)
		if fieldName == "" {
			continue
		}
		//set values
		fieldValue := paramValues.Field(i)
		switch fieldValue.Kind() {
		case reflect.Ptr, reflect.Array, reflect.Map, reflect.Struct:
			continue
		}
		ret[fieldName] = toString(fieldValue)
	}
	return ret
}

func toString(fieldValue reflect.Value) string {
	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String()
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(fieldValue.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(fieldValue.Uint(), 10)
	default:
		return fmt.Sprintf("%v", fieldValue.Interface())
	}
}

// summary
func summaryURLEncode(i interface{}) ([]byte, error) {
	params := paramsMap(i)
	delete(params, "sign")

	u := url.Values{}
	for k, v := range params {
		u.Set(k, v)
	}
	//按字典升序排序并URL编码
	secretStr := u.Encode()

	return []byte(secretStr), nil
}

func summaryNoEncode(i interface{}) ([]byte, error) {
	params := paramsMap(i)
	delete(params, "sign")

	u := url.Values{}
	for k, v := range params {
		u.Set(k, v)
	}

	var buf strings.Builder
	keys := make([]string, 0, len(u))
	for k := range u {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := u[k]
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}

	//按字典升序排序
	secretStr := buf.String()

	return []byte(secretStr), nil
}

// GetSummary 得到摘要
func GetSummary(i interface{}) ([]byte, error) {
	return summaryNoEncode(i)
}

func Secp256k1VerifyChain33(msg, sig, pubKey []byte) (b bool) {
	c, err := crypto.Load(types.GetSignName("", types.SECP256K1), -1)
	if err != nil {
		return false
	}

	pub, err := c.PubKeyFromBytes(pubKey)
	if err != nil {
		return false
	}

	signature, err := c.SignatureFromBytes(sig)
	if err != nil {
		return false
	}

	return pub.VerifyBytes(msg, signature)
}

func Secp256k1Verify(msg, sig, pubKey []byte) (b bool) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			log15.Error("secp256k1 VerifySignature failed", "err", err)
			b = false
		}
	}()
	return 1 == secp256k1.VerifySignature(msg, sig, pubKey)
}
