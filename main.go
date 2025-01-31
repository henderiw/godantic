//go:generate go run pkg/genvalidate/generator.go ./apis

package main

import (
	networkv1alpha1 "github.com/henderiw/godantic/apis/kubenet/apis/network/v1alpha1"
)

func main() {
	//validategenerator := genvalidate.NewGenerator("./apis")
	//validategenerator.Generate()

	x := networkv1alpha1.Dummy(1)
	if err := x.Validate(); err != nil {
		panic(err)
	}

}
