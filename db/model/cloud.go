// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons 
// to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or 
// substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package model

//TableName
func (t *RegionUserInfo) TableName() string {
	return "user_region_info"
}

//RegionUserInfo
type RegionUserInfo struct {
	Model
	EID            string `gorm:"column:eid;size:34" json:"eid"`
	APIRange       string `gorm:"column:api_range;size:24" json:"api_range"`
	RegionTag      string `gorm:"column:region_tag;size:24" json:"region_tag"`
	ValidityPeriod int    `gorm:"column:validity_period;size:10" json:"validity_period"`
	Token          string `gorm:"column:token;size:32" json:"token"`
	CA             string `gorm:"column:ca;size:4096" json:"ca"`
	Key            string `gorm:"column:key;size:4096" json:"key"`
}
