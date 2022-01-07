package docker

type manifest struct {
	conf config
}

type config struct {
	digest string
}
