package db

type Source struct {
	source  string
	like    bool
	comment bool
}

type LikedIP struct {
	ip string
}

type Like struct {
	uri   string
	count int
}
