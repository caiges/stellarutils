stellarutils
============

A collection of helper methods for working with the stellar network and related functionality.

## Federated Users

```go
package main
import (
	"fmt"
	"github.com/caiges/stellarutils"
)
func main() {
	address, err := stellarutils.ResolveFederationUser("caiges@stellar.org")
	if err != nil {
		fmt.Printf("Problem resolving address: %v", err)
	}
	fmt.Println(address.FederationJSON.DestinationAddress)
}
```
