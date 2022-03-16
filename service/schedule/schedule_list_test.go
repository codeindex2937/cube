package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLinkedNodeCompare(t *testing.T) {
	now := time.Now()
	n1 := NewScheduleNode(now, "1", "node1")
	n2 := NewScheduleNode(now.Add(1*time.Second), "2", "node2")

	assert.Equal(t, -1, n1.Compare(n2))
	assert.Equal(t, 1, n2.Compare(n1))
	assert.Equal(t, 0, n1.Compare(n1))
}

func TestScheduleListSearch(t *testing.T) {
	now := time.Now()
	targetSched := now.Add(4 * time.Second)
	schedList := NewScheduleList()
	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	schedList.Insert(targetSched, "4", "string4")
	schedList.Insert(now.Add(2*time.Second), "2", "string2")

	ret, nextSched, err := schedList.Search("4")
	assert.NoError(t, err)
	assert.Equal(t, "string4", ret)
	assert.Equal(t, nextSched, targetSched)
}

func TestScheduleListPop(t *testing.T) {
	now := time.Now()
	schedList := NewScheduleList()
	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	schedList.Insert(now.Add(4*time.Second), "4", "string4")
	schedList.Insert(now.Add(2*time.Second), "2", "string2")

	ret, err := schedList.Pop()
	assert.NoError(t, err)
	assert.Equal(t, "string2", ret)
	ret, err = schedList.Pop()
	assert.NoError(t, err)
	assert.Equal(t, "string3", ret)
	ret, err = schedList.Pop()
	assert.NoError(t, err)
	assert.Equal(t, "string4", ret)
	ret, err = schedList.Pop()
	assert.Equal(t, ErrNoItem, err)
	assert.Nil(t, ret)
}

func TestScheduleListRemove(t *testing.T) {
	now := time.Now()
	schedList := NewScheduleList()
	schedList.Insert(now.Add(4*time.Second), "4", "string4")
	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	assert.Nil(t, schedList.Remove("3"))
	assert.Nil(t, schedList.Remove("4"))
	assert.Equal(t, ErrNoItem, schedList.Remove("4"))
}

func TestScheduleListIsEmpty(t *testing.T) {
	now := time.Now()
	schedList := NewScheduleList()
	assert.True(t, schedList.IsEmpty())

	schedList.Insert(now.Add(4*time.Second), "4", "string4")
	assert.False(t, schedList.IsEmpty())

	value, err := schedList.Pop()
	assert.NoError(t, err)
	assert.True(t, schedList.IsEmpty())
	assert.Equal(t, "string4", value.(string))

	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	assert.False(t, schedList.IsEmpty())

	err = schedList.Remove("3")
	assert.NoError(t, err)
	assert.True(t, schedList.IsEmpty())
}
