package check_tool

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

var functions MapFunc

var TAG = "chk"

var ErrorKCHECK = errors.New("error unexpected")

func init() {
	functions = make(MapFunc)
	register()
}
func register() {
	functions["nonil"] = noNilFunc
	functions["nosp"] = noSpacesStartAndEnd
	functions["sword"] = sword
	functions["txt"] = textFunc
	functions["email"] = emailFunc
	functions["num"] = numberFunc
	functions["decimal"] = decimalFunc
	functions["len"] = lengthFunc
	functions["max"] = maxLengthFunc
	functions["min"] = minLengthFunc
	functions["rgx"] = regularExpression
}

// Fields SkipFields list of fields that were not taken into account when performing the verification
type Fields []string

type TagParamExtractor interface {
	GetTagValue(fieldName string) (value string, ok bool)
}

type paramExtractor struct {
	reflectValue reflect.Value
	reflectType  reflect.Type
}

func (ex *paramExtractor) GetTagValue(fieldName string) (value string, ok bool) {
	rsf, found := ex.reflectType.FieldByName(fieldName)
	if !found {
		return
	}

	if rsf.Type.Kind() == reflect.String {
		value = rsf.Tag.Get(TAG)
		ok = true
	}

	return
}

// AddFunc allows you to register a new custom function, which will be associated with the indicated tag key, if the
// tag key already exists, it will be replaced, for example, if the tag key `num` is used, it replaces the existing one.
// ValidFunc receives as its first parameter an object with the data of the field to verify, it includes the name and
// the value, and as its second parameter it receives the value after the `=` of the tag key, for example, the tag is
// `len` and it receives a value `len=10` the 10 is sent as the second parameter in string format.
// Note: the registration of new functions mustn't be done at runtime, this could generate access problems
// for the goroutines
func AddFunc(tagKey string, f ValidFunc) {
	functions[tagKey] = f
}

func (of *Fields) isContain(field string) bool {
	for _, v := range *of {
		if v == field {
			return true
		}
	}

	return false
}

func Valid(i interface{}) error {
	return ValidWithOmit(i, Fields{})
}

func ValidWithSelect(i interface{}, selected Fields) error {
	return valid(i, selected, false)
}

func ValidWithOmit(i interface{}, skips Fields) error {
	return valid(i, skips, true)
}

func reflectValueAndType(i interface{}) (*reflect.Value, *reflect.Type, error) {
	var rValue reflect.Value
	var rType reflect.Type = reflect.TypeOf(i)
	if rType == nil {
		log.Println("ERROR: nil value was received")
		return nil, nil, ErrorKCHECK
	}

	switch rType.Kind() {
	case reflect.Struct:
		rValue = reflect.ValueOf(i)
	case reflect.Ptr:
		if rType.Elem().Kind() == reflect.Struct {
			rValue = reflect.ValueOf(i).Elem()
			rType = rType.Elem()
		} else {
			log.Printf("ERROR: a structure was type expected, invalid type `%v`\n", rType)
			return nil, nil, ErrorKCHECK
		}
	}

	return &rValue, &rType, nil
}

func valid(i interface{}, fields Fields, isOmit bool) error {
	reflectValue, reflectType, err := reflectValueAndType(i)
	if err != nil {
		return err
	}

	for i := 0; i < (*reflectType).NumField(); i++ {
		rField := (*reflectType).Field(i)
		rValue := reflectValue.Field(i)
		if rField.Type.Kind() == reflect.String {
			tagValues := rField.Tag.Get(TAG)

			if isOmit {
				if tagValues == "" || fields.isContain(rField.Name) {
					continue
				}
			} else {
				if !fields.isContain(rField.Name) {
					continue
				}
			}

			atom := Atom{Name: SplitCamelCase(rField.Name), Value: rValue.String()}
			if err = ValidTarget(tagValues, atom); err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidTarget(tags string, atom Atom) error {
	tags = StandardSpace(tags)
	keys := strings.Split(tags, " ")
	for _, key := range keys {
		if f, ok := functions[key]; ok {
			if err := f(atom, ""); err != nil {
				return err
			}

		} else {
			valid, fKey, keyValues := SplitKeyValue(key)
			if valid {
				if function, okk := functions[fKey]; okk {
					if err := function(atom, keyValues); err != nil {
						return err
					}

				} else {
					log.Printf("ERROR: tag value `%s` invalid in `%s` field\n", key, atom.Name)
					return ErrorKCHECK
				}

			} else {
				log.Printf("ERROR: tag value `%s` invalid in `%s` field\n", key, atom.Name)
				return ErrorKCHECK
			}
		}
	}

	return nil
}

func BuildTagParamExtractor(entry interface{}) (TagParamExtractor, error) {
	reflectValue, reflectType, err := reflectValueAndType(entry)
	if err != nil {
		return nil, err
	}

	return &paramExtractor{reflectValue: *reflectValue, reflectType: *reflectType}, nil
}
