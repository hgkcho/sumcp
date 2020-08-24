package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	if err := os.Chdir("../../testdata"); err != nil {
		panic(err)
	}
	type args struct {
		sources []string
		target  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				sources: []string{"a.txt", "b.txt"},
				target:  "c.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.sources, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}

			var wantStr, actualStr string
			for _, source := range tt.args.sources {
				f, err := os.Open(source)
				if err != nil {
					fmt.Println(err)
				}
				defer f.Close()
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					wantStr += scanner.Text()
				}
			}

			f, err := os.Open(tt.args.target)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				actualStr += scanner.Text()
			}

			if actualStr != wantStr {
				t.Errorf("run() actual data:\n %v \n, but want \n %v", actualStr, wantStr)
			}
		})
	}
}
