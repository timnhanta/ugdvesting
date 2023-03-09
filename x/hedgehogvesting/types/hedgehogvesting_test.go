package types_test

import (
	"testing"
	"time"

	"github.com/timnhanta/ugdvesting/x/hedgehogvesting/types"
)

func TestGetUnvestedAmount(t *testing.T) {
	// unvested should be 0, if vesting not started
	timeNow := time.Now()
	timeNow = timeNow.In(time.FixedZone("CET", 0))
	amount := "1000" + types.Denom
	duration := "P10M"
	part := int64(10)
	timeStart := timeNow.Add(10 * time.Minute)
	formattedTimeStart := timeStart.UTC().Format("2006-01-02T15:04:05.999999Z")

	vesting := types.Vesting{
		Amount:   amount,
		Start:    formattedTimeStart,
		Duration: duration,
		Parts:    part,
	}
	unvested := types.GetUnvestedAmount(vesting)

	expected := 0.0
	if unvested != expected {
		t.Errorf("unvested = %v, expected = %v", unvested, expected)
	}

	// unvested should be 0, if vesting is done
	timeNow = time.Now()
	timeNow = timeNow.In(time.FixedZone("CET", 0))
	timeStart = timeNow.Add(-11 * time.Minute)
	formattedTimeStart = timeStart.UTC().Format("2006-01-02T15:04:05.999999Z")

	vesting = types.Vesting{
		Amount:   amount,
		Start:    formattedTimeStart,
		Duration: duration,
		Parts:    part,
	}
	unvested = types.GetUnvestedAmount(vesting)

	expected = 0.0
	if unvested != expected {
		t.Errorf("unvested = %v, expected = %v", unvested, expected)
	}

	// unvested should be 600, if vesting progress is 400ugd/1000ugd
	timeNow = time.Now()
	timeNow = timeNow.In(time.FixedZone("CET", 0))
	timeStart = timeNow.Add(-4 * time.Minute)
	formattedTimeStart = timeStart.UTC().Format("2006-01-02T15:04:05.999999Z")

	vesting = types.Vesting{
		Amount:   amount,
		Start:    formattedTimeStart,
		Duration: duration,
		Parts:    part,
	}
	unvested = types.GetUnvestedAmount(vesting)

	expected = 600.0
	if unvested != expected {
		t.Errorf("unvested = %v, expected = %v", unvested, expected)
	}
}
