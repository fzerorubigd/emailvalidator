package emailvalidator

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixture struct {
	email      string
	free       ValidationState
	disposable ValidationState
	blackList  ValidationState
	fail       bool
}

var (
	testFixtures = []fixture{
		{
			email:      "test.with.dot@gmail.com",
			free:       ValidationStateTrue,
			disposable: ValidationStateFalse,
			blackList:  ValidationStateFalse,
		},
		{
			email:      "test.with.dot+extra@gmail.com",
			free:       ValidationStateTrue,
			disposable: ValidationStateFalse,
			blackList:  ValidationStateFalse,
		},
		{
			email:      "test@things.10mail.org",
			free:       ValidationStateFalse,
			disposable: ValidationStateTrue,
			blackList:  ValidationStateFalse,
		},
		{
			email:      "test@things.more.10mail.org",
			free:       ValidationStateFalse,
			disposable: ValidationStateTrue,
			blackList:  ValidationStateFalse,
		},
		{
			email:      "iub65391@bcaoo.com",
			free:       ValidationStateFalse,
			disposable: ValidationStateTrue,
			blackList:  ValidationStateFalse,
		},
		{
			email:      "abuse@example.com",
			free:       ValidationStateFalse,
			disposable: ValidationStateFalse,
			blackList:  ValidationStateTrue,
		},
		{
			email: "fail@iub65391@bcaoo.com",
			fail:  true,
		},
		{
			email: "fail@localhost",
			fail:  true,
		},
		{
			email: "fail@localhost.invalidtld",
			fail:  true,
		},
		{
			email: "fa illong@gmail.com",
			fail:  true,
		},
		{
			email: "fa il@mysite.com",
			fail:  true,
		},
		{
			email: ".fail@mysite.com",
			fail:  true,
		},
		{
			email: "fail.@mysite.com",
			fail:  true,
		},
		{
			email: ".fail@gmail.com",
			fail:  true,
		},
		{
			email: "fail.@gmail.com",
			fail:  true,
		},
		{
			email: "fail@gmail.com",
			fail:  true,
		},
		{
			email: "fail+extra@gmail.com",
			fail:  true,
		},
		{
			email: "faillong+ex tra@gmail.com",
			fail:  true,
		},
		{
			email: "faillong+ex,tra@gmail.com",
			fail:  true,
		},
		{
			email: "fail.+extra@gmail.com",
			fail:  true,
		},
		{
			email: "fail<user>@gmail.com",
			fail:  true,
		},
		{
			email: strings.Repeat("a", 65) + "@mydomain.com",
			fail:  true,
		},
		{
			email: "valid@mydomain" + strings.Repeat("a", 255) + ".com",
			fail:  true,
		},
	}
)

func TestValidate(t *testing.T) {
	for _, tf := range testFixtures {
		res, err := Validate(tf.email)
		if err != nil {
			require.True(t, tf.fail)
			continue
		}
		assert.Equal(t, tf.free, res.FreeProvider)
		assert.Equal(t, tf.disposable, res.Disposable)
		assert.Equal(t, tf.blackList, res.BlackList)
	}
}

func TestValidateMX(t *testing.T) {
	chk := CheckMX(0, false)
	res, err := Validate("validemail@gmail.com", chk)
	require.Error(t, err)

	res, err = Validate("email@google.com", CheckMX(time.Second, false))
	require.NoError(t, err)
	assert.Equal(t, ValidationStateTrue, res.MXValidation)
	assert.Equal(t, ValidationStateFalse, res.Disposable)
	assert.Equal(t, ValidationStateFalse, res.FreeProvider)

	res, err = Validate("email@ifsomeonebuythisdomainandrunitsomewherethistestfails.com", CheckMX(time.Second, true))
	require.NoError(t, err)
	assert.Equal(t, ValidationStateFalse, res.MXValidation)
	assert.Equal(t, ValidationStateFalse, res.Disposable)
	assert.Equal(t, ValidationStateFalse, res.FreeProvider)
}

func TestJSONResult(t *testing.T) {
	res := ValidationResult{
		FreeProvider: ValidationStateNotChecked,
		Disposable:   ValidationStateFalse,
		MXValidation: ValidationStateTrue,
	}

	b, err := json.Marshal(res)
	require.NoError(t, err)

	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	require.NoError(t, err)

	assert.Equal(t, map[string]interface{}{
		"free_provider": nil,
		"disposable":    false,
		"mx_validation": true,
		"black_list":    nil,
	}, m)

	res = ValidationResult{
		FreeProvider: -1,
		Disposable:   -2,
		MXValidation: -3,
	}

	_, err = json.Marshal(&res)
	require.Error(t, err)

}
