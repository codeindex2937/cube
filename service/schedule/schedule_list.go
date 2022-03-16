package schedule

import (
	"errors"
	"time"
)

var ErrNoItem = errors.New("key not found")

type ScheduleList struct {
	head *ScheduleNode
}

//NewScheduleList : Init structure for basic Sorted Linked List.
func NewScheduleList() *ScheduleList {
	return &ScheduleList{head: nil}
}

func (s *ScheduleList) NextSchedule() time.Time {
	if s.head == nil {
		return time.Time{}
	}
	return s.head.sortValue
}

func (s *ScheduleList) IsEmpty() bool {
	return s.head == nil
}

func (s *ScheduleList) Insert(sortKey time.Time, key interface{}, value interface{}) {
	if s.head == nil {
		s.head = NewScheduleNode(sortKey, key, value)
		return
	}

	var currentNode *ScheduleNode
	currentNode = s.head
	var previousNode *ScheduleNode
	var found bool
	newNode := NewScheduleNode(sortKey, key, value)

	for {
		if currentNode.Compare(newNode) >= 0 {
			if previousNode != nil {
				newNode.next = previousNode.next
				previousNode.next = newNode
			} else {
				newNode.next = s.head
				s.head = newNode
			}
			found = true
			break
		}

		if currentNode.next == nil {
			break
		}

		previousNode = currentNode
		currentNode = currentNode.next
	}

	if !found {
		currentNode.next = newNode
	}
}

func (s *ScheduleList) Search(key interface{}) (value interface{}, nextSched time.Time, err error) {
	currentNode := s.head
	for {
		if currentNode.key == key {
			return currentNode.value, currentNode.sortValue, nil
		}

		if currentNode.next == nil {
			break
		}
		currentNode = currentNode.next
	}
	return nil, time.Time{}, ErrNoItem
}

func (s *ScheduleList) Remove(key interface{}) error {
	if s.head == nil {
		return ErrNoItem
	}

	currentNode := s.head
	var previousNode *ScheduleNode
	for {
		if currentNode.key == key {
			if previousNode != nil {
				previousNode.next = currentNode.next
			} else {
				s.head = currentNode.next
			}
			return nil
		}

		if currentNode.next == nil {
			break
		}
		previousNode = currentNode
		currentNode = currentNode.next
	}
	return ErrNoItem
}

func (s *ScheduleList) Pop() (interface{}, error) {
	if s.head == nil {
		return nil, ErrNoItem
	}

	popped := s.head
	s.head = popped.next
	return popped.value, nil
}

func (s *ScheduleList) DisplayAll() {
	log.Infof("")
	log.Infof("head->")
	currentNode := s.head
	for {
		log.Infof("[key:%v][val:%v]->", currentNode.sortValue, currentNode.sortValue)
		if currentNode.next == nil {
			break
		}
		currentNode = currentNode.next
	}
	log.Infof("nil\n")
}
