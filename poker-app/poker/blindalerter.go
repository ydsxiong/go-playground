package poker

import (
	"fmt"
	"io"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

/**
Remember that any type can implement an interface, not just structs.
If you are making a library that exposes an interface with one function defined,
then it is a common idiom to also expose a MyInterfaceFunc type.
This type will be a func which will also implement your interface.
That way users of your interface have the option to implement your interface with just a function;
rather than having to create an empty struct type.
*/
type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (ba BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	ba(duration, amount, to)
}

func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}
