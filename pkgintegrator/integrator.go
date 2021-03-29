// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package pkgintegrator

import "github.com/gin-gonic/gin"

type PKGIntegrator struct {
	integrations []gin.HandlerFunc
}

var integrator *PKGIntegrator

func New() *PKGIntegrator {
	integrator = &PKGIntegrator{}

	return integrator
}

func Resolve() *PKGIntegrator {
	return integrator
}

func (i *PKGIntegrator) Integrate(pkgIntegration gin.HandlerFunc) {
	i.integrations = append(i.integrations, pkgIntegration)
}

func (i *PKGIntegrator) GetIntegrations() []gin.HandlerFunc {
	return i.integrations
}
