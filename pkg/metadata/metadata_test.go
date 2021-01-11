package metadata

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

type Data struct {
	Id     int      `annotation:"id"`
	Name   string   `annotation:"name"`
	Age    uint     `label:"age"`
	Skills []string `label:"skills"`
}

func TestParse(t *testing.T) {
	type args struct {
		metadata v1.ObjectMeta
		prefix   string
	}
	tests := []struct {
		name    string
		args    args
		want    Data
		wantErr bool
	}{
		{
			name: "",
			args: args{
				metadata: v1.ObjectMeta{
					Annotations: map[string]string{
						"prefix/id":   "1",
						"prefix/name": "John",
					},
					Labels: map[string]string{
						"prefix/age":    "30",
						"prefix/skills": "cooking,swimming,driving",
					},
				},
				prefix: "prefix",
			},
			want: Data{
				Id:     1,
				Name:   "John",
				Age:    30,
				Skills: []string{"cooking", "swimming", "driving"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := Data{}
			if err := Parse(tt.args.metadata, &data, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(data, tt.want) {
				t.Errorf("Parse() struct = %v, want %v", data, tt.want)
			}
		})
	}
}
