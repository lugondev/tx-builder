// tokucore
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
// Copyright (c) 2018 Stanislav Fomichev.
// Copyright 2019 by KeyFuse Labs
// BSD License

package secp256k1

import (
	"crypto/elliptic"
	"math/big"
	"sync"
)

var (
	initonce  sync.Once
	secp256k1 *secp256k1Curve
	three     = new(big.Int).SetUint64(3)
)

type secp256k1Curve struct {
	elliptic.CurveParams
}

func initSECP256K1() {
	// http://www.secg.org/sec2-v2.pdf
	secp256k1 = &secp256k1Curve{elliptic.CurveParams{Name: "secp256k1"}}
	secp256k1.P, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f", 16)
	secp256k1.N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	secp256k1.Gx, _ = new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	secp256k1.Gy, _ = new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	secp256k1.BitSize = 256
}

// SECP256K1 -- returns secp256k1 curve
func SECP256K1() elliptic.Curve {
	initonce.Do(initSECP256K1)
	return secp256k1
}

// Params -- returns elliptic curve params
func (curve *secp256k1Curve) Params() *elliptic.CurveParams {
	return &curve.CurveParams
}

// IsOnCurve -- check the point (x,y) is on the curve or not.
func (curve *secp256k1Curve) IsOnCurve(x, y *big.Int) bool {
	// y² = x³ + b
	var y2, x3 big.Int

	y2.Mul(y, y)         //y²
	y2.Mod(&y2, curve.P) //y²%P

	x3.Mul(x, x)         //x²
	x3.Mul(&x3, x)       //x³
	x3.Add(&x3, curve.B) //x³+B
	x3.Mod(&x3, curve.P) //(x³+B)%P
	return x3.Cmp(&y2) == 0
}

// zForAffine --
// returns a Jacobian Z value for the affine point (x, y). If x and
// y are zero, it assumes that they represent the point at infinity because (0,
// 0) is not on the any of the curves handled here.
func zForAffine(x, y *big.Int) *big.Int {
	z := new(big.Int)
	if x.Sign() != 0 || y.Sign() != 0 {
		z.SetInt64(1)
	}
	return z
}

// affineFromJacobian --
// reverses the Jacobian transform. If the point is ∞ it returns 0, 0.
// For a given (x, y) position on the curve, the Jacobian coordinates are (x1, y1, z1)
// where x = x1/z1² and y = y1/z1³. The greatest speedups come when the whole
// calculation can be performed within the transform (as in ScalarMult and
// ScalarBaseMult). But even for Add and Double, it's faster to apply and
// reverse the transform than to operate in affine coordinates.
func (curve *secp256k1Curve) affineFromJacobian(x, y, z *big.Int) (xOut, yOut *big.Int) {
	if z.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}

	zinv := new(big.Int).ModInverse(z, curve.P)
	zinvsq := new(big.Int).Mul(zinv, zinv)

	xOut = new(big.Int).Mul(x, zinvsq)
	xOut.Mod(xOut, curve.P)
	zinvsq.Mul(zinvsq, zinv)
	yOut = new(big.Int).Mul(y, zinvsq)
	yOut.Mod(yOut, curve.P)
	return
}

// Add -- returns the sum of (x1,y1) and (x2,y2).
func (curve *secp256k1Curve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	z1 := zForAffine(x1, y1)
	z2 := zForAffine(x2, y2)
	return curve.affineFromJacobian(curve.addJacobian(x1, y1, z1, x2, y2, z2))
}

