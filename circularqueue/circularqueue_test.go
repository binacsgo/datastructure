package circularqueue

import (
	"testing"
)

type fields struct {
	head int
	tail int
	size int
	data []interface{}
}

func TestCircularQueue_Size(t *testing.T) {
	tests := []struct {
		name string
		size int
		want int
	}{
		{
			name: "zero",
			size: 0,
			want: defaultCircularQueueSize,
		},
		{
			name: "non-zero",
			size: 10,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewCircularQueue(tt.size)
			if got := q.Size(); got != tt.want {
				t.Errorf("CircularQueue.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_CountUsed(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "empty queue",
			fields: fields{
				head: 0,
				tail: 0,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 0,
		},
		{
			name: "non-empty queue 1",
			fields: fields{
				head: 0,
				tail: 4,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 4,
		},
		{
			name: "non-empty queue 2",
			fields: fields{
				head: 4,
				tail: 0,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 7,
		},
		{
			name: "full queue",
			fields: fields{
				head: 0,
				tail: 10,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &CircularQueue{
				head: tt.fields.head,
				tail: tt.fields.tail,
				size: tt.fields.size,
				data: tt.fields.data,
			}
			if got := q.CountUsed(); got != tt.want {
				t.Errorf("CircularQueue.CountUsed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_CountUnused(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "empty queue",
			fields: fields{
				head: 0,
				tail: 0,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 10,
		},
		{
			name: "non-empty queue 1",
			fields: fields{
				head: 0,
				tail: 4,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 6,
		},
		{
			name: "non-empty queue 2",
			fields: fields{
				head: 4,
				tail: 0,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 3,
		},
		{
			name: "full queue",
			fields: fields{
				head: 0,
				tail: 10,
				size: 10,
				data: make([]interface{}, 11),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &CircularQueue{
				head: tt.fields.head,
				tail: tt.fields.tail,
				size: tt.fields.size,
				data: tt.fields.data,
			}
			if got := q.CountUnused(); got != tt.want {
				t.Errorf("CircularQueue.CountUnused() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularQueue_All(t *testing.T) {
	q := NewCircularQueue(10)
	for i := 0; i < 10; i++ {
		if ok := q.EnQueue(i); !ok {
			t.Errorf("Step 1: EnQueue() = %v, want %v", ok, true)
		}
	}

	if item := q.Front(); item.(int) != 0 {
		t.Errorf("Step 2: Front() = %v, want %v", item.(int), 0)
	}

	if item, ok := q.DeQueue(); !ok || item.(int) != 0 {
		t.Errorf("Step 3: DeQueue() = %v, %v, want %v, %v", item.(int), ok, 0, true)
	}

	if item := q.Rear(); item.(int) != 9 {
		t.Errorf("Step 4: Rear() = %v, want %v", item.(int), 9)
	}
}
