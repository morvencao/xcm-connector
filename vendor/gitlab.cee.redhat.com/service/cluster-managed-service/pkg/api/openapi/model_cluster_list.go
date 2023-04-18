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

// ClusterList struct for ClusterList
type ClusterList struct {
	Kind  string    `json:"kind"`
	Href  *string   `json:"href,omitempty"`
	Page  *int32    `json:"page,omitempty"`
	Size  *int32    `json:"size,omitempty"`
	Total *int32    `json:"total,omitempty"`
	Items []Cluster `json:"items"`
}

// NewClusterList instantiates a new ClusterList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterList(kind string, items []Cluster) *ClusterList {
	this := ClusterList{}
	this.Kind = kind
	this.Items = items
	return &this
}

// NewClusterListWithDefaults instantiates a new ClusterList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterListWithDefaults() *ClusterList {
	this := ClusterList{}
	return &this
}

// GetKind returns the Kind field value
func (o *ClusterList) GetKind() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Kind
}

// GetKindOk returns a tuple with the Kind field value
// and a boolean to check if the value has been set.
func (o *ClusterList) GetKindOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Kind, true
}

// SetKind sets field value
func (o *ClusterList) SetKind(v string) {
	o.Kind = v
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *ClusterList) GetHref() string {
	if o == nil || o.Href == nil {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetHrefOk() (*string, bool) {
	if o == nil || o.Href == nil {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *ClusterList) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *ClusterList) SetHref(v string) {
	o.Href = &v
}

// GetPage returns the Page field value if set, zero value otherwise.
func (o *ClusterList) GetPage() int32 {
	if o == nil || o.Page == nil {
		var ret int32
		return ret
	}
	return *o.Page
}

// GetPageOk returns a tuple with the Page field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetPageOk() (*int32, bool) {
	if o == nil || o.Page == nil {
		return nil, false
	}
	return o.Page, true
}

// HasPage returns a boolean if a field has been set.
func (o *ClusterList) HasPage() bool {
	if o != nil && o.Page != nil {
		return true
	}

	return false
}

// SetPage gets a reference to the given int32 and assigns it to the Page field.
func (o *ClusterList) SetPage(v int32) {
	o.Page = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *ClusterList) GetSize() int32 {
	if o == nil || o.Size == nil {
		var ret int32
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetSizeOk() (*int32, bool) {
	if o == nil || o.Size == nil {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *ClusterList) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// SetSize gets a reference to the given int32 and assigns it to the Size field.
func (o *ClusterList) SetSize(v int32) {
	o.Size = &v
}

// GetTotal returns the Total field value if set, zero value otherwise.
func (o *ClusterList) GetTotal() int32 {
	if o == nil || o.Total == nil {
		var ret int32
		return ret
	}
	return *o.Total
}

// GetTotalOk returns a tuple with the Total field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetTotalOk() (*int32, bool) {
	if o == nil || o.Total == nil {
		return nil, false
	}
	return o.Total, true
}

// HasTotal returns a boolean if a field has been set.
func (o *ClusterList) HasTotal() bool {
	if o != nil && o.Total != nil {
		return true
	}

	return false
}

// SetTotal gets a reference to the given int32 and assigns it to the Total field.
func (o *ClusterList) SetTotal(v int32) {
	o.Total = &v
}

// GetItems returns the Items field value
func (o *ClusterList) GetItems() []Cluster {
	if o == nil {
		var ret []Cluster
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *ClusterList) GetItemsOk() ([]Cluster, bool) {
	if o == nil {
		return nil, false
	}
	return o.Items, true
}

// SetItems sets field value
func (o *ClusterList) SetItems(v []Cluster) {
	o.Items = v
}

func (o ClusterList) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["kind"] = o.Kind
	}
	if o.Href != nil {
		toSerialize["href"] = o.Href
	}
	if o.Page != nil {
		toSerialize["page"] = o.Page
	}
	if o.Size != nil {
		toSerialize["size"] = o.Size
	}
	if o.Total != nil {
		toSerialize["total"] = o.Total
	}
	if true {
		toSerialize["items"] = o.Items
	}
	return json.Marshal(toSerialize)
}

type NullableClusterList struct {
	value *ClusterList
	isSet bool
}

func (v NullableClusterList) Get() *ClusterList {
	return v.value
}

func (v *NullableClusterList) Set(val *ClusterList) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterList) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterList(val *ClusterList) *NullableClusterList {
	return &NullableClusterList{value: val, isSet: true}
}

func (v NullableClusterList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}