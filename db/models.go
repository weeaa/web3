package db

type (

	// User is 'user' Table
	User struct {
		BaseAddress     string //base_address text
		Status          string //status text
		TwitterUsername string //twitter_username text
		TwitterName     string //twitter_name text
		TwitterURL      string //twitter_url text
		UserId          int    //user_id integer
	}
)
