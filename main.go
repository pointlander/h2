// Copyright 2021 The h2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/itsubaki/q"
)

func main() {
	rand.Seed(1)

	type System struct {
		qsim *q.Q
		q0   q.Qubit
		q1   q.Qubit
	}

	base := func(theta float64) System {
		qsim := q.New()
		q0 := qsim.Zero()
		q1 := qsim.Zero()
		qsim.RX(math.Pi, q0)
		qsim.RY(math.Pi/2, q1)
		qsim.RX(-math.Pi/2, q0)
		qsim.CNOT(q1, q0)
		qsim.RZ(theta, q0)
		qsim.CNOT(q1, q0)
		qsim.RX(math.Pi/2, q0)
		qsim.RY(-math.Pi/2, q1)
		return System{
			qsim: qsim,
			q0:   q0,
			q1:   q1,
		}
	}
	xx := func(theta float64) System {
		system := base(theta)
		system.qsim.RY(-math.Pi/2, system.q0)
		system.qsim.RY(-math.Pi/2, system.q1)
		return system
	}
	yy := func(theta float64) System {
		system := base(theta)
		system.qsim.RX(math.Pi/2, system.q0)
		system.qsim.RX(math.Pi/2, system.q1)
		return system
	}

	sim := func(theta float64) float64 {
		sz := base(theta)
		sxx := xx(theta)
		syy := yy(theta)
		signs := map[string]map[string]float64{
			"zi": {
				"00": 1,
				"01": 1,
				"10": -1,
				"11": -1,
			},
			"iz": {
				"00": 1,
				"01": -1,
				"10": 1,
				"11": -1,
			},
			"zz": {
				"00": 1,
				"01": -1,
				"10": -1,
				"11": 1,
			},
			"xx": {
				"00": 1,
				"01": -1,
				"10": -1,
				"11": 1,
			},
			"yy": {
				"00": 1,
				"01": -1,
				"10": -1,
				"11": 1,
			},
		}
		e := map[string]float64{}
		for _, state := range sz.qsim.State() {
			v := e["zi"]
			v += signs["zi"][state.BinaryString[0]] * state.Probability
			e["zi"] = v

			v = e["iz"]
			v += signs["iz"][state.BinaryString[0]] * state.Probability
			e["iz"] = v

			v = e["zz"]
			v += signs["zz"][state.BinaryString[0]] * state.Probability
			e["zz"] = v
		}
		for _, state := range sxx.qsim.State() {
			v := e["xx"]
			v += signs["xx"][state.BinaryString[0]] * state.Probability
			e["xx"] = v
		}
		for _, state := range syy.qsim.State() {
			v := e["yy"]
			v += signs["yy"][state.BinaryString[0]] * state.Probability
			e["yy"] = v
		}
		sum := 0.0
		for _, value := range e {
			sum += value
		}
		return sum
	}

	for v := 0.0; v < 2*math.Pi; v += .1 {
		fmt.Println(sim(v))
	}
}
