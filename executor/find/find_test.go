package find

import "testing"

type Finder []int

func (f Finder) Len() int {
	return len(f)
}

func (f Finder) Compare(i int, val interface{}) int {
	return f[i] - val.(int)
}

func TestFind(t *testing.T) {
	var fd Finder = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	type args struct {
		data  Interface
		input interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test",
			args: args{
				data:  fd,
				input: 4,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.args.data, tt.args.input); got != tt.want {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
