package log

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLevel_String(t *testing.T) {
	Convey("log level test", t, func() {
		lv := LevelDebug
		So(lv.String(), ShouldEqual, "debug")
	})
}
