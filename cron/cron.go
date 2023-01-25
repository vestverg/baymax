package cron

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/vestverg/baymax/bits"
)

type CronFieldType int

const (
	Second CronFieldType = 1 + iota
	Minute
	Hour
	DOM // day of month
	Month
	DOW // day of week

)

var fieldRange = map[CronFieldType]FieldRange{
	Second: {
		min: 0,
		max: 59,
	},
	Minute: {
		min: 0,
		max: 59,
	},
	Hour: {
		min: 0,
		max: 23,
	},
	DOM: {
		min: 1,
		max: 31,
	},
	Month: {
		min: 1,
		max: 12,
		names: map[string]int{
			"jan": 1,
			"feb": 2,
			"mar": 3,
			"apr": 4,
			"may": 5,
			"jun": 6,
			"jul": 7,
			"aug": 8,
			"sep": 9,
			"oct": 10,
			"nov": 11,
			"dec": 12,
		},
	},
	DOW: {
		min: 1,
		max: 7,
		names: map[string]int{
			"sun": 0,
			"mon": 1,
			"tue": 2,
			"wed": 3,
			"thu": 4,
			"fri": 5,
			"sat": 6,
		},
	},
}

type FieldRange struct {
	min   int
	max   int
	names map[string]int
}

func (fr *FieldRange) IsValid(value int) bool {
	if value < fr.min || value > fr.max {
		return false
	}
	return true
}

//CronExpression Parse the given crontab expression  string into a CronExpression. The string has six single space-separated cron and date fields:
//┌───────────── second (0-59)
//│ ┌───────────── minute (0 - 59)
//│ │ ┌───────────── hour (0 - 23)
//│ │ │ ┌───────────── day of the month (1 - 31)
//│ │ │ │ ┌───────────── month (1 - 12) (or JAN-DEC)
//│ │ │ │ │ ┌───────────── day of the week (0 - 7)
//│ │ │ │ │ │          (0 or 7 is Sunday, or MON-SUN)
//│ │ │ │ │ │
//* * * * * *
//The following rules apply:
//A field may be an asterisk ( *), which always stands for "first-last". For the "day of the month" or "day of the week" fields, a question mark ( ?) may be used instead of an asterisk.
//Ranges of numbers are expressed by two numbers separated with a hyphen ( -). The specified range is inclusive.
//Following a range (or *) with /n specifies the interval of the number's value through the range.
//English names can also be used for the "month" and "day of week" fields. Use the first three letters of the particular day or month (case does not matter).
//The "day of month" and "day of week" fields can contain a L-character, which stands for "last", and has a different meaning in each field:
//In the "day of month" field, L stands for "the last day of the month". If followed by an negative offset (i.e. L-n), it means " nth-to-last day of the month". If followed by W (i.e. LW), it means "the last weekday of the month".
//In the "day of week" field, L stands for "the last day of the week". If prefixed by a number or three-letter name (i.e. dL or DDDL), it means "the last day of week d (or DDD) in the month".
//The "day of month" field can be nW, which stands for "the nearest weekday to day of the month n". If n falls on Saturday, this yields the Friday before it. If n falls on Sunday, this yields the Monday after, which also happens if n is 1 and falls on a Saturday (i.e. 1W stands for "the first weekday of the month").
//The "day of week" field can be d#n (or DDD#n), which stands for "the n-th day of week d (or DDD) in the month".

type CronExpression struct {
	fields     []CronField
	expression string
}

func (cr *CronExpression) Next(from time.Time) time.Time {
	next := from
	for i := 0; i < 366; i++ {
		temp := next
		for _, field := range cr.fields {
			temp = field.NextOrSame(temp)
		}
		if temp == next {
			return next
		}
		next = temp
	}

	return next
}

type CronField struct {
	fieldType  CronFieldType
	fieldRange FieldRange
	bits       int64
}

