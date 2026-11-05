average
[![TravisCI](https://travis-ci.org/prep/average.svg?branch=master)](https://travis-ci.org/prep/average.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/prep/average)](https://goreportcard.com/report/github.com/prep/average)
[![GoDoc](https://godoc.org/github.com/prep/average?status.svg)](https://godoc.org/github.com/prep/average)
=======
This stupidly named Go package contains a single struct that is used to implement counters on a sliding time window.

Usage
-----
```go

import (
    "fmt"

    "github.com/prep/average"
)

func main() {
    // Create a SlidingWindow that has a window of 15 minutes, with a
    // granulity of 1 minute.
    sw := average.MustNew(15 * time.Minute, time.Minute)
    defer sw.Stop()

    // Do some work.
    sw.Add(15)
    // Do some more work.
    sw.Add(22)
    // Do even more work.
    sw.Add(22)

    fmt.Printf("Average of last  1m: %f\n", sw.Average(time.Minute)
    fmt.Printf("Average of last  5m: %f\n", sw.Average(5 * time.Minute)
    fmt.Printf("Average of last 15m: %f\n\n", sw.Average(15 * time.Minute)

    total, numSamples := sw.Total(15 * time.Minute)
    fmt.Printf("Counter has a total of %d over %d samples", total, numSamples)
}
```

License
-------
This software is created for MessageBird B.V. and distributed under the BSD-style license found in the LICENSE file.
