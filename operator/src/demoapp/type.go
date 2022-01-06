package demoapp

import (
	"crypto/sha512"
	"fmt"
)

type DemoApp struct {
	Name      string
	Namespace string
}

func NewDemoApp(name string, namespace string) *DemoApp {
	if name == "" {
		name = defaultName
	}

	return &DemoApp{
		Name:      name,
		Namespace: namespace,
	}
}

func calculateHash(object interface{}) string {
	hash := sha512.New()
	hash.Write([]byte(fmt.Sprintf("%v", object)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
