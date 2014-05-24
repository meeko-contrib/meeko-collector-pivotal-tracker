// Copyright (c) 2013-2014 The meeko-collector-pivotal-tracker AUTHORS
//
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package data

type StringValueChange struct {
	Original *StringValue
	Current  *StringValue
}

type StringValue struct {
	Value     string
	UpdatedAt float64
}
