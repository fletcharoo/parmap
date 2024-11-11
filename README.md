# parmap
[parmap](https://github.com/fletcharoo/parmap) offers a simple generic implementation of parallel map functionality.

Features:
* Parallel execution
* Generic
* Deterministic
* Verbose error handling
* No race conditions

TODO:
* Implement a way to set max goroutines for the library to start

## Installation
`go get github.com/fletcharoo/parmap`

## Usage
Example:
```go
import (
  "fmt"
  "github.com/fletcharoo/parmap"
)

func main() {
  inputs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
  do := func(i int) (int, error) {
    return i * 5, nil
  }

  result, err := parmap.Do(inputs, do)

  if err != nil {
    panic(fmt.Sprintf("Unepxected error: %s", err))
  }

  fmt.Println(result)
}
```
Output:
```
[5 10 15 20 25 30 35 40 45 50]
```

## Contributing
I'm always open to constructive criticism and help, so if you do wish to contribute to this repo, please abide by the following process:
* Create an issue that describes the bug/feature
* If you wish to fix/implement this yourself, please create a PR and link to the issue you created (ensure that all code you update/add is well tested)
