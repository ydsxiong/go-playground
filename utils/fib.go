package utils

import (
	"fmt"
	"reflect"
)

func fib(n int) uint64 {
	if n == 1 || n == 2 {
		return 1
	}

	return fib(n-2) + fib(n-1)
}

func fib_v2() func() uint64 {
	var current, next uint64 = 0, 1
	return func() uint64 {
		current, next = next, current+next
		return current
	}
}

type T struct {
	A int
	B string
}

type Fooer interface {
	Foo1()
	Foo2()
	Foo3()
}


type SuperFooer struct {
	Fooer
	string
}

func sumOf10naturalnumbersbyimperativeapproach() int {

	sf := SuperFooer{}
	sf.Foo1() // obtained through embedding/composition, but not actually implmenting any of its methods, so it would panic here at runtime.


	var x float64 = 3.4
	fmt.Println("value:", reflect.ValueOf(x).String())

	var x1 uint8 = 'x'
	v1 := reflect.ValueOf(x1)
	fmt.Println("type of:", reflect.TypeOf(x1))
	fmt.Println("type kind:", reflect.TypeOf(x1).Kind())
	fmt.Println("type string:", reflect.TypeOf(x1).String())

	fmt.Println("value string:", v1.String())
	fmt.Println("value:", v1.Uint())
	fmt.Println("value type:", v1.Type())                            // uint8.
	fmt.Println("value kind is uint8: ", v1.Kind() == reflect.Uint8) // true.
	x1 = uint8(v1.Uint())                                            // v.Uint returns a uint64.

	type MyInt int
	var x2 MyInt = 7
	v2 := reflect.ValueOf(x2)

	fmt.Println("type of:", reflect.TypeOf(x2))
	fmt.Println("type kind:", reflect.TypeOf(x2).Kind())
	fmt.Println("type string:", reflect.TypeOf(x2).String())

	fmt.Println("value string:", v2.String())
	fmt.Println("value:", v2.Int())
	fmt.Println("value type:", v2.Type())     // uint8.
	fmt.Println("value kind is: ", v2.Kind()) // true.
	x2 = MyInt(v2.Int())

	fmt.Println("===================")

	t := T{23, "skidoo"}
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	s.Field(0).SetInt(77)
	s.Field(1).SetString("Sunset Strip")
	fmt.Println("t is now", t)

	sum, num, n := 0, 10, 1
	for num > 0 {
		if n*n%5 == 0 {
			sum += n
			num--
		}
		n++
	}

	return sum
}

func sumOf10naturalnumbersbyRecursionapproach() int {
	return sumUpNumber(1, 10, 0)
}

func sumUpNumber(n, times, sum int) int {
	if times == 0 {
		return sum
	}
	nextN := n + 1
	var nextSum int
	var nexttimes int
	if n*n%5 == 0 {
		nextSum = sum + n
		nexttimes = times - 1
	} else {
		nextSum = sum
		nexttimes = times
	}
	return sumUpNumber(nextN, nexttimes, nextSum)
}

func sumOf10naturalnumbersbyfunctionalapproach() int {
	// Reduce(Filter(pred, createValues), sum, uint64).(uint64)
	return 0
}


