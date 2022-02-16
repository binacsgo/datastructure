package circularqueue

const defaultCircularQueueSize = 10

type CircularQueue struct {
	head int
	tail int
	size int
	data []interface{}
}

func NewCircularQueue(size int) *CircularQueue {
	if size <= 0 {
		size = defaultCircularQueueSize
	}
	return &CircularQueue{
		head: 0,
		tail: 0,
		size: size,
		data: make([]interface{}, size+1),
	}
}

func (q *CircularQueue) Size() int {
	return q.size
}

func (q *CircularQueue) CountUsed() int {
	if q.tail >= q.head {
		return q.tail - q.head
	}
	return q.tail + q.size + 1 - q.head
}

func (q *CircularQueue) CountUnused() int {
	return q.Size() - q.CountUsed()
}

func (q *CircularQueue) IsEmpty() bool {
	return q.tail == q.head
}

func (q *CircularQueue) IsFull() bool {
	return (q.tail+1)%(q.size+1) == q.head
}

func (q *CircularQueue) EnQueue(item interface{}) bool {
	if q.IsFull() {
		return false
	}
	q.data[q.tail] = item
	q.tail = (q.tail + 1) % (q.size + 1)
	return true
}

func (q *CircularQueue) DeQueue() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	item := q.data[q.head]
	q.head = (q.head + 1) % (q.size + 1)
	return item, true
}

func (q *CircularQueue) Front() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.data[q.head]
}

func (q *CircularQueue) Rear() interface{} {
	if q.IsEmpty() {
		return nil
	}
	idx := (q.tail - 1 + q.size + 1) % (q.size + 1)
	return q.data[idx]
}
