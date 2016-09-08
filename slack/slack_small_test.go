//
// +build small

/*
http://www.apache.org/licenses/LICENSE-2.0.txt



Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package slack

import (
	"fmt"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSlackPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create Slack Collector", t, func() {
		bCol := NewSlackCollector()
		Convey("So bCol should not be nil", func() {
			So(bCol, ShouldNotBeNil)
		})
		Convey("So bCol should be of Slack type", func() {
			So(bCol, ShouldHaveSameTypeAs, &Slack{})
		})
		Convey("bCol.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := bCol.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			testConfig := make(map[string]ctypes.ConfigValue)
			testConfig["token"] = ctypes.ConfigValueStr{Value: "example-token"}
			testConfig["max_element"] = ctypes.ConfigValueInt{Value: 42}
			cfg, errs := configPolicy.Get([]string{vendor, name}).Process(testConfig)
			Convey("So config policy should process testConfig and return a config", func() {
				So(cfg, ShouldNotBeNil)
			})
			testConfig["max_element"] = ctypes.ConfigValueStr{Value: "this-is-not-an-integer"}
			cfg, errs = configPolicy.Get([]string{vendor, name}).Process(testConfig)
			Convey("So config policy should not return a config after processing invalid testConfig", func() {
				So(cfg, ShouldBeNil)
			})
			Convey("So testConfig processing should return errors", func() {
				So(errs.HasErrors(), ShouldBeTrue)
			})
		})
	})
}

type mystruct struct {
	MyString     string
	MySubStruct  mysubstruct
	MyFloatArray []float64
	MyPointer    *int
	MyMapPointer map[string]*mystruct
}

type mysubstruct struct {
	MyString     string
	MyFloatArray []float64
	MyPointer    *int
	MyMapPointer map[string]*mystruct
}

var (
	plop = 42
)

func newMystructPointer(i string) *mystruct {
	return &mystruct{
		MyString:     "mystructpointer" + i,
		MySubStruct:  newSubStruct("frompointer" + i),
		MyFloatArray: make([]float64, 3),
		MyPointer:    &plop,
		MyMapPointer: make(map[string]*mystruct),
	}
}

func newSubStruct(i string) mysubstruct {
	return mysubstruct{
		MyString:     "mysubstruct" + i,
		MyFloatArray: make([]float64, 3),
		MyPointer:    &plop,
		MyMapPointer: make(map[string]*mystruct),
	}
}

func newStruct() mystruct {
	m := make(map[string]*mystruct)
	for i := 1; i < 4; i++ {
		istr := fmt.Sprintf("%v", i)
		m[istr] = newMystructPointer(istr)
	}
	return mystruct{
		MyString:     "mystruct",
		MySubStruct:  newSubStruct("0"),
		MyFloatArray: make([]float64, 3),
		MyPointer:    &plop,
		MyMapPointer: m,
	}
}

var (
	testMap = map[string]string{"test-mysubstruct-mystring": "mysubstruct0", "test-mymappointer-2-mysubstruct-mystring": "mysubstructfrompointer2", "test-mymappointer-2-mysubstruct-myfloatarray-0": "0", "test-mymappointer-2-myfloatarray-0": "0", "test-mymappointer-3-mysubstruct-mystring": "mysubstructfrompointer3", "test-mymappointer-3-myfloatarray-1": "0", "test-mymappointer-1-mysubstruct-mystring": "mysubstructfrompointer1", "test-mystring": "mystruct", "test-mymappointer-2-mystring": "mystructpointer2", "test-mymappointer-2-myfloatarray-1": "0", "test-mymappointer-1-mystring": "mystructpointer1", "test-mymappointer-2-mysubstruct-mypointer": "42", "test-mymappointer-3-myfloatarray-0": "0", "test-mymappointer-3-myfloatarray-2": "0", "test-mymappointer-1-mysubstruct-myfloatarray-0": "0", "test-mymappointer-1-mysubstruct-mypointer": "42", "test-mysubstruct-myfloatarray-1": "0", "test-mysubstruct-mypointer": "42", "test-myfloatarray-0": "0", "test-myfloatarray-2": "0", "test-mymappointer-2-mypointer": "42", "test-mymappointer-1-mysubstruct-myfloatarray-1": "0", "test-mymappointer-1-mysubstruct-myfloatarray-2": "0", "test-mymappointer-1-myfloatarray-2": "0", "test-mysubstruct-myfloatarray-0": "0", "test-mymappointer-2-mysubstruct-myfloatarray-1": "0", "test-mymappointer-2-mysubstruct-myfloatarray-2": "0", "test-mymappointer-3-mypointer": "42", "test-mymappointer-1-mypointer": "42", "test-mysubstruct-myfloatarray-2": "0", "test-myfloatarray-1": "0", "test-mymappointer-2-myfloatarray-2": "0", "test-mymappointer-3-mystring": "mystructpointer3", "test-mymappointer-3-mysubstruct-myfloatarray-1": "0", "test-mymappointer-3-mysubstruct-myfloatarray-2": "0", "test-mymappointer-1-myfloatarray-1": "0", "test-mypointer": "42", "test-mymappointer-3-mysubstruct-myfloatarray-0": "0", "test-mymappointer-3-mysubstruct-mypointer": "42", "test-mymappointer-1-myfloatarray-0": "0"}
)

func TestGetTags(t *testing.T) {
	Convey("GetTags should return all the elements of any data type in a map", t, func() {
		myStr := newStruct()
		m, err := GetTags(myStr, "test", "-")
		Convey("So GetTags should not return an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("GetTags returns a non-empty map", func() {
			So(m, ShouldNotBeEmpty)
		})

		Convey("GetTags elements match testMap", func() {
			everythingIsFine := true
			for key, value := range m {
				elt, ok := testMap[key]
				if ok != true || elt != value {
					fmt.Printf("KEY1 %v: %v, KEY2 %v: %v", key, value, ok, elt)
					everythingIsFine = false
					break
				}
			}
			So(everythingIsFine, ShouldBeTrue)
		})
	})

}
