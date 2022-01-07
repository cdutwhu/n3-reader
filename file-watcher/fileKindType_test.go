package filewatcher

import "testing"

func Test_getFileType(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want EmFileType
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				file: "../executable/n3-reader/config.json",
			},
			want: Text,
		},
		{
			name: "",
			args: args{
				file: "../executable/n3-reader/n3-reader",
			},
			want: Executable,
		},
		{
			name: "",
			args: args{
				file: "/home/qmiao/Downloads/MLY-zh-cn.pdf",
			},
			want: Text,
		},
		{
			name: "",
			args: args{
				file: "/home/qmiao/Downloads/魔神英雄伝2_1.rmvb",
			},
			want: Video,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileType(tt.args.file); got != tt.want {
				t.Errorf("getFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}
