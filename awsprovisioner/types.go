// This source code file is AUTO-GENERATED by github.com/taskcluster/jsonschema2go

package awsprovisioner

import (
	"encoding/json"
	"errors"

	"github.com/taskcluster/taskcluster-client-go/tcclient"
)

type (
	// A worker launchSpecification and required metadata
	CreateWorkerTypeRequest struct {

		// True if this worker type is allowed on demand instances.  Currently
		// ignored
		CanUseOndemand bool `json:"canUseOndemand"`

		// True if this worker type is allowed spot instances.  Currently ignored
		// as all instances are Spot
		CanUseSpot bool `json:"canUseSpot"`

		InstanceTypes []struct {

			// This number represents the number of tasks that this instance type
			// is capable of running concurrently.  This is used by the provisioner
			// to know how many pending tasks to offset a pending instance of this
			// type by
			Capacity float64 `json:"capacity"`

			// InstanceType name for Amazon.
			InstanceType string `json:"instanceType"`

			// LaunchSpecification entries unique to this InstanceType
			LaunchSpec json.RawMessage `json:"launchSpec"`

			// Scopes which should be included for this InstanceType.  Scopes must
			// be composed of printable ASCII characters and spaces.
			Scopes []string `json:"scopes"`

			// Static Secrets unique to this InstanceType
			Secrets json.RawMessage `json:"secrets"`

			// UserData entries unique to this InstanceType
			UserData json.RawMessage `json:"userData"`

			// This number is a relative measure of performance between two instance
			// types.  It is multiplied by the spot price from Amazon to figure out
			// which instance type is the cheapest one
			Utility float64 `json:"utility"`
		} `json:"instanceTypes"`

		// Launch Specification entries which are used in all regions and all instance types
		LaunchSpec json.RawMessage `json:"launchSpec"`

		// Maximum number of capacity units to be provisioned.
		MaxCapacity float64 `json:"maxCapacity"`

		// Maximum price we'll pay.  Like minPrice, this takes into account the
		// utility factor when figuring out what the actual SpotPrice submitted
		// to Amazon will be
		MaxPrice float64 `json:"maxPrice"`

		// Minimum number of capacity units to be provisioned.  A capacity unit
		// is an abstract unit of capacity, where one capacity unit is roughly
		// one task which should be taken off the queue
		MinCapacity float64 `json:"minCapacity"`

		// Minimum price to pay for an instance.  A Price is considered to be the
		// Amazon Spot Price multiplied by the utility factor of the InstantType
		// as specified in the instanceTypes list.  For example, if the minPrice
		// is set to $0.5 and the utility factor is 2, the actual minimum bid
		// used will be $0.25
		MinPrice float64 `json:"minPrice"`

		Regions []struct {

			// LaunchSpecification entries unique to this Region
			LaunchSpec struct {

				// Per-region AMI ImageId
				ImageID string `json:"ImageId"`
			} `json:"launchSpec"`

			// The Amazon AWS Region being configured.  Example: us-west-1
			Region string `json:"region"`

			// Scopes which should be included for this Region.  Scopes must be
			// composed of printable ASCII characters and spaces.
			Scopes []string `json:"scopes"`

			// Static Secrets unique to this Region
			Secrets json.RawMessage `json:"secrets"`

			// UserData entries unique to this Region
			UserData json.RawMessage `json:"userData"`
		} `json:"regions"`

		// A scaling ratio of `0.2` means that the provisioner will attempt to keep
		// the number of pending tasks around 20% of the provisioned capacity.
		// This results in pending tasks waiting 20% of the average task execution
		// time before starting to run.
		// A higher scaling ratio often results in better utilization and longer
		// waiting times. For workerTypes running long tasks a short scaling ratio
		// may be prefered, but for workerTypes running quick tasks a higher scaling
		// ratio may increase utilization without major delays.
		// If using a scaling ratio of 0, the provisioner will attempt to keep the
		// capacity of pending spot requests equal to the number of pending tasks.
		ScalingRatio float64 `json:"scalingRatio"`

		// Scopes to issue credentials to for all regions Scopes must be composed of
		// printable ASCII characters and spaces.
		Scopes []string `json:"scopes"`

		// Static secrets entries which are used in all regions and all instance types
		Secrets json.RawMessage `json:"secrets"`

		// UserData entries which are used in all regions and all instance types
		UserData json.RawMessage `json:"userData"`
	}

	// All of the launch specifications for a worker type
	GetAllLaunchSpecsResponse json.RawMessage

	// A Secret
	GetSecretRequest struct {

		// The date at which the secret is no longer guarunteed to exist
		Expiration tcclient.Time `json:"expiration"`

		// List of strings which are scopes for temporary credentials to give
		// to the worker through the secret system.  Scopes must be composed of
		// printable ASCII characters and spaces.
		Scopes []string `json:"scopes"`

		// Free form object which contains the secrets stored
		Secrets json.RawMessage `json:"secrets"`

		// A Slug ID which is the uniquely addressable token to access this
		// set of secrets
		Token string `json:"token"`

		// A string describing what the secret will be used for
		WorkerType string `json:"workerType"`
	}

	// Secrets from the provisioner
	GetSecretResponse struct {

		// Generated Temporary credentials from the Provisioner
		Credentials struct {
			AccessToken string `json:"accessToken"`

			Certificate string `json:"certificate"`

			ClientID string `json:"clientId"`
		} `json:"credentials"`

		// Free-form object which contains secrets from the worker type definition
		Data json.RawMessage `json:"data"`
	}

	// A worker launchSpecification and required metadata
	GetWorkerTypeRequest struct {

		// True if this worker type is allowed on demand instances.  Currently
		// ignored
		CanUseOndemand bool `json:"canUseOndemand"`

		// True if this worker type is allowed spot instances.  Currently ignored
		// as all instances are Spot
		CanUseSpot bool `json:"canUseSpot"`

		InstanceTypes []struct {

			// This number represents the number of tasks that this instance type
			// is capable of running concurrently.  This is used by the provisioner
			// to know how many pending tasks to offset a pending instance of this
			// type by
			Capacity float64 `json:"capacity"`

			// InstanceType name for Amazon.
			InstanceType string `json:"instanceType"`

			// LaunchSpecification entries unique to this InstanceType
			LaunchSpec json.RawMessage `json:"launchSpec"`

			// Scopes which should be included for this InstanceType.  Scopes must
			// be composed of printable ASCII characters and spaces.
			Scopes []string `json:"scopes"`

			// Static Secrets unique to this InstanceType
			Secrets json.RawMessage `json:"secrets"`

			// UserData entries unique to this InstanceType
			UserData json.RawMessage `json:"userData"`

			// This number is a relative measure of performance between two instance
			// types.  It is multiplied by the spot price from Amazon to figure out
			// which instance type is the cheapest one
			Utility float64 `json:"utility"`
		} `json:"instanceTypes"`

		// ISO Date string (e.g. new Date().toISOString()) which represents the time
		// when this worker type definition was last altered (inclusive of creation)
		LastModified tcclient.Time `json:"lastModified"`

		// Launch Specification entries which are used in all regions and all instance types
		LaunchSpec json.RawMessage `json:"launchSpec"`

		// Maximum number of capacity units to be provisioned.
		MaxCapacity float64 `json:"maxCapacity"`

		// Maximum price we'll pay.  Like minPrice, this takes into account the
		// utility factor when figuring out what the actual SpotPrice submitted
		// to Amazon will be
		MaxPrice float64 `json:"maxPrice"`

		// Minimum number of capacity units to be provisioned.  A capacity unit
		// is an abstract unit of capacity, where one capacity unit is roughly
		// one task which should be taken off the queue
		MinCapacity float64 `json:"minCapacity"`

		// Minimum price to pay for an instance.  A Price is considered to be the
		// Amazon Spot Price multiplied by the utility factor of the InstantType
		// as specified in the instanceTypes list.  For example, if the minPrice
		// is set to $0.5 and the utility factor is 2, the actual minimum bid
		// used will be $0.25
		MinPrice float64 `json:"minPrice"`

		Regions []struct {

			// LaunchSpecification entries unique to this Region
			LaunchSpec struct {

				// Per-region AMI ImageId
				ImageID string `json:"ImageId"`
			} `json:"launchSpec"`

			// The Amazon AWS Region being configured.  Example: us-west-1
			Region string `json:"region"`

			// Scopes which should be included for this Region.  Scopes must be
			// composed of printable ASCII characters and spaces.
			Scopes []string `json:"scopes"`

			// Static Secrets unique to this Region
			Secrets json.RawMessage `json:"secrets"`

			// UserData entries unique to this Region
			UserData json.RawMessage `json:"userData"`
		} `json:"regions"`

		// A scaling ratio of `0.2` means that the provisioner will attempt to keep
		// the number of pending tasks around 20% of the provisioned capacity.
		// This results in pending tasks waiting 20% of the average task execution
		// time before starting to run.
		// A higher scaling ratio often results in better utilization and longer
		// waiting times. For workerTypes running long tasks a short scaling ratio
		// may be prefered, but for workerTypes running quick tasks a higher scaling
		// ratio may increase utilization without major delays.
		// If using a scaling ratio of 0, the provisioner will attempt to keep the
		// capacity of pending spot requests equal to the number of pending tasks.
		ScalingRatio float64 `json:"scalingRatio"`

		// Scopes to issue credentials to for all regions.  Scopes must be composed
		// of printable ASCII characters and spaces.
		Scopes []string `json:"scopes"`

		// Static secrets entries which are used in all regions and all instance types
		Secrets json.RawMessage `json:"secrets"`

		// UserData entries which are used in all regions and all instance types
		UserData json.RawMessage `json:"userData"`

		// The ID of the workerType
		//
		// Syntax:     ^[A-Za-z0-9+/=_-]{1,22}$
		WorkerType string `json:"workerType"`
	}

	ListWorkerTypes []string
)

// MarshalJSON calls json.RawMessage method of the same name. Required since
// GetAllLaunchSpecsResponse is of type json.RawMessage...
func (this *GetAllLaunchSpecsResponse) MarshalJSON() ([]byte, error) {
	x := json.RawMessage(*this)
	return (&x).MarshalJSON()
}

// UnmarshalJSON is a copy of the json.RawMessage implementation.
func (this *GetAllLaunchSpecsResponse) UnmarshalJSON(data []byte) error {
	if this == nil {
		return errors.New("GetAllLaunchSpecsResponse: UnmarshalJSON on nil pointer")
	}
	*this = append((*this)[0:0], data...)
	return nil
}
