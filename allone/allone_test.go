package allone

import (
	"testing"
)

/**
 * Your AllOne object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Inc(key);
 * obj.Dec(key);
 * param_3 := obj.GetMaxKey();
 * param_4 := obj.GetMinKey();
 */

func TestAllOne(t *testing.T) {
	type args struct {
		op []string
		in []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "normal",
			args: args{
				op: []string{"inc", "inc", "getMaxKey", "getMinKey", "inc", "getMaxKey", "getMinKey"},
				in: []string{"hello", "hello", "", "", "leet", "", ""},
			},
			want: []string{"", "", "hello", "hello", "", "hello", "leet"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := Constructor()
			for i := range tt.args.op {
				switch tt.args.op[i] {
				case "inc":
					this.Inc(tt.args.in[i])
				case "dec":
					this.Dec(tt.args.in[i])
				case "getMaxKey":
					if key := this.GetMaxKey(); key != tt.want[i] {
						t.Errorf("GetMaxKey get %v, want %v\n", key, tt.want[i])
					}
				case "getMinKey":
					if key := this.GetMinKey(); key != tt.want[i] {
						t.Errorf("GetMinKey get %v, want %v\n", key, tt.want[i])
					}
				default:
					t.Errorf("Invalid operation: %v\n", tt.args.op[i])
				}
			}
		},
		)
	}
}
