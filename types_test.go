package asql_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/asql"
)

func TestSortDirection(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		testCases := []struct {
			name string

			value asql.SortDirection

			expectErr bool
		}{
			{
				name: "OK/Empty",

				value: asql.SortDirectionNone,
			},
			{
				name: "OK/Asc",

				value: asql.SortDirectionAsc,
			},
			{
				name: "OK/Desc",

				value: asql.SortDirectionDesc,
			},
			{
				name: "KO/Invalid",

				value: "invalid",

				expectErr: true,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				customValidator := validator.New(validator.WithRequiredStructEnabled())
				asql.RegisterSortDirection(customValidator)

				toValidate := struct {
					Value asql.SortDirection `validate:"omitempty,sort_direction"`
				}{
					Value: testCase.value,
				}

				err := customValidator.Struct(toValidate)

				if testCase.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}

		t.Run("OtherTypes", func(t *testing.T) {
			customValidator := validator.New(validator.WithRequiredStructEnabled())
			asql.RegisterSortDirection(customValidator)

			toValidate := struct {
				Value interface{} `validate:"omitempty,sort_direction"`
			}{
				Value: "asc",
			}

			err := customValidator.Struct(toValidate)
			require.Error(t, err)
		})
	})
}
