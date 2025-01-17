package types_test

import (
	"clousel/infra/repo/types"
	"testing"
	"time"
)

// 2024-12-23 18:27:12

func TestTimeString(t *testing.T) {
	var st types.TimeString
	st.SetStr("2025-01-05 16:27:05")
	tm := st.Time()
	tm = tm.Add(time.Hour)
	st.SetTime(tm)
	if st.Str() != "2025-01-05 17:27:05" {
		t.Errorf("Fail to add 1 hour to time container %s != %s", tm.Format(time.DateTime), st.Str())
	}
}
