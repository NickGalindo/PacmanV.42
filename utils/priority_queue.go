package utils

import (
  "container/heap"
)

type PriorityQueueItem[T any] struct {
  Value T
  Priority int
  Index int
}

type PriorityQueue[T any] []*PriorityQueueItem[T] 

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
  return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
  pq[i], pq[j] = pq[j], pq[i]
  pq[i].Index = i
  pq[j].Index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
  n := len(*pq)
  item := x.(*PriorityQueueItem[T])
  item.Index = n
  *pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
  old := *pq
  n := len(old)
  item := old[n-1]
  old[n-1] = nil // make GC go purrr
  item.Index = -1 // For safety?
  *pq = old[0:n-1]
  return item
}

func (pq *PriorityQueue[T]) update(item *PriorityQueueItem[T], value T, priority int) {
  item.Value = value
  item.Priority = priority
  heap.Fix(pq, item.Index) // The Push and Pop have to have any for this to work. Why? No fucking clue. This is basically Heapify
}

func (pq *PriorityQueue[T]) PushHeap(item *PriorityQueueItem[T]) {
  heap.Push(pq, item)
}

func (pq *PriorityQueue[T]) PopHeap() *PriorityQueueItem[T]{
  return heap.Pop(pq).(*PriorityQueueItem[T])
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
  pq := &PriorityQueue[T]{}
  heap.Init(pq)
  return pq
}
