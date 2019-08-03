package emailvalidator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixture struct {
	email      string
	free       bool
	disposable bool
	fail       bool
}

var (
	testFixtures = []fixture{
		{
			email:      "test.with.dot@gmail.com",
			free:       true,
			disposable: false,
		},
		{
			email:      "test.with.dot+extra@gmail.com",
			free:       true,
			disposable: false,
		},
		{
			email:      "test@things.10mail.org",
			free:       false,
			disposable: true,
		},
		{
			email:      "test@things.more.10mail.org",
			free:       false,
			disposable: true,
		},
		{
			email:      "iub65391@bcaoo.com",
			free:       false,
			disposable: true,
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
	}
)

func TestValidate(t *testing.T) {
	for _, tf := range testFixtures {
		free, disposable, err := Validate(tf.email)
		if err != nil {
			require.True(t, tf.fail)
			continue
		}
		assert.Equal(t, tf.free, free)
		assert.Equal(t, tf.disposable, disposable)
	}
}

func TestValidateMX(t *testing.T) {
	chk := CheckMX(0)
	_, _, err := Validate("validemail@gmail.com", chk)
	require.Error(t, err)

	chk = CheckMX(time.Second)

	f, d, err := Validate("email@google.com", chk)
	require.NoError(t, err)
	assert.False(t, d)
	assert.False(t, f)

	f, d, err = Validate("email@ifsomeonebuythisdomainandrunitsomewherethistestfails.com", chk)
	require.Error(t, err)
	assert.False(t, d)
	assert.False(t, f)

}
