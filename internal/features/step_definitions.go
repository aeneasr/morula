package features

import . "github.com/gucumber/gucumber"
import "fmt"

func init() {

	When(`^running "(.+?)"$`, func(s1 string) {
		fmt.Println("testing")
	})

}