func (cr *CronField) getHigherFieldType() CronFieldType {
	fieldType := cr.fieldType + 1
	if fieldType >= 6 {
		fieldType = Month
	}
	return fieldType
}

func Reset(ft CronFieldType, date time.Time) time.Time {
	switch ft {
	case Second:
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), 0, date.Location())
	case Minute:
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, date.Location())
	case Hour:
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 0, 0, 0, date.Location())
	case DOM:
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	case Month:
		return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	case DOW:
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	}
	return date
}

func (cr *CronField) NextOrSame(date time.Time) time.Time {
	current := cr.getPartOfTime(date)
	next := cr.Next(current)
	if next == -1 {
		date = cr.shift(date)
		next = cr.Next(0)
	}
	if next == current {
		return date
	} else {
		current = cr.getPartOfTime(date)
		count := 0
		for current != next && count < 366 {
			date = cr.evaluateNext(date, next)
			current = cr.getPartOfTime(date)
			next = cr.Next(current)
			if next == -1 {
				date = cr.shift(date)
				next = cr.Next(0)
			}
			count++
		}
		date = Reset(cr.fieldType, date)
	}
	return date
}

func (cr *CronField) evaluateNext(date time.Time, next int) time.Time {
	current := cr.getPartOfTime(date)
	if current < next {
		date = cr.Add(date, next-current)
	} else {
		amount := cr.fieldRange.max - current + next + 1 - cr.fieldRange.min
		date = cr.Add(date, amount)
	}
	return date
}

func (cr *CronField) getPartOfTime(date time.Time) int {
	switch cr.fieldType {
	case Minute:
		return date.Minute()
	case Second:
		return date.Second()
	case Hour:
		return date.Hour()
	case DOM:
		return date.Day()
	case Month:
		return int(date.Month())
	case DOW:
		return int(date.Weekday())
	}
	return -1

}
func (cr *CronField) GetBit(idx int) int {
	return int(cr.bits>>int64(idx)) % 2

}
func (cr *CronField) Next(idx int) int {
	result := cr.bits & (math.MaxInt64 << idx)
	if result != 0 {
		return bits.NumberOfTrailingZeros64(result)
	}
	return -1

}

func (cr *CronField) SetBits(fr *FieldRange) {
	if fr.min == fr.max {
		cr.SetBit(fr.min)
		return
	}
	minMask := math.MaxInt64 << fr.min
	maxMask := math.MaxInt64 >> (fr.max + 1)
	cr.bits |= int64(minMask & maxMask)

}

func (cr *CronField) SetBit(idx int) {
	cr.bits = bits.SetBit(cr.bits, idx)
}

func (cr *CronField) shift(date time.Time) time.Time {
	switch cr.fieldType {
	case Minute:
		return Reset(cr.getHigherFieldType(), date.Add(time.Hour))
	case Second:
		return Reset(cr.getHigherFieldType(), date.Add(time.Minute))
	case Hour:
		return Reset(cr.getHigherFieldType(), date.Add(24*time.Hour))
	case DOM: //???
		y, m, _ := date.Date()
		return time.Date(y, m+1, 1, 0, 0, 0, 0, date.Location())
	case DOW:
		cr.Add(date, int(date.Weekday()-time.Monday))
	case Month:
		y, _, _ := date.Date()
		return time.Date(y+1, 1, 1, 0, 0, 0, 0, date.Location())
	}
	return date

}

