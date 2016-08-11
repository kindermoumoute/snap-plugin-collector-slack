package slack

import (
	"fmt"
	"reflect"
	"strings"
)

// GetTags return all the fields of an object in a single map of tags.
// - rootprefix is the name of the object.
// - sep is the string separator between 2 layers of deepness (e.g. "test-myfloatarray-1").
func GetTags(object interface{}, rootprefix, sep string) (map[string]string, error) {
	tags := make(map[string]string)
	return tags, GetTagsRec(reflect.ValueOf(object), tags, rootprefix, sep)
}

func GetTagsRec(obj reflect.Value, tags map[string]string, prefix, sep string) error {
	var err error
	switch obj.Kind() {
	case reflect.Slice:
		for i := 0; i < obj.Len(); i++ {
			nextPrefix := fmt.Sprintf("%s%s%d", prefix, sep, i)
			err = GetTagsRec(obj.Index(i), tags, nextPrefix, sep)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		keys := obj.MapKeys()
		for _, key := range keys {
			nextPrefix := fmt.Sprintf("%s%s%s", prefix, sep, key)
			err = GetTagsRec(obj.MapIndex(key), tags, nextPrefix, sep)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		fieldsCount := obj.NumField()
		objType := obj.Type()
		for i := 0; i < fieldsCount; i++ {
			field := obj.Field(i)
			objTypeField := objType.Field(i)
			key := objTypeField.Name
			if objTypeField.PkgPath == "" {
				if !field.IsValid() {
					return fmt.Errorf("No such field: %s in obj", key)
				}

				err = GetTagsRec(field, tags, prefix+sep+key, sep)
				if err != nil {
					return err
				}
			}
		}
	case reflect.Ptr:
		v := reflect.Indirect(obj)
		if !v.IsValid() || !v.CanInterface() {
			// If reflect can't call Interface() on v, we can't go deeper even if
			// len(ns) > 1. Therefore, we should just return nil here.
			return nil
		}
		err = GetTagsRec(v, tags, prefix, sep)

	default:
		tmpVal := fmt.Sprintf("%v", obj)
		if tmpVal != "" {
			tags[strings.ToLower(prefix)] = tmpVal
		}
		return nil
	}
	return err
}
