package timeutils

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

/*
type timeUtilsSuite struct {
	suite.Suite
	goMockCtrl *gomock.Controller
}

func TestTimeUtilsSuite(t *testing.T) {
	suite.Run(t, new(timeUtilsSuite))
}

func (t *timeUtilsSuite) SetupTest() {
	t.goMockCtrl = gomock.NewController(t.T())
}

func (t *timeUtilsSuite) TearDownTest() {
	t.goMockCtrl.Finish()
}

func (t *timeUtilsSuite) TestConvertStringToTime() {
	testcases := []struct {
		name     string
		input    string
		expected time.Time
	}{}
	for _, tc := range testcases {
		t.Run(tc.name, func() {
			ctxMock := gomock.Any()
			actual:=t.
		}
	}
}*/

func TestConvertStringToTime(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{name: "from pg database", input: "2022-03-30 16:54:48", expected: SetDateTime(2022, 3, 30, 16, 54, 48)},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			actual := ConvertStringToTime(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGenLoyaltyPeriodByDate(t *testing.T) {
	testcases := []struct {
		name     string
		input    time.Time
		expected func() (StartAt, EndAt time.Time)
	}{
		{
			name:  "2022-03-15",
			input: SetDateTime(2022, 3, 15, 16, 54, 48),
			expected: func() (StartAt, EndAt time.Time) {
				StartAt = SetDateTime(2022, 3, 01, 0, 0, 0)
				EndAt = EndOfDayGMT07(SetDateTime(2022, 3, 31, 23, 0, 0))
				return
			},
		},
		{
			name:  "2022-03-01",
			input: SetDateTime(2022, 3, 01, 16, 54, 48),
			expected: func() (StartAt, EndAt time.Time) {
				StartAt = SetDateTime(2022, 3, 01, 0, 0, 0)
				EndAt = EndOfDayGMT07(SetDateTime(2022, 3, 31, 23, 0, 0))
				return
			},
		},
		{
			name:  "2022-03-31",
			input: SetDateTime(2022, 3, 31, 16, 54, 48),
			expected: func() (StartAt, EndAt time.Time) {
				StartAt = SetDateTime(2022, 3, 01, 0, 0, 0)
				EndAt = EndOfDayGMT07(SetDateTime(2022, 3, 31, 23, 0, 0))
				return
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			startAtExpected, endAtExpected := tc.expected()
			startAtActual, endAtActual := GenLoyaltyPeriodByDate(tc.input)
			assert.Equal(t, startAtExpected, startAtActual)
			assert.Equal(t, endAtExpected, endAtActual)
		})
	}
}
