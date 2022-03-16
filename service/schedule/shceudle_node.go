package schedule

import "time"

type ScheduleNode struct {
	sortValue time.Time
	key       interface{}
	value     interface{}
	next      *ScheduleNode
}

func NewScheduleNode(sortValue time.Time, key interface{}, value interface{}) *ScheduleNode {
	return &ScheduleNode{sortValue, key, value, nil}
}

func (node *ScheduleNode) Compare(that *ScheduleNode) int {
	if node.sortValue.Before(that.sortValue) {
		return -1
	}
	if node.sortValue.After(that.sortValue) {
		return 1
	}
	return 0
}

func (node *ScheduleNode) Value() *interface{} {
	return &node.value
}
