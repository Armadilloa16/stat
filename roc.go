// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"sort"
)

// ROC returns paired false positive rate (FPR) and true positive rate
// (TPR) values corresponding to cutoffs, i.e. particular points on the
// receiver operator characteristic (ROC) curve obtained when y is
// treated as a binary classifier for classes with weights.
//
// The input y and cutoffs must be sorted, and values in y must correspond
// to values in classes and weights. SortWeightedLabeled can be used to
// sort y together with classes and weights.
//
// For a given cutoff value, observations corresponding to entries in y
// greater than the cutoff value are classified as false, while those
// less than (or equal to) the cutoff value are classified as true. These
// assigned class labels are compared with the true values in the classes
// slice and used to calculate the FPR and TPR.
//
// If weights is nil, all weights are treated as 1.
//
// If cutoffs is nil or empty all possible cutoffs are calculated,
// resulting in fpr and tpr having length one greater than the number of
// unique values in y. Otherwise fpr and tpr will be returned with the
// same length as cutoffs. EquallySpaced can be used to generate
// equally spaced cutoffs.
//
// More details about ROC curves are available at
// https://en.wikipedia.org/wiki/Receiver_operating_characteristic
func ROC(cutoffs, y []float64, classes []bool, weights []float64) (tpr, fpr []float64) {
	if len(y) != len(classes) {
		panic("stat: slice length mismatch")
	}
	if weights != nil && len(y) != len(weights) {
		panic("stat: slice length mismatch")
	}
	if !sort.Float64sAreSorted(y) {
		panic("stat: input must be sorted")
	}
	if !sort.Float64sAreSorted(cutoffs) {
		panic("stat: input must be sorted")
	}
	if len(y) == 0 {
		return nil, nil
	}
	var bin int
	if len(cutoffs) == 0 {
		cutoffs = make([]float64, len(y)+1)
		cutoffs[0] = math.Nextafter(y[0], y[0]-1)
		// Choose all possible cutoffs but remove duplicate values
		// in y.
		for i, u := range y {
			if i != 0 && u != y[i-1] {
				bin++
			}
			cutoffs[bin+1] = u
		}
		cutoffs = cutoffs[0:(bin + 2)]
	}

	tpr = make([]float64, len(cutoffs))
	fpr = make([]float64, len(cutoffs))
	bin = 0
	var nPos, nNeg float64
	for i, u := range classes {
		// Update the bin until it matches the next y value
		// (skip empty bins).
		for (bin < len(cutoffs)) && (y[i] > cutoffs[bin]) {
			if bin == (len(cutoffs) - 1) {
				break
			}
			bin++
			tpr[bin] = tpr[bin-1]
			fpr[bin] = fpr[bin-1]
		}
		var posWeight, negWeight float64 = 0, 1
		if weights != nil {
			negWeight = weights[i]
		}
		if u {
			posWeight, negWeight = negWeight, posWeight
		}
		nPos += posWeight
		nNeg += negWeight
		if y[i] <= cutoffs[bin] {
			tpr[bin] += posWeight
			fpr[bin] += negWeight
		}
	}

	invNeg := 1 / nNeg
	invPos := 1 / nPos
	for i := range tpr {
		tpr[i] *= invPos
		fpr[i] *= invNeg
	}

	return tpr, fpr
}

// EquallySpaced returns n equally spaced values,
// beggining eps less than min and ending at max. n must be
// at least two, and min must not be greater than max.
func EquallySpaced(min, max float64, n int) (cutoffs []float64) {
	if max < min {
		panic("stat: max < min")
	}
	if n < 2 {
		panic("stat: n too small")
	}

	cutoffs = make([]float64, n)
	w := (max - min) / float64(n-1)
	for i := range cutoffs {
		cutoffs[i] = min + w*float64(i)
	}
	cutoffs[0] = math.Nextafter(min, min-1)

	return cutoffs
}
