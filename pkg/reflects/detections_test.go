package reflects

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type User struct {
	Name      string
	Locations Locations
	CreatedAt time.Time
}
type Locations struct {
	Longtitude int
	House      []House
	DieArrays  []int
}
type House struct {
	Dumb bool
}

func Test_StructChangeDetection(t *testing.T) {

	// struct with field is default of golang
	type input struct {
		name               string
		inputStruct        interface{}
		changedInputStruct interface{}
		bypassFiled        map[string]bool
	}

	type Expected struct {
		Result map[string][2]interface{}
	}

	testCases := []struct {
		Name   string
		Input  input
		Result Expected
	}{
		{
			Name: "Model with common fields",
			Input: input{
				name: "",
				inputStruct: &User{
					Name: "DongTQ",
					Locations: Locations{
						Longtitude: 0,
						House: []House{
							{
								Dumb: true,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				changedInputStruct: &User{
					Name: "DongTQ",
					Locations: Locations{
						Longtitude: 0,
						House: []House{
							{
								Dumb: true,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				bypassFiled: map[string]bool{
					"CreatedAt": true,
				},
			},
			Result: Expected{
				Result: map[string][2]interface{}{},
			},
		},
		{
			Name: "Change level 1",
			Input: input{
				name: "",
				inputStruct: &User{
					Name: "DongTQ 01",
					Locations: Locations{
						Longtitude: 0,
						House: []House{
							{
								Dumb: true,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				changedInputStruct: &User{
					Name: "DongTQ",
					Locations: Locations{
						Longtitude: 0,
						House: []House{
							{
								Dumb: true,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				bypassFiled: map[string]bool{
					"CreatedAt": true,
				},
			},
			Result: Expected{
				Result: map[string][2]interface{}{
					"Name": {
						"DongTQ O1", "DongTQ",
					},
				},
			},
		},
		{
			Name: "Change level 2",
			Input: input{
				name: "",
				inputStruct: &User{
					Name: "DongTQ",
					Locations: Locations{
						Longtitude: 0,
						House: []House{
							{
								Dumb: true,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				changedInputStruct: &User{
					Name: "DongTQ",
					Locations: Locations{
						Longtitude: 10,
						House: []House{
							{
								Dumb: true,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				bypassFiled: map[string]bool{
					"CreatedAt": true,
				},
			},
			Result: Expected{
				Result: map[string][2]interface{}{},
			},
		},
		{
			Name: "Change level 3",
			Input: input{
				name: "",
				inputStruct: &User{
					Name: "DongTQ",
					Locations: Locations{
						Longtitude: 0,
						House: []House{
							{
								Dumb: true,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				changedInputStruct: &User{
					Name: "DongTQ",
					Locations: Locations{
						Longtitude: 0,
						House: []House{
							{
								Dumb: false,
							},
						},
						DieArrays: []int{1, 2, 3, 4},
					},
				},
				bypassFiled: map[string]bool{
					"CreatedAt": true,
				},
			},
			Result: Expected{
				Result: map[string][2]interface{}{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			result := StructChangeDetection(tc.Input.inputStruct, tc.Input.changedInputStruct, tc.Input.bypassFiled)
			totalResult := reflect.DeepEqual(result, tc.Result.Result)
			fmt.Println(result)
			fmt.Println(tc.Result.Result)
			if !totalResult {
				t.Fail()
			}
			// arrTest := []int{}
			// len(arrTest)
		})
	}

}
