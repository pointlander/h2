// Copyright 2021 The h2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"math/rand"
  "fmt"

	"github.com/itsubaki/q"
)

func main() {
	rand.Seed(1)

	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	qsim.RX(math.Pi, q0)
	qsim.RY(math.Pi/2, q1)
	qsim.RX(-math.Pi/2, q0)
	qsim.CNOT(q1, q0)
	qsim.RZ(2, q0)
	qsim.CNOT(q1, q0)
	qsim.RX(math.Pi/2, q0)
  qsim.RY(-math.Pi/2, q1)
  max, binary := 0.0, []string{}
  for _, state := range qsim.State() {
    if state.Probability > max {
      max, binary = state.Probability, state.BinaryString
    }
  }
  fmt.Println(binary)
}
