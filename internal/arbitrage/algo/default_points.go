package algo

import (
	"math/big"
)

var criticalPoints []*big.Int

func init() {
	values := []int{
		1, 5, 8, 10, 12, 15, 18, 20, 25, 30, 35, 40, 45, 50, 55, 60, 70, 75, 80, 85, 90, 95, 100, 105, 110,
		115, 125, 135, 145, 150, 155, 160, 170, 180, 190, 200, 210, 220, 240, 250, 300, 350, 400, 450,
		500, 550, 600, 700, 800, 900, 1000, 1100, 1200, 1300, 1400, 1500, 1600, 1700, 1800, 1900,
		2000, 2100, 2200, 2300, 2500, 2700, 2900, 3100, 3300, 3500, 3700, 3900, 4100, 4300, 4500,
		4700, 4900, 5100, 5300, 5500, 5700, 5900, 6100, 6300, 6500, 6700, 6900, 7100, 7300, 7500, 7700, 7900, 8100, 8300, 8500,
		8700, 9000, 9200, 9500, 9700,
		10000, 10500, 11000, 11500, 12000,
	}

	criticalPoints = make([]*big.Int, len(values))

	tenToThePowerOf16 := new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil)

	for i, v := range values {
		tmp := new(big.Int).Mul(big.NewInt(int64(v)), tenToThePowerOf16)

		criticalPoints[i] = tmp
	}
}

func PrepareNewPointsWithMerge(min *big.Int, max *big.Int, splitParam *big.Int, valueToMerge *big.Int) ([]*big.Int, *big.Int, int) {
	var points []*big.Int

	diff := new(big.Int).Sub(max, min)
	bestValueIndex := 0
	merged := false
	if valueToMerge == nil {
		merged = true
	}
	if diff.Cmp(zero) == 0 {
		return nil, zero, 0
	}

	for i := int64(0); ; i++ {
		mul := new(big.Int).Mul(diff, big.NewInt(i))
		offset := new(big.Int).Quo(mul, splitParam)
		point := new(big.Int).Add(min, offset)

		if !merged {
			// to avoid duplication
			if valueToMerge.Cmp(point) == 0 {
				merged = true
				bestValueIndex = int(i)
			}
			if valueToMerge.Cmp(point) == -1 {
				bestValueIndex = int(i)
				points = append(points, valueToMerge)
				merged = true
			}
		}

		if point.Cmp(max) == 1 {
			break
		}

		points = append(points, point)
	}

	return points, diff, bestValueIndex
}
