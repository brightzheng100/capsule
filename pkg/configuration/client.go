// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package configuration

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
	machineryerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	capsulev1alpha1 "github.com/clastix/capsule/api/v1alpha1"
)

// capsuleConfiguration is the Capsule Configuration retrieval mode
// using a closure that provides the desired configuration.
type capsuleConfiguration struct {
	retrievalFn func() *capsulev1alpha1.CapsuleConfiguration
}

func NewCapsuleConfiguration(client client.Client, name string) Configuration {
	return &capsuleConfiguration{retrievalFn: func() *capsulev1alpha1.CapsuleConfiguration {
		config := &capsulev1alpha1.CapsuleConfiguration{}

		if err := client.Get(context.Background(), types.NamespacedName{Name: name}, config); err != nil {
			if machineryerr.IsNotFound(err) {
				return &capsulev1alpha1.CapsuleConfiguration{
					Spec: capsulev1alpha1.CapsuleConfigurationSpec{
						UserGroups:                     []string{"capsule.clastix.io"},
						ForceTenantPrefix:              false,
						ProtectedNamespaceRegexpString: "",
					},
				}
			}
			panic(errors.Wrap(err, "Cannot retrieve Capsule configuration with name "+name))
		}

		return config
	}}
}

func (c capsuleConfiguration) ProtectedNamespaceRegexp() (*regexp.Regexp, error) {
	expr := c.retrievalFn().Spec.ProtectedNamespaceRegexpString
	if len(expr) == 0 {
		return nil, nil
	}
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot compile the protected namespace regexp")
	}

	return r, nil
}

func (c capsuleConfiguration) ForceTenantPrefix() bool {
	return c.retrievalFn().Spec.ForceTenantPrefix
}

func (c capsuleConfiguration) UserGroups() []string {
	return c.retrievalFn().Spec.UserGroups
}
