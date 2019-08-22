package utils

import (
	"testing"
)

type Person struct {
	Name    string
	Profile *Profile
}

type Profile struct {
	Age  int
	City string
}
type Job struct {
	Id   int
	Desc []Profile
}
type Job2 struct {
	Id   int
	Desc [2]Profile
}

type Job3 struct {
	Id   int
	Desc map[string]Profile
}

func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"a simple string input",
			"stephenx1",
			[]string{"stephenx1"},
		},
		{
			"slice of string input",
			[]string{"stephen2", "xiong2"},
			[]string{"stephen2", "xiong2"},
		},
		{
			"struct with one string field",
			struct {
				S string
				I int
			}{"stephenx3", 30},
			[]string{"stephenx3"},
		},
		{
			"struct with nested field",
			Person{
				"stephen4",
				&Profile{30, "london4"},
			},
			[]string{"stephen4", "london4"},
		},
		{
			"struct with pointer field",
			Person{
				"xiong5",
				&Profile{30, "london5"},
			},
			[]string{"xiong5", "london5"},
		},
		{
			"slice",
			[]Profile{
				{30, "brighton6"},
				{31, "bristol6"},
			},
			[]string{"brighton6", "bristol6"},
		},
		{
			"array",
			[2]Profile{
				{33, "London7"},
				{34, "Reykjavík7"},
			},
			[]string{"London7", "Reykjavík7"},
		},
		{
			"struct with slice field",
			Job{
				1,
				[]Profile{
					{30, "brighton8"},
					{31, "bristol8"},
				},
			},
			[]string{"brighton8", "bristol8"},
		},
		{
			"struct with array field",
			Job2{
				1,
				[2]Profile{
					{30, "brighton9"},
					{31, "bristol9"},
				},
			},
			[]string{"brighton9", "bristol9"},
		},
		{
			"Maps",
			map[string]string{
				"Foo": "Bar10",
				"Baz": "Boz10",
			},
			[]string{"Bar10", "Boz10"},
		},
		{
			"struct with nested Maps",
			Job3{
				1,
				map[string]Profile{
					"Foo": Profile{30, "Bar11"},
					"Baz": Profile{31, "Boz11"},
				},
			},
			[]string{"Bar11", "Boz11"},
		},
	}

	for _, test := range cases {
		t.Run("Run with object type of input", func(t *testing.T) {

			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			// with the map being in the test cases who element order may not be guarranteed, can't do exact equal any more on the collections
			// if !reflect.DeepEqual(test.expectedCalls, got) {
			// 	t.Errorf("In test '%s', got %v, want %v", test.name, got, test.expectedCalls)
			// }
			for _, data := range test.ExpectedCalls {
				assertContains(t, test.Name, got, data)
			}
		})
	}
}

func assertContains(t *testing.T, testName string, haystack []string, needle string) {
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}
	if !contains {
		t.Errorf("In test '%s', expected got %+v to contain '%s' but it didnt", testName, haystack, needle)
	}
}
