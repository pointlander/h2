// Copyright 2021 The h2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/itsubaki/q"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Coefficients are the coefficients for the quantum simulation
type Coefficients struct {
	R            float64
	One          float64
	Z0           float64
	Z1           float64
	Z0Z1         float64
	X0X1         float64
	Y0Y1         float64
	t0           float64
	Ordering     string
	TrotterError float64
}

var Coeff = [...]Coefficients{
	Coefficients{0.20, 2.8489, 0.5678, -1.4508, 0.6799, 0.0791, 0.0791, 1.500, "Z0 · X0X1 · Z1 · Y0Y1", 0.0124},
	Coefficients{0.25, 2.1868, 0.5449, -1.2870, 0.6719, 0.0798, 0.0798, 1.590, "Z0 · Y0Y1 · Z1 · X0X1", 0.0521},
	Coefficients{0.30, 1.7252, 0.5215, -1.1458, 0.6631, 0.0806, 0.0806, 1.770, "X0X1 · Z0 · Y0Y1 · Z1", 0.0111},
	Coefficients{0.35, 1.3827, 0.4982, -1.0226, 0.6537, 0.0815, 0.0815, 2.080, "Z0 · X0X1 · Z1 · Y0Y1", 0.0368},
	Coefficients{0.40, 1.1182, 0.4754, -0.9145, 0.6438, 0.0825, 0.0825, 2.100, "Z0 · X0X1 · Z1 · Y0Y1", 0.0088},
	Coefficients{0.45, 0.9083, 0.4534, -0.8194, 0.6336, 0.0835, 0.0835, 2.310, "X0X1 · Z0 · Y0Y1 · Z1", 0.0141},
	Coefficients{0.50, 0.7381, 0.4325, -0.7355, 0.6233, 0.0846, 0.0846, 2.580, "Z0 · X0X1 · Z1 · Y0Y1", 0.0672},
	Coefficients{0.55, 0.5979, 0.4125, -0.6612, 0.6129, 0.0858, 0.0858, 2.700, "Z0 · X0X1 · Z1 · Y0Y1", 0.0147},
	Coefficients{0.60, 0.4808, 0.3937, -0.5950, 0.6025, 0.0870, 0.0870, 2.250, "Z0 · X0X1 · Z1 · Y0Y1", 0.0167},
	Coefficients{0.65, 0.3819, 0.3760, -0.5358, 0.5921, 0.0883, 0.0883, 3.340, "Z1 · X0X1 · Z0 · Y0Y1", 0.0175},
	Coefficients{0.70, 0.2976, 0.3593, -0.4826, 0.5818, 0.0896, 0.0896, 0.640, "Z0 · Y0Y1 · Z1 · X0X1", 0.0171},
	Coefficients{0.75, 0.2252, 0.3435, -0.4347, 0.5716, 0.0910, 0.0910, 0.740, "Z0 · Y0Y1 · Z1 · X0X1", 0.0199},
	Coefficients{0.80, 0.1626, 0.3288, -0.3915, 0.5616, 0.0925, 0.0925, 0.790, "Z0 · Y0Y1 · Z1 · X0X1", 0.0291},
	Coefficients{0.85, 0.1083, 0.3149, -0.3523, 0.5518, 0.0939, 0.0939, 3.510, "Z0 · X0X1 · Z1 · Y0Y1", 0.0254},
	Coefficients{0.90, 0.0609, 0.3018, -0.3168, 0.5421, 0.0954, 0.0954, 3.330, "Z0 · X0X1 · Z1 · Y0Y1", 0.0283},
	Coefficients{0.95, 0.0193, 0.2895, -0.2845, 0.5327, 0.0970, 0.0970, 4.090, "X0X1 · Z0 · Y0Y1 · Z1", 0.0328},
	Coefficients{1.00, -0.0172, 0.2779, -0.2550, 0.5235, 0.0986, 0.0986, 4.360, "Z0 · X0X1 · Z1 · Y0Y1", 0.0362},
	Coefficients{1.05, -0.0493, 0.2669, -0.2282, 0.5146, 0.1002, 0.1002, 4.650, "Z1 · X0X1 · Z0 · Y0Y1", 0.0405},
	Coefficients{1.10, -0.0778, 0.2565, -0.2036, 0.5059, 0.1018, 0.1018, 4.280, "Z1 · X0X1 · Z0 · Y0Y1", 0.0243},
	Coefficients{1.15, -0.1029, 0.2467, -0.1810, 0.4974, 0.1034, 0.1034, 5.510, "Z0 · X0X1 · Z1 · Y0Y1", 0.0497},
	Coefficients{1.20, -0.1253, 0.2374, -0.1603, 0.4892, 0.1050, 0.1050, 5.950, "Z0 · Y0Y1 · Z1 · X0X1", 0.0559},
	Coefficients{1.25, -0.1452, 0.2286, -0.1413, 0.4812, 0.1067, 0.1067, 6.360, "X0X1 · Z1 · Y0Y1 · Z0", 0.0585},
	Coefficients{1.30, -0.1629, 0.2203, -0.1238, 0.4735, 0.1083, 0.1083, 0.660, "Z1 · X0X1 · Z0 · Y0Y1", 0.0905},
	Coefficients{1.35, -0.1786, 0.2123, -0.1077, 0.4660, 0.1100, 0.1100, 9.810, "Z0 · X0X1 · Z1 · Y0Y1", 0.0694},
	Coefficients{1.40, -0.1927, 0.2048, -0.0929, 0.4588, 0.1116, 0.1116, 9.930, "Z0 · X0X1 · Z1 · Y0Y1", 0.0755},
	Coefficients{1.45, -0.2053, 0.1976, -0.0792, 0.4518, 0.1133, 0.1133, 5.680, "Y0Y1 · Z0 · X0X1 · Z1", 0.0142},
	Coefficients{1.50, -0.2165, 0.1908, -0.0666, 0.4451, 0.1149, 0.1149, 10.200, "Z1 · X0X1 · Z0 · Y0Y1", 0.0885},
	Coefficients{1.55, -0.2265, 0.1843, -0.0549, 0.4386, 0.1165, 0.1165, 9.830, "Z0 · X0X1 · Z1 · Y0Y1", 0.0917},
	Coefficients{1.60, -0.2355, 0.1782, -0.0442, 0.4323, 0.1181, 0.1181, 8.150, "Z0 · Y0Y1 · Z1 · X0X1", 0.0416},
	Coefficients{1.65, -0.2436, 0.1723, -0.0342, 0.4262, 0.1196, 0.1196, 8.240, "X0X1 · Z0 · Y0Y1 · Z1", 0.0488},
	Coefficients{1.70, -0.2508, 0.1667, -0.0251, 0.4204, 0.1211, 0.1211, 0.520, "Z1 · X0X1 · Z0 · Y0Y1", 0.0450},
	Coefficients{1.75, -0.2573, 0.1615, -0.0166, 0.4148, 0.1226, 0.1226, 0.520, "Z0 · Y0Y1 · Z1 · X0X1", 0.0509},
	Coefficients{1.80, -0.2632, 0.1565, -0.0088, 0.4094, 0.1241, 0.1241, 1.010, "Z0 · X0X1 · Z1 · Y0Y1", 0.0663},
	Coefficients{1.85, -0.2684, 0.1517, -0.0015, 0.4042, 0.1256, 0.1256, 0.530, "Z1 · X0X1 · Z0 · Y0Y1", 0.0163},
	Coefficients{1.90, -0.2731, 0.1472, 0.0052, 0.3992, 0.1270, 0.1270, 1.090, "X0X1 · Z0 · Z1 · Y0Y1", 0.0017},
	Coefficients{1.95, -0.2774, 0.1430, 0.0114, 0.3944, 0.1284, 0.1284, 0.610, "X0X1 · Z1 · Z0 · Y0Y1", 0.0873},
	Coefficients{2.00, -0.2812, 0.1390, 0.0171, 0.3898, 0.1297, 0.1297, 1.950, "Z1 · Z0 · X0X1 · Y0Y1", 0.0784},
	Coefficients{2.05, -0.2847, 0.1352, 0.0223, 0.3853, 0.1310, 0.1310, 4.830, "X0X1 · Y0Y1 · Z0 · Z1", 0.0947},
	Coefficients{2.10, -0.2879, 0.1316, 0.0272, 0.3811, 0.1323, 0.1323, 1.690, "Y0Y1 · X0X1 · Z0 · Z1", 0.0206},
	Coefficients{2.15, -0.2908, 0.1282, 0.0317, 0.3769, 0.1335, 0.1335, 0.430, "X0X1 · Y0Y1 · Z0 · Z1", 0.0014},
	Coefficients{2.20, -0.2934, 0.1251, 0.0359, 0.3730, 0.1347, 0.1347, 1.750, "Z0 · Z1 · X0X1 · Y0Y1", 0.0107},
	Coefficients{2.25, -0.2958, 0.1221, 0.0397, 0.3692, 0.1359, 0.1359, 11.500, "X0X1 · Z1 · Z0 · Y0Y1", 0.0946},
	Coefficients{2.30, -0.2980, 0.1193, 0.0432, 0.3655, 0.1370, 0.1370, 0.420, "Z0 · Z1 · X0X1 · Y0Y1", 0.0370},
	Coefficients{2.35, -0.3000, 0.1167, 0.0465, 0.3620, 0.1381, 0.1381, 0.470, "Z1 · Z0 · Y0Y1 · X0X1", 0.0762},
	Coefficients{2.40, -0.3018, 0.1142, 0.0495, 0.3586, 0.1392, 0.1392, 10.100, "X0X1 · Z1 · Z0 · Y0Y1", 0.0334},
	Coefficients{2.45, -0.3035, 0.1119, 0.0523, 0.3553, 0.1402, 0.1402, 11.200, "Z0 · Z1 · X0X1 · Y0Y1", 0.0663},
	Coefficients{2.50, -0.3051, 0.1098, 0.0549, 0.3521, 0.1412, 0.1412, 0.580, "Z0 · Y0Y1 · X0X1 · Z1", 0.0296},
	Coefficients{2.55, -0.3066, 0.1078, 0.0572, 0.3491, 0.1422, 0.1422, 11.000, "Z0 · Z1 · X0X1 · Y0Y1", 0.0550},
	Coefficients{2.60, -0.3079, 0.1059, 0.0594, 0.3461, 0.1432, 0.1432, 11.000, "Z0 · X0X1 · Y0Y1 · Z1", 0.0507},
	Coefficients{2.65, -0.3092, 0.1042, 0.0614, 0.3433, 0.1441, 0.1441, 11.040, "Z1 · X0X1 · Y0Y1 · Z0", 0.0490},
	Coefficients{2.70, -0.3104, 0.1026, 0.0632, 0.3406, 0.1450, 0.1450, 0.400, "Z0 · Z1 · Y0Y1 · X0X1", 0.0471},
	Coefficients{2.75, -0.3115, 0.1011, 0.0649, 0.3379, 0.1458, 0.1458, 0.450, "Y0Y1 · Z0 · Z1 · X0X1", 0.0061},
	Coefficients{2.80, -0.3125, 0.0997, 0.0665, 0.3354, 0.1467, 0.1467, 0.950, "Z0 · Y0Y1 · X0X1 · Z1", 0.0368},
	Coefficients{2.85, -0.3135, 0.0984, 0.0679, 0.3329, 0.1475, 0.1475, 10.600, "Z0 · X0X1 · Y0Y1 · Z1", 0.0324},
}

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

	sim := func(coeff int, theta float64) float64 {
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
		coefficients := Coeff[coeff]
		return coefficients.One +
			e["zi"]*coefficients.Z0 +
			e["iz"]*coefficients.Z1 +
			e["zz"]*coefficients.Z0Z1 +
			e["xx"]*coefficients.X0X1 +
			e["yy"]*coefficients.Y0Y1
	}

	getMinEnergy := func(coeff int) float64 {
		points := make(plotter.XYs, 0, 8)
		min := math.MaxFloat64
		for v := 0.0; v < 2*math.Pi; v += .1 {
			result := sim(coeff, v)
			if result < min {
				min = result
			}
			points = append(points, plotter.XY{X: v, Y: result})
			fmt.Println(result)
		}

		if coeff == 11 {
			p := plot.New()

			p.Title.Text = "energy vs theta"
			p.X.Label.Text = "theta"
			p.Y.Label.Text = "energy"

			scatter, err := plotter.NewScatter(points)
			if err != nil {
				panic(err)
			}
			scatter.GlyphStyle.Radius = vg.Length(1)
			scatter.GlyphStyle.Shape = draw.CircleGlyph{}
			p.Add(scatter)

			err = p.Save(8*vg.Inch, 8*vg.Inch, "energy.png")
			if err != nil {
				panic(err)
			}
		}

		return min
	}

	points := make(plotter.XYs, 0, 8)
	min, distance := math.MaxFloat64, 0.0
	for i, coefficients := range Coeff {
		energy := getMinEnergy(i)
		if energy < min {
			min, distance = energy, coefficients.R
		}
		points = append(points, plotter.XY{X: coefficients.R, Y: energy})
	}

	p := plot.New()

	p.Title.Text = "energy vs distance"
	p.X.Label.Text = "distance"
	p.Y.Label.Text = "energy"

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Radius = vg.Length(1)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(scatter)

	err = p.Save(8*vg.Inch, 8*vg.Inch, "min_energy.png")
	if err != nil {
		panic(err)
	}

	fmt.Println(min, distance)
}
