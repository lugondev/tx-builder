package builder

import (
	"math"
	"testing"
)

func TestCalc(t *testing.T) {
	f := 1.1
	t.Log(int64(f))
	t.Log(int64(math.Ceil(f)))
	t.Log(int64(math.Round(f)))
}
