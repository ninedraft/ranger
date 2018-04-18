package ranger

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type intRanger struct {
	min int
	max int
}

type intRangerJSON struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func IntRanger(min, max int) intRanger {
	if min >= max {
		panic(fmt.Sprintf("[ranger IntRanger] min parameter must be less than max, but min=%d, max=%d", min, max))
	}
	return intRanger{
		min: min,
		max: max,
	}
}

func DefaultIntRanger() intRanger {
	return IntRanger(0, 1)
}

func (ranger intRanger) Containing(i int) bool {
	return ranger.min <= i && i <= ranger.max
}

func (ranger intRanger) String() string {
	return strconv.Itoa(ranger.min) + ".." + strconv.Itoa(ranger.max)
}

func (ranger intRanger) Min() int {
	return ranger.min
}

func (ranger intRanger) Max() int {
	return ranger.max
}

func (ranger intRanger) Bounds() (min int, max int) {
	return ranger.Min(), ranger.Max()
}

func (ranger intRanger) In(r intRanger) bool {
	return ranger.min >= r.min && ranger.max <= r.max
}

func (ranger intRanger) toJSON() intRangerJSON {
	return intRangerJSON{
		Min: ranger.min,
		Max: ranger.max,
	}
}

func (ranger intRanger) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(ranger.toJSON())
	return data, err
}

func (ranger *intRanger) UnmarshalJSON(p []byte) error {
	var jsonRanger intRangerJSON
	err := json.Unmarshal(p, &jsonRanger)
	if err != nil {
		return err
	}
	*ranger = IntRanger(jsonRanger.Min, jsonRanger.Max)
	return nil
}

func (ranger intRanger) Iter() intRangerIterator {
	return ranger.IterWithStep(1)
}

func (ranger intRanger) IterWithStep(step int) intRangerIterator {
	return intRangerIterator{
		intRanger: ranger,
		value:     ranger.min,
		step:      step,
		counter:   0,
	}
}

func (ranger intRangerIterator) Shift(off int) intRanger {
	return IntRanger(ranger.min+off, ranger.max+off)
}

type intRangerIterator struct {
	intRanger
	value   int
	step    int
	counter int
}

func (iter intRangerIterator) String() string {
	return strconv.Itoa(iter.Value())
}

func (iter *intRangerIterator) Next() bool {
	min, max := iter.Bounds()
	nextVal := iter.counter*iter.step + min
	if nextVal <= max {
		iter.value = nextVal
		iter.counter++
		return true
	}
	return false
}

func (iter intRangerIterator) Value() int {
	return iter.value
}
