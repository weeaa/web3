package db

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// I know it's not the best testing function
// but I have 1 table so it doesn't really matter.
func TestDatabase(t *testing.T) {

	pg, err := New()
	if err != nil {
		assert.Error(t, err)
	}

	defer pg.RemoveUser("0xTest")

	u := &User{
		BaseAddress:     "0xTest",
		Status:          "whaleTest",
		TwitterName:     "wjhTest",
		TwitterUsername: "whj_Test",
		TwitterURL:      "https://x.com/weeaa",
	}

	if err = pg.InsertUser(u); err != nil {
		assert.Error(t, err)
	}

	user, err := pg.GetUser(u.BaseAddress)
	if err != nil {
		assert.Error(t, err)
	}

	if user.TwitterURL != u.TwitterURL {
		assert.Error(t, fmt.Errorf("expected %s, got %s", u.TwitterURL, user.TwitterURL))
	}

	if err = pg.RemoveUser(u.BaseAddress); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}
