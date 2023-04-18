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

// ClusterLabelAllOf struct for ClusterLabelAllOf
type ClusterLabelAllOf struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

// NewClusterLabelAllOf instantiates a new ClusterLabelAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterLabelAllOf() *ClusterLabelAllOf {
	this := ClusterLabelAllOf{}
	return &this
}

// NewClusterLabelAllOfWithDefaults instantiates a new ClusterLabelAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterLabelAllOfWithDefaults() *ClusterLabelAllOf {
	this := ClusterLabelAllOf{}
	return &this
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *ClusterLabelAllOf) GetKey() string {
	if o == nil || o.Key == nil {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLabelAllOf) GetKeyOk() (*string, bool) {
	if o == nil || o.Key == nil {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *ClusterLabelAllOf) HasKey() bool {
	if o != nil && o.Key != nil {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *ClusterLabelAllOf) SetKey(v string) {
	o.Key = &v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *ClusterLabelAllOf) GetValue() string {
	if o == nil || o.Value == nil {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLabelAllOf) GetValueOk() (*string, bool) {
	if o == nil || o.Value == nil {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *ClusterLabelAllOf) HasValue() bool {
	if o != nil && o.Value != nil {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *ClusterLabelAllOf) SetValue(v string) {
	o.Value = &v
}

func (o ClusterLabelAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Key != nil {
		toSerialize["key"] = o.Key
	}
	if o.Value != nil {
		toSerialize["value"] = o.Value
	}
	return json.Marshal(toSerialize)
}

type NullableClusterLabelAllOf struct {
	value *ClusterLabelAllOf
	isSet bool
}

func (v NullableClusterLabelAllOf) Get() *ClusterLabelAllOf {
	return v.value
}

func (v *NullableClusterLabelAllOf) Set(val *ClusterLabelAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterLabelAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterLabelAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterLabelAllOf(val *ClusterLabelAllOf) *NullableClusterLabelAllOf {
	return &NullableClusterLabelAllOf{value: val, isSet: true}
}

func (v NullableClusterLabelAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterLabelAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
