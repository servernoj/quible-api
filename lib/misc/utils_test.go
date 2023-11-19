package misc

import (
	"reflect"
	"testing"
)

func TestPickFields(t *testing.T) {
	tests := []struct {
		name   string
		data   interface{}
		fields []string
		want   map[string]interface{}
	}{
		{
			name: "Pick single field",
			data: struct {
				ID   int
				Name string
			}{ID: 1, Name: "Test"},
			fields: []string{"Name"},
			want:   map[string]interface{}{"Name": "Test"},
		},
		{
			name: "Pick multiple fields",
			data: struct {
				ID    int
				Name  string
				Email string
			}{ID: 2, Name: "John", Email: "john@example.com"},
			fields: []string{"ID", "Email"},
			want:   map[string]interface{}{"ID": 2, "Email": "john@example.com"},
		},
		{
			name:   "Pick non-existent field",
			data:   struct{ Name string }{Name: "Jane"},
			fields: []string{"Age"},
			want:   map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PickFields(tt.data, tt.fields...)
			if !reflect.DeepEqual(got, tt.want) {
				for k, v := range got {
					t.Logf("got[%v] = %v (%T)", k, v, v)
				}
				for k, v := range tt.want {
					t.Logf("want[%v] = %v (%T)", k, v, v)
				}
				t.Errorf("PickFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
