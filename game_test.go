package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Move(t *testing.T) {
	type testData struct {
		arr []int
		res []int
	}

	datas := []*testData{
		{arr: []int{0, 0, 0, 0}, res: []int{0, 0, 0, 0}},

		{arr: []int{2, 0, 0, 0}, res: []int{2, 0, 0, 0}},
		{arr: []int{0, 2, 0, 0}, res: []int{2, 0, 0, 0}},
		{arr: []int{0, 0, 2, 0}, res: []int{2, 0, 0, 0}},
		{arr: []int{0, 0, 0, 2}, res: []int{2, 0, 0, 0}},

		{arr: []int{2, 2, 0, 0}, res: []int{2, 2, 0, 0}},
		{arr: []int{2, 0, 2, 0}, res: []int{2, 2, 0, 0}},
		{arr: []int{2, 0, 0, 2}, res: []int{2, 2, 0, 0}},
		{arr: []int{0, 2, 2, 0}, res: []int{2, 2, 0, 0}},
		{arr: []int{0, 2, 0, 2}, res: []int{2, 2, 0, 0}},
		{arr: []int{0, 0, 2, 2}, res: []int{2, 2, 0, 0}},

		{arr: []int{2, 2, 2, 0}, res: []int{2, 2, 2, 0}},
		{arr: []int{2, 2, 0, 2}, res: []int{2, 2, 2, 0}},
		{arr: []int{2, 0, 2, 2}, res: []int{2, 2, 2, 0}},
		{arr: []int{0, 2, 2, 2}, res: []int{2, 2, 2, 0}},

		{arr: []int{2, 2, 2, 2}, res: []int{2, 2, 2, 2}},
	}
	for _, data := range datas {
		move(data.arr)
		assert.Equal(t, data.res, data.arr)
	}

}

func Test_Merge(t *testing.T) {
	type testData struct {
		arr []int
		res []int
	}

	datas := []*testData{
		{arr: []int{0, 0, 0, 0}, res: []int{0, 0, 0, 0}},

		{arr: []int{2, 0, 0, 0}, res: []int{2, 0, 0, 0}},
		{arr: []int{0, 2, 0, 0}, res: []int{2, 0, 0, 0}},
		{arr: []int{0, 0, 2, 0}, res: []int{2, 0, 0, 0}},
		{arr: []int{0, 0, 0, 2}, res: []int{2, 0, 0, 0}},

		{arr: []int{2, 2, 0, 0}, res: []int{4, 0, 0, 0}},
		{arr: []int{2, 0, 2, 0}, res: []int{4, 0, 0, 0}},
		{arr: []int{2, 0, 0, 2}, res: []int{4, 0, 0, 0}},
		{arr: []int{0, 2, 2, 0}, res: []int{4, 0, 0, 0}},
		{arr: []int{0, 2, 0, 2}, res: []int{4, 0, 0, 0}},
		{arr: []int{0, 0, 2, 2}, res: []int{4, 0, 0, 0}},

		{arr: []int{2, 2, 2, 0}, res: []int{4, 2, 0, 0}},
		{arr: []int{2, 2, 0, 2}, res: []int{4, 2, 0, 0}},
		{arr: []int{2, 0, 2, 2}, res: []int{4, 2, 0, 0}},
		{arr: []int{0, 2, 2, 2}, res: []int{4, 2, 0, 0}},

		{arr: []int{2, 2, 2, 2}, res: []int{4, 4, 0, 0}},
	}
	for _, data := range datas {
		move(data.arr)
		merge(data.arr)
		move(data.arr)
		assert.Equal(t, data.res, data.arr)
	}

}
