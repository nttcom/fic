/*
Copyright 2020 NTT Limited and NTT Communications Corporation All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"github.com/nttcom/go-fic"
	"github.com/nttcom/go-fic/fic/utils"
	"github.com/spf13/viper"
)

func NewClient() (*fic.ServiceClient, error) {
	opts := fic.AuthOptions{
		IdentityEndpoint: viper.GetString("auth_url"),
		Username:         viper.GetString("username"),
		Password:         viper.GetString("password"),
		TenantID:         viper.GetString("tenant_id"),
		DomainID:         "default",
	}

	provider, err := utils.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	client, err := utils.NewEriV1(provider, fic.EndpointOpts{Region: "gl1"})
	if err != nil {
		return nil, err
	}

	return client, nil
}
