package protocol

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*
@Time : 		2023/3/28 9:46
@Author : 		xumx01
@File :			fast_test.go
@Software:		GOLAND
*/

func overwrite(ori map[string]interface{}, newMap map[string]interface{}) {
	for k, v := range newMap {
		i := ori[k]
		if i == reflect.Zero(reflect.TypeOf(v)).Interface() {
			ori[k] = v
			fmt.Printf("%s: %v overwrite to %v\n", k, i, v)
		} else {
			fmt.Printf("%s: %v don't need to overwrite\n", k, i)
		}
	}
}

func TestOverwrite(t *testing.T) {
	ori := map[string]interface{}{
		"a": "",
		"b": 0,
		"c": 0.0,
	}
	newMap := map[string]interface{}{
		"a": "1",
		"b": 2,
		"c": 3.0,
	}
	Convey("change values when original values are zero val", t, func() {
		overwrite(ori, newMap)
		Convey("it should be overwritten", func() {
			So(ori["a"], ShouldEqual, "1")
			So(ori["b"], ShouldEqual, 2)
			So(ori["c"], ShouldEqual, 3.0)
		})
	})

	ori = map[string]interface{}{
		"a": "4",
		"b": 5,
		"c": 6.0,
	}
	Convey("change values when original values are not zero val", t, func() {
		overwrite(ori, newMap)
		Convey("it should not be overwritten", func() {
			So(ori["a"], ShouldEqual, "4")
			So(ori["b"], ShouldEqual, 5)
			So(ori["c"], ShouldEqual, 6.0)
		})
	})
}
