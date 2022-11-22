package fenwick

import (
	"testing"
)

type Manager struct {
	f Fenwick
	n int
}

func Constructor(n int) Manager {
	return Manager{
		f: New(),
		n: n,
	}
}

func (this *Manager) Reserve() int {
	l, r := 1, this.n
	for l < r {
		m := (l + r) >> 1
		if this.f.Sum(m) >= m {
			l = m + 1
		} else {
			r = m
		}
	}
	this.f.Add(l, 1)
	return l
}

func (this *Manager) Unreserve(seatNumber int) {
	this.f.Add(seatNumber, -1)
}

func TestFenwick(t *testing.T) {
	obj := Constructor(N)
	ops := []string{"reserve", "reserve", "unreserve", "reserve", "reserve", "reserve", "reserve", "unreserve"}
	wants := []int{1, 2, 2, 2, 3, 4, 5, 5}

	if len(ops) != len(wants) {
		t.Errorf("len(ops) != len(got)\n")
	}

	for i := 0; i < len(ops); i++ {
		op, want := ops[i], wants[i]
		switch op {
		case "reserve":
			got := obj.Reserve()
			if want != got {
				t.Errorf("Unexpect result, want: %v, got: %v\n", want, got)
			}
		case "unreserve":
			obj.Unreserve(want)
		default:
			t.Errorf("Invalid op: %v\n", op)
		}
	}
}
