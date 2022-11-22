package fenwick

const N int = 1e6 + 10

type Fenwick interface {
	Add(int, int)
	Sum(int) int
}

type fenwick struct {
	items []int
}

func New() Fenwick {
	return &fenwick{
		items: make([]int, N),
	}
}

func lowbit(x int) int {
	return x & -x
}

func (f *fenwick) Add(x, v int) {
	for i := x; i < N; i += lowbit(i) {
		f.items[i] += v
	}
}

func (f *fenwick) Sum(x int) int {
	ret := 0
	for i := x; i != 0; i -= lowbit(i) {
		ret += f.items[i]
	}
	return ret
}
