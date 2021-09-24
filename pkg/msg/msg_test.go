package msg

import "testing"

// TestMsg does not really test anything; it is just a useful tool
// (with `go test -v ./...`) to visually see if the color codes are working
func TestMsg(t *testing.T) {

	initMap()

	Debug("This is a test of the debug command")

	for i, k := range ColorMapKeys {
		Msg(ColorMap[k], "%d: This is a sample %s message\n", i, k)
	}

}
