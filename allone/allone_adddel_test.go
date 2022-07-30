package allone

import (
	"testing"
)

/**
 * Your AllOneAddDel object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Add(x);
 * param_2 := obj.Del(x);
 */

func TestAllOneAddDel(t *testing.T) {
	type args struct {
		op []string
		in []int
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{
			name: "normal",
			args: args{
				op: []string{"add", "del", "add", "del", "add"},
				in: []int{1, 2, 2, 1, 2},
			},
			want: []bool{true, false, true, true, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := AllOneAddDelConstructor()
			var ret bool
			for i := range tt.args.op {
				switch tt.args.op[i] {
				case "add":
					ret = this.Add(tt.args.in[i])
				case "del":
					ret = this.Del(tt.args.in[i])
				default:
					t.Errorf("Invalid operation: %v\n", tt.args.op[i])
				}

				if ret != tt.want[i] {
					t.Errorf("Index %v get %v, want %v\n", i, ret, tt.want[i])
				}
			}
		},
		)
	}
}
