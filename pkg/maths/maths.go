package maths

import "errors"

func Avg(s []float64) (float64, error) {
	if len(s) == 0 {
		return 0.0, errors.New("the size of array is zero")
	}
	sum := 0.0
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s)), nil
}

func Max(s []float64) (float64, error) {
	if len(s) == 0 {
		return 0.0, errors.New("the size of array is zero")
	}
	max := s[0]
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	return max, nil
}

func Min(s []float64) (float64, error) {
	if len(s) == 0 {
		return 0.0, errors.New("the size of array is zero")
	}
	min := s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
	}
	return min, nil
}