func (cr *CronField) Add(date time.Time, value int) time.Time {
	switch cr.fieldType {
	case Minute:
		return date.Add(time.Duration(value) * time.Minute)
	case Second:
		return date.Add(time.Duration(value) * time.Second)
	case Hour:
		return date.Add(time.Duration(value) * time.Hour)
	case DOM, DOW:
		return time.Date(date.Year(), date.Month(), date.Day()+value, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
	case Month:
		return time.Date(date.Year(), date.Month()+time.Month(value), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
	}
	return date

}

func Parse(value string) (*CronExpression, error) {
	if value == "" {
		return nil, fmt.Errorf("expression string must not be empty")
	}
	value = strings.Replace(value, "?", "*", -1)
	//todo: add macros
	fields := strings.Split(value, " ")
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid cron expression, cron expression must consist of 6 fields")
	}
	seconds, err := parseField(fields[0], Second)
	if err != nil {
		return nil, fmt.Errorf("failed to parse seconds:f %w", err)
	}
	minutes, err := parseField(fields[1], Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to parse minutes: %w", err)
	}
	hours, err := parseField(fields[2], Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hours: %w", err)
	}
	dom, err := parseField(fields[3], DOM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse days of moth: %w", err)
	}
	months, err := parseField(fields[4], Month)
	if err != nil {
		return nil, fmt.Errorf("failed to parse months: %w", err)
	}
	dow, err := parseField(fields[5], DOW)
	if err != nil {
		return nil, fmt.Errorf("failed to parse days of week: %w", err)
	}

	return &CronExpression{
		fields:     []CronField{*dow, *months, *dom, *hours, *minutes, *seconds},
		expression: value,
	}, nil
}

func parseField(field string, fieldType CronFieldType) (*CronField, error) {
	if field == "" {
		return nil, fmt.Errorf("field is empty")
	}
	fr := fieldRange[fieldType]
	cronField := CronField{fieldType: fieldType, fieldRange: fr}
	parts := strings.Split(field, ",")
	for _, part := range parts {
		slash := strings.Index(part, "/")
		if slash == -1 {
			r, err := parseRange(part, fieldType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse field %w", err)
			}
			cronField.SetBits(r)
		} else {
			rangeStr := part[:slash]

			r, err := parseRange(rangeStr, fieldType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse field %w", err)
			}
			if strings.Index(rangeStr, "-") == -1 {
				// if it's not interval
				r = &FieldRange{min: r.min, max: fr.max}
			}

			deltaStr := part[slash+1:]
			delta, err := strconv.Atoi(deltaStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse delta %w", err)
			}
			if delta <= 0 {
				return nil, fmt.Errorf("delta is negative")
			}
			if delta == 1 {
				cronField.SetBits(r)
			}
			for i := r.min; i <= r.max; i += delta {
				cronField.SetBit(i)
			}
		}
	}
	return &cronField, nil
}

func parseRange(val string, fieldType CronFieldType) (*FieldRange, error) {
	fr, ok := fieldRange[fieldType]
	if !ok {
		return nil, fmt.Errorf("unknown field type")
	}
	if val == "*" {
		return &fr, nil
	}
	hypen := strings.Index(val, "-")
	if hypen == -1 {
		res, err := parseNameOrVal(val, fr)
		if err != nil {
			return nil, fmt.Errorf("invalid range value")
		}
		if !fr.IsValid(res) {
			return nil, fmt.Errorf("value out of range")
		}
		return &FieldRange{min: res, max: res}, nil
	}

	min, err := parseNameOrVal(val[0:hypen], fr)
	if err != nil {
		return nil, fmt.Errorf("invalid min range value")
	}
	if !fr.IsValid(min) {
		return nil, fmt.Errorf("min value out of range")
	}

	max, err := parseNameOrVal(val[hypen+1:], fr)
	if err != nil {
		return nil, fmt.Errorf("invalid max range value")
	}
	if !fr.IsValid(min) {
		return nil, fmt.Errorf("max value out of range")
	}

	return &FieldRange{
		min: min,
		max: max,
	}, nil

}

func parseNameOrVal(val string, fr FieldRange) (int, error) {
	if fr.names != nil {
		if v, ok := fr.names[strings.ToLower(val)]; ok {
			return v, nil
		}
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return -1, fmt.Errorf("invalid range value")
	}
	if v < 0 {
		return -1, fmt.Errorf("invalid range value")
	}
	return v, err
}
