/*
Cluster Self-Managed Service API

Cluster Self-Managed Service API

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"time"
)

// ClusterStatus struct for ClusterStatus
type ClusterStatus struct {
	Type             *string    `json:"type,omitempty"`
	State            *string    `json:"state,omitempty"`
	Message          *string    `json:"message,omitempty"`
	UpdatedTimestamp *time.Time `json:"updated_timestamp,omitempty"`
}

// NewClusterStatus instantiates a new ClusterStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterStatus() *ClusterStatus {
	this := ClusterStatus{}
	return &this
}

// NewClusterStatusWithDefaults instantiates a new ClusterStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterStatusWithDefaults() *ClusterStatus {
	this := ClusterStatus{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *ClusterStatus) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterStatus) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ClusterStatus) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ClusterStatus) SetType(v string) {
	o.Type = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *ClusterStatus) GetState() string {
	if o == nil || o.State == nil {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterStatus) GetStateOk() (*string, bool) {
	if o == nil || o.State == nil {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *ClusterStatus) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *ClusterStatus) SetState(v string) {
	o.State = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ClusterStatus) GetMessage() string {
	if o == nil || o.Message == nil {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterStatus) GetMessageOk() (*string, bool) {
	if o == nil || o.Message == nil {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ClusterStatus) HasMessage() bool {
	if o != nil && o.Message != nil {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *ClusterStatus) SetMessage(v string) {
	o.Message = &v
}

// GetUpdatedTimestamp returns the UpdatedTimestamp field value if set, zero value otherwise.
func (o *ClusterStatus) GetUpdatedTimestamp() time.Time {
	if o == nil || o.UpdatedTimestamp == nil {
		var ret time.Time
		return ret
	}
	return *o.UpdatedTimestamp
}

// GetUpdatedTimestampOk returns a tuple with the UpdatedTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterStatus) GetUpdatedTimestampOk() (*time.Time, bool) {
	if o == nil || o.UpdatedTimestamp == nil {
		return nil, false
	}
	return o.UpdatedTimestamp, true
}

// HasUpdatedTimestamp returns a boolean if a field has been set.
func (o *ClusterStatus) HasUpdatedTimestamp() bool {
	if o != nil && o.UpdatedTimestamp != nil {
		return true
	}

	return false
}

// SetUpdatedTimestamp gets a reference to the given time.Time and assigns it to the UpdatedTimestamp field.
func (o *ClusterStatus) SetUpdatedTimestamp(v time.Time) {
	o.UpdatedTimestamp = &v
}

func (o ClusterStatus) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	if o.State != nil {
		toSerialize["state"] = o.State
	}
	if o.Message != nil {
		toSerialize["message"] = o.Message
	}
	if o.UpdatedTimestamp != nil {
		toSerialize["updated_timestamp"] = o.UpdatedTimestamp
	}
	return json.Marshal(toSerialize)
}

type NullableClusterStatus struct {
	value *ClusterStatus
	isSet bool
}

func (v NullableClusterStatus) Get() *ClusterStatus {
	return v.value
}

func (v *NullableClusterStatus) Set(val *ClusterStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterStatus(val *ClusterStatus) *NullableClusterStatus {
	return &NullableClusterStatus{value: val, isSet: true}
}

func (v NullableClusterStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}