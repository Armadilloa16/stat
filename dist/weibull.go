// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dist

import (
	"math"
	"math/cmplx"
	"math/rand"
)

// Weibull distribution. Valid range for x is [0,+∞).
type Weibull struct {
	// Shape parameter of the distribution. A value of 1 represents
	// the exponential distribution. A value of 2 represents the
	// Rayleigh distribution. Valid range is (0,+∞).
	K float64
	// Scale parameter of the distribution. Valid range is (0,+∞).
	Lambda float64
	// Source of random numbers
	Source *rand.Rand
}

// CDF computes the value of the cumulative density function at x.
func (w Weibull) CDF(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return 1 - cmplx.Abs(cmplx.Exp(w.LogCDF(x)))
	}
}

// DLogProbDX returns the derivative of the log of the probability with
// respect to the input x.
//
// Special cases are:
//  DLogProbDX(0) = NaN
func (w Weibull) DLogProbDX(x float64) float64 {
	if x > 0 {
		return (-w.K*math.Pow(x/w.Lambda, w.K) + w.K - 1) / x
	}
	if x < 0 {
		return 0
	}
	return math.NaN()
}

// DLogProbDParam returns the derivative of the log of the probability with
// respect to the parameters of the distribution. The deriv slice must have length
// equal to the number of parameters of the distribution.
//
// The order is ∂LogProb / ∂K and then ∂LogProb / ∂λ.
//
// Special cases are:
//  The derivative at 0 is NaN.
func (w Weibull) DLogProbDParam(x float64, deriv []float64) {
	if len(deriv) != w.NumParameters() {
		panic("weibull: slice length mismatch")
	}
	if x > 0 {
		deriv[0] = 1/w.K + math.Log(x) - math.Log(w.Lambda) - (math.Log(x)-math.Log(w.Lambda))*math.Pow(x/w.Lambda, w.K)
		deriv[1] = (w.K * (math.Pow(x/w.Lambda, w.K) - 1)) / w.Lambda
		return
	}
	if x < 0 {
		deriv[0] = 0
		deriv[1] = 0
		return
	}
	deriv[0] = math.NaN()
	deriv[0] = math.NaN()
	return
}

// Entropy returns the entropy of the distribution.
func (w Weibull) Entropy() float64 {
	return eulerGamma*(1-1/w.K) + math.Log(w.Lambda/w.K) + 1
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (w Weibull) ExKurtosis() float64 {
	return (-6*w.gammaIPow(1, 4) + 12*w.gammaIPow(1, 2)*math.Gamma(1+2/w.K) - 3*w.gammaIPow(2, 2) - 4*math.Gamma(1+1/w.K)*math.Gamma(1+3/w.K) + math.Gamma(1+4/w.K)) / math.Pow(math.Gamma(1+2/w.K)-w.gammaIPow(1, 2), 2)
}

// gammIPow is a shortcut for computing the gamma function to a power.
func (w Weibull) gammaIPow(i, pow float64) float64 {
	return math.Pow(math.Gamma(1+i/w.K), pow)
}

// LogCDF computes the value of the log of the cumulative density function at x.
func (w Weibull) LogCDF(x float64) complex128 {
	if x < 0 {
		return 0
	} else {
		return cmplx.Log(-1) + complex(-math.Pow(x/w.Lambda, w.K), 0)
	}
}

// LogProb computes the natural logarithm of the value of the probability
// density function at x. Zero is returned if x is less than zero.
//
// Special cases occur when x == 0, and the result depends on the shape
// parameter as follows:
//  If 0 < K < 1, LogProb returns +Inf.
//  If K == 1, LogProb returns 0.
//  If K > 1, LogProb returns -Inf.
func (w Weibull) LogProb(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return math.Log(w.K) - math.Log(w.Lambda) + (w.K-1)*(math.Log(x)-math.Log(w.Lambda)) - math.Pow(x/w.Lambda, w.K)
	}
}

// Survival returns the log of the survival function (complementary CDF) at x.
func (w Weibull) LogSurvival(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return -math.Pow(x/w.Lambda, w.K)
	}
}

// Mean returns the mean of the probability distribution.
func (w Weibull) Mean() float64 {
	return w.Lambda * math.Gamma(1+1/w.K)
}

// Median returns the median of the normal distribution.
func (w Weibull) Median() float64 {
	return w.Lambda * math.Pow(ln2, 1/w.K)
}

// Mode returns the mode of the normal distribution.
//
// The mode is NaN in the special case where the K (shape) parameter
// is less than 1.
func (w Weibull) Mode() float64 {
	if w.K > 1 {
		return w.Lambda * math.Pow((w.K-1)/w.K, 1/w.K)
	} else if w.K == 1 {
		return 0
	} else {
		return math.NaN()
	}
}

// NumParameters returns the number of parameters in the distribution.
func (Weibull) NumParameters() int {
	return 2
}

// Prob computes the value of the probability density function at x.
func (w Weibull) Prob(x float64) float64 {
	if x < 0 {
		return 0
	} else {
		return math.Exp(w.LogProb(x))
	}
}

// Quantile returns the inverse of the cumulative probability distribution.
func (w Weibull) Quantile(p float64) float64 {
	if p < 0 || p > 1 {
		panic("weibull: percentile out of bounds")
	}
	return w.Lambda * math.Pow(-math.Log(1-p), 1/w.K)
}

// Rand returns a random sample drawn from the distribution.
func (w Weibull) Rand() float64 {
	var rnd float64
	if w.Source == nil {
		rnd = rand.Float64()
	} else {
		rnd = w.Source.Float64()
	}
	return w.Quantile(rnd)
}

// Skewness returns the skewness of the distribution.
func (w Weibull) Skewness() float64 {
	stdDev := w.StdDev()
	firstGamma, firstGammaSign := math.Lgamma(1 + 3/w.K)
	logFirst := firstGamma + 3*(math.Log(w.Lambda)-math.Log(stdDev))
	logSecond := math.Log(3) + math.Log(w.Mean()) + 2*math.Log(stdDev) - 3*math.Log(stdDev)
	logThird := 3 * (math.Log(w.Mean()) - math.Log(stdDev))
	return float64(firstGammaSign)*math.Exp(logFirst) - math.Exp(logSecond) - math.Exp(logThird)
}

// StdDev returns the standard deviation of the probability distribution.
func (w Weibull) StdDev() float64 {
	return math.Sqrt(w.Variance())
}

// Survival returns the survival function (complementary CDF) at x.
func (w Weibull) Survival(x float64) float64 {
	return math.Exp(w.LogSurvival(x))
}

// Variance returns the variance of the probability distribution.
func (w Weibull) Variance() float64 {
	return math.Pow(w.Lambda, 2) * (math.Gamma(1+2/w.K) - w.gammaIPow(1, 2))
}
