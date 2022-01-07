package docker

type manifest struct {
	Conf config `json:"config"`
}

type config struct {
	Digest string `json:"digest"`
}
