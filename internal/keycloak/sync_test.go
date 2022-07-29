package keycloak

import (
	"strings"
	"testing"
)

func TestConstructGroupUpdateRequestBody(t *testing.T) {
	var testCases = map[string]struct {
		input  *Group
		expect string
	}{
		"simple update": {
			input: &Group{
				ID: "f6697da3-016a-43cd-ba9f-3f5b91b45302",
				GroupUpdateRepresentation: GroupUpdateRepresentation{
					Name: "foo",
					Attributes: map[string][]string{
						"group-lagoon-project-ids": {`{"foo":[1,2,34]}`},
						"lagoon-projects":          {`1,2,34`},
					},
				},
			},
			expect: `{"name":"foo","attributes":{` +
				`"group-lagoon-project-ids":["{\"foo\":[1,2,34]}"],` +
				`"lagoon-projects":["1,2,34"]` +
				`}}`,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(tt *testing.T) {
			reader, err := constructGroupUpdateRequestBody(tc.input)
			if err != nil {
				tt.Fatal(err)
			}
			result := strings.Builder{}
			_, err = reader.WriteTo(&result)
			if err != nil {
				tt.Fatal(err)
			}
			if result.String() != tc.expect {
				tt.Fatalf("got %s, expected %s", result.String(), tc.expect)
			}
		})
	}
}
