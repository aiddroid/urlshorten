package urlshorten

type Storage interface {
	Shorten(url string, expire int) (string, error)
	UnShorten(eid string) (string, error)
	ShortenInfo(eid string) (interface{}, error)
}
