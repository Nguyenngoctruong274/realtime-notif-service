package tests

import (
	"github.com/stretchr/testify/require"
	"testing"
	"yes4all/ads-noti-api/pkg/utils/array"
	"yes4all/ads-noti-api/services/ads-noti/model/entity"
)

func TestMergeMapArrays(t *testing.T) {
	t.Run("test it should return when merge success", func(t *testing.T) {
		// arrange
		map1 := map[string][]entity.Keywords{
			"adgroup1": {
				{
					ID:    1,
					AwsID: "keyword_1",
					Bid:   1,
				},
			},
			"adgroup2": {
				{
					ID:    2,
					AwsID: "keyword_2",
					Bid:   1,
				},
			},
		}

		map2 := map[string][]entity.Keywords{
			"adgroup1": {
				{
					ID:    3,
					AwsID: "keyword_3",
					Bid:   1,
				},
			},
			"adgroup3": {
				{
					ID:    4,
					AwsID: "keyword_4",
					Bid:   1,
				},
			},
		}

		// act
		result := array.MergeMapArrays(map1, map2)

		// assert
		require.NotEmpty(t, result)
		require.Equal(t, map[string][]entity.Keywords{
			"adgroup1": {
				{
					ID:    1,
					AwsID: "keyword_1",
					Bid:   1,
				},
				{
					ID:    3,
					AwsID: "keyword_3",
					Bid:   1,
				},
			},
			"adgroup2": {
				{
					ID:    2,
					AwsID: "keyword_2",
					Bid:   1,
				},
			},
			"adgroup3": {
				{
					ID:    4,
					AwsID: "keyword_4",
					Bid:   1,
				},
			},
		}, result)
	})
}
