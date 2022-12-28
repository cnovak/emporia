# emporia-api-go


`emporia-api-go` is Go API Client for the [Emporia](https://www.emporiaenergy.com/) Vue energy monitoring system API. This API is not officially supported by Emporia.

The library can be invoked directly to pull back some basic info but requires your email and password to be added to a keys.json file, which is then replaced with the access tokens.

[Mike Marvin](https://github.com/magico13)'s [API documentation](https://github.com/magico13/PyEmVue/blob/master/api_docs.md) can help to understand the underlaying API, or use the OpenAPI file in the `/docs` directory to understand the API.

## Installation

`go get -u github.com/cnovak/emporia`

## Usage
You need a username and password to initialize the client:

```
package main

import (
	"fmt"
	"os"

	emporia "github.com/cnovak/emporia"
)

func main() {
	username := os.Getenv("EMPORIA_USERNAME")
	password := os.Getenv("EMPORIA_PASSWORD")

	client, err := emporia.NewClient(username, password)
	if err != nil {
		panic(err)
	}

	devices, err := client.GetDevices()
	if err != nil {
		panic(err)
	}
	fmt.Printf("============== \nDevices:\n============== \n%+v\n\n", devices)

}


```

## Credits
Thanks to [Mike Marvin](https://github.com/magico13) for his work in the [PyEmVue project](https://github.com/magico13/PyEmVue) documenting the Emporia Energy API.

## Disclaimer
This project is not affiliated with or endorsed by Emporia Energy.