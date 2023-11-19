package misc

import (
	"reflect"
	"testing"
)

// TestPickFields verifies that the PickFields function accurately selects the specified JSON fields from a struct.
// It checks the function's ability to handle empty field arrays and translate struct field names to JSON field names using struct tags.
func TestPickFields(t *testing.T) {
	tests := []struct {
		name   string
		data   interface{}            // Struct from which fields are to be picked
		fields []string               // list of JSON Fields to be picked from the struct
		want   map[string]interface{} // Expected map after picking the fields
	}{
		{
			name: "Pick single field",
			data: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{ID: 1, Name: "Test"},
			fields: []string{"name"},
			want:   map[string]interface{}{"name": "Test"},
		},
		{
			name: "Pick multiple fields",
			data: struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			}{ID: 2, Name: "John", Email: "john@example.com"},
			fields: []string{"id", "email"},
			want:   map[string]interface{}{"id": float64(2), "email": "john@example.com"},
		},
		{
			name: "Pick non-existent field",
			data: struct {
				Name string `json:"name"`
			}{Name: "Jane"},
			fields: []string{"age"},
			want:   map[string]interface{}{},
		},
		{
			name: "Empty fields array",
			data: struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			}{ID: 3, Name: "Jane", Email: "jane@example.com"},
			fields: []string{}, // No fields specified, expecting an empty result.
			want:   map[string]interface{}{},
		},
		{
			name: "JSON tags translation",
			data: struct {
				ID    int    `json:"id"`
				Name  string `json:"full_name"`
				Email string `json:"email_address"`
			}{ID: 4, Name: "Doe", Email: "doe@example.com"},
			fields: []string{"full_name", "email_address"}, // Use JSON tag names to pick the fields.
			want:   map[string]interface{}{"full_name": "Doe", "email_address": "doe@example.com"},
		},
		{
			name: "Struct fields vs JSON field names",
			data: struct {
				ID    int    `json:"id"`
				Name  string `json:"full_name"`
				Email string `json:"email_address"`
			}{ID: 5, Name: "Smith", Email: "smith@example.com"},
			fields: []string{"Name", "Email"}, // Incorrect field names, as they do not match the JSON tag names.
			want:   map[string]interface{}{},  // Expect an empty map as a result.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PickFields(tt.data, tt.fields...)
			if !reflect.DeepEqual(got, tt.want) {
				// Log details if there's a discrepancy between what we got and what we wanted.
				for k, v := range got {
					t.Logf("got[%v] = %v (%T)", k, v, v)
				}
				for k, v := range tt.want {
					t.Logf("want[%v] = %v (%T)", k, v, v)
				}
				t.Errorf("PickFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}