// addJacobian --
// takes two points in Jacobian coordinates, (x1, y1, z1) and
// (x2, y2, z2) and returns their sum, also in Jacobian form.
func (curve *secp256k1Curve) addJacobian(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	x3, y3, z3 := new(big.Int), new(big.Int), new(big.Int)
	if z1.Sign() == 0 {
		x3.Set(x2)
		y3.Set(y2)
		z3.Set(z2)
		return x3, y3, z3
	}
	if z2.Sign() == 0 {
		x3.Set(x1)
		y3.Set(y1)
		z3.Set(z1)
		return x3, y3, z3
	}

	z1z1 := new(big.Int).Mul(z1, z1) //z1²
	z1z1.Mod(z1z1, curve.P)          //z1²%P
	z2z2 := new(big.Int).Mul(z2, z2) //z2²
	z2z2.Mod(z2z2, curve.P)          //z2²%P

	u1 := new(big.Int).Mul(x1, z2z2) //x1*z2²
	u1.Mod(u1, curve.P)              //(x1*z2²)%P
	u2 := new(big.Int).Mul(x2, z1z1) //x2*z1²
	u2.Mod(u2, curve.P)              //(x2*z1²)%P
	h := new(big.Int).Sub(u2, u1)    //u2-u1
	xEqual := h.Sign() == 0
	if h.Sign() == -1 {
		h.Add(h, curve.P)
	}
	i := new(big.Int).Lsh(h, 1)
	i.Mul(i, i)
	j := new(big.Int).Mul(h, i)

	s1 := new(big.Int).Mul(y1, z2)
	s1.Mul(s1, z2z2)
	s1.Mod(s1, curve.P)
	s2 := new(big.Int).Mul(y2, z1)
	s2.Mul(s2, z1z1)
	s2.Mod(s2, curve.P)
	r := new(big.Int).Sub(s2, s1)
	if r.Sign() == -1 {
		r.Add(r, curve.P)
	}
	yEqual := r.Sign() == 0
	if xEqual && yEqual {
		return curve.doubleJacobian(x1, y1, z1)
	}
	r.Lsh(r, 1) //r²
	v := new(big.Int).Mul(u1, i)

	x3.Set(r)
	x3.Mul(x3, x3)
	x3.Sub(x3, j)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	x3.Mod(x3, curve.P)

	y3.Set(r)
	v.Sub(v, x3)
	y3.Mul(y3, v)
	s1.Mul(s1, j)
	s1.Lsh(s1, 1)
	y3.Sub(y3, s1)
	y3.Mod(y3, curve.P)

	z3.Add(z1, z2)
	z3.Mul(z3, z3)
	z3.Sub(z3, z1z1)
	z3.Sub(z3, z2z2)
	z3.Mul(z3, h)
	z3.Mod(z3, curve.P)
	return x3, y3, z3
}

// Double -- returns 2*(x,y).
func (curve *secp256k1Curve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	z1 := zForAffine(x1, y1)
	return curve.affineFromJacobian(curve.doubleJacobian(x1, y1, z1))
}

// doubleJacobian -- takes a point in Jacobian coordinates, (x, y, z),
// and returns its double, also in Jacobian form.
func (curve *secp256k1Curve) doubleJacobian(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	var a, b, c, d, e, f, x3, y3, z3 big.Int

	a.Mul(x, x)
	a.Mod(&a, curve.P)
	b.Mul(y, y)
	b.Mod(&b, curve.P)
	c.Mul(&b, &b)
	c.Mod(&c, curve.P)

	d.Add(x, &b)
	d.Mul(&d, &d)
	d.Sub(&d, &a)
	d.Sub(&d, &c)
	d.Lsh(&d, 1)
	if d.Sign() < 0 {
		d.Add(&d, curve.P)
	} else {
		d.Mod(&d, curve.P)
	}

	e.Mul(three, &a)
	e.Mod(&e, curve.P)
	f.Mul(&e, &e)
	f.Mod(&f, curve.P)

	x3.Lsh(&d, 1)
	x3.Sub(&f, &x3)
	if x3.Sign() < 0 {
		x3.Add(&x3, curve.P)
	} else {
		x3.Mod(&x3, curve.P)
	}

	y3.Sub(&d, &x3)
	y3.Mul(&e, &y3)
	c.Lsh(&c, 3)
	y3.Sub(&y3, &c)
	if y3.Sign() < 0 {
		y3.Add(&y3, curve.P)
	} else {
		y3.Mod(&y3, curve.P)
	}

	z3.Mul(y, z)
	z3.Lsh(&z3, 1)
	z3.Mod(&z3, curve.P)
	return &x3, &y3, &z3
}

// ScalarMult -- returns k*(Bx,By) where k is a number in big-endian form.
func (curve *secp256k1Curve) ScalarMult(Bx, By *big.Int, k []byte) (*big.Int, *big.Int) {
	Bz := new(big.Int).SetInt64(1)
	x, y, z := new(big.Int), new(big.Int), new(big.Int)

	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			x, y, z = curve.doubleJacobian(x, y, z)
			if byte&0x80 == 0x80 {
				x, y, z = curve.addJacobian(Bx, By, Bz, x, y, z)
			}
			byte <<= 1
		}
	}
	return curve.affineFromJacobian(x, y, z)
}

// ScalarBaseMult -- returns k*G.
func (curve *secp256k1Curve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	return curve.ScalarMult(curve.Gx, curve.Gy, k)
}
