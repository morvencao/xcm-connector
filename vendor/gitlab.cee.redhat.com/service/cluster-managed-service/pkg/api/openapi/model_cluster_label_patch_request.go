/*
Cluster Self-Managed Service API

Cluster Self-Managed Service API

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// ClusterLabelPatchRequest struct for ClusterLabelPatchRequest
type ClusterLabelPatchRequest struct {
	Value *string `json:"value,omitempty"`
}

// NewClusterLabelPatchRequest instantiates a new ClusterLabelPatchRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterLabelPatchRequest() *ClusterLabelPatchRequest {
	this := ClusterLabelPatchRequest{}
	return &this
}

// NewClusterLabelPatchRequestWithDefaults instantiates a new ClusterLabelPatchRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterLabelPatchRequestWithDefaults() *ClusterLabelPatchRequest {
	this := ClusterLabelPatchRequest{}
	return &this
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *ClusterLabelPatchRequest) GetValue() string {
	if o == nil || o.Value == nil {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLabelPatchRequest) GetValueOk() (*string, bool) {
	if o == nil || o.Value == nil {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *ClusterLabelPatchRequest) HasValue() bool {
	if o != nil && o.Value != nil {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *ClusterLabelPatchRequest) SetValue(v string) {
	o.Value = &v
}

func (o ClusterLabelPatchRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Value != nil {
		toSerialize["value"] = o.Value
	}
	return json.Marshal(toSerialize)
}

type NullableClusterLabelPatchRequest struct {
	value *ClusterLabelPatchRequest
	isSet bool
}

func (v NullableClusterLabelPatchRequest) Get() *ClusterLabelPatchRequest {
	return v.value
}

func (v *NullableClusterLabelPatchRequest) Set(val *ClusterLabelPatchRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterLabelPatchRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterLabelPatchRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterLabelPatchRequest(val *ClusterLabelPatchRequest) *NullableClusterLabelPatchRequest {
	return &NullableClusterLabelPatchRequest{value: val, isSet: true}
}

func (v NullableClusterLabelPatchRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterLabelPatchRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
