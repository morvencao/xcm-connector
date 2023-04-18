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

// ClusterLabelList struct for ClusterLabelList
type ClusterLabelList struct {
	Kind  string         `json:"kind"`
	Href  *string        `json:"href,omitempty"`
	Page  *int32         `json:"page,omitempty"`
	Size  *int32         `json:"size,omitempty"`
	Total *int32         `json:"total,omitempty"`
	Items []ClusterLabel `json:"items"`
}

// NewClusterLabelList instantiates a new ClusterLabelList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterLabelList(kind string, items []ClusterLabel) *ClusterLabelList {
	this := ClusterLabelList{}
	this.Kind = kind
	this.Items = items
	return &this
}

// NewClusterLabelListWithDefaults instantiates a new ClusterLabelList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterLabelListWithDefaults() *ClusterLabelList {
	this := ClusterLabelList{}
	return &this
}

// GetKind returns the Kind field value
func (o *ClusterLabelList) GetKind() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Kind
}

// GetKindOk returns a tuple with the Kind field value
// and a boolean to check if the value has been set.
func (o *ClusterLabelList) GetKindOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Kind, true
}

// SetKind sets field value
func (o *ClusterLabelList) SetKind(v string) {
	o.Kind = v
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *ClusterLabelList) GetHref() string {
	if o == nil || o.Href == nil {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLabelList) GetHrefOk() (*string, bool) {
	if o == nil || o.Href == nil {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *ClusterLabelList) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *ClusterLabelList) SetHref(v string) {
	o.Href = &v
}

// GetPage returns the Page field value if set, zero value otherwise.
func (o *ClusterLabelList) GetPage() int32 {
	if o == nil || o.Page == nil {
		var ret int32
		return ret
	}
	return *o.Page
}

// GetPageOk returns a tuple with the Page field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLabelList) GetPageOk() (*int32, bool) {
	if o == nil || o.Page == nil {
		return nil, false
	}
	return o.Page, true
}

// HasPage returns a boolean if a field has been set.
func (o *ClusterLabelList) HasPage() bool {
	if o != nil && o.Page != nil {
		return true
	}

	return false
}

// SetPage gets a reference to the given int32 and assigns it to the Page field.
func (o *ClusterLabelList) SetPage(v int32) {
	o.Page = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *ClusterLabelList) GetSize() int32 {
	if o == nil || o.Size == nil {
		var ret int32
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLabelList) GetSizeOk() (*int32, bool) {
	if o == nil || o.Size == nil {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *ClusterLabelList) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// SetSize gets a reference to the given int32 and assigns it to the Size field.
func (o *ClusterLabelList) SetSize(v int32) {
	o.Size = &v
}

// GetTotal returns the Total field value if set, zero value otherwise.
func (o *ClusterLabelList) GetTotal() int32 {
	if o == nil || o.Total == nil {
		var ret int32
		return ret
	}
	return *o.Total
}

// GetTotalOk returns a tuple with the Total field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLabelList) GetTotalOk() (*int32, bool) {
	if o == nil || o.Total == nil {
		return nil, false
	}
	return o.Total, true
}

// HasTotal returns a boolean if a field has been set.
func (o *ClusterLabelList) HasTotal() bool {
	if o != nil && o.Total != nil {
		return true
	}

	return false
}

// SetTotal gets a reference to the given int32 and assigns it to the Total field.
func (o *ClusterLabelList) SetTotal(v int32) {
	o.Total = &v
}

// GetItems returns the Items field value
func (o *ClusterLabelList) GetItems() []ClusterLabel {
	if o == nil {
		var ret []ClusterLabel
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *ClusterLabelList) GetItemsOk() ([]ClusterLabel, bool) {
	if o == nil {
		return nil, false
	}
	return o.Items, true
}

// SetItems sets field value
func (o *ClusterLabelList) SetItems(v []ClusterLabel) {
	o.Items = v
}

func (o ClusterLabelList) MarshalJSON() ([]byte, error) {
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

type NullableClusterLabelList struct {
	value *ClusterLabelList
	isSet bool
}

func (v NullableClusterLabelList) Get() *ClusterLabelList {
	return v.value
}

func (v *NullableClusterLabelList) Set(val *ClusterLabelList) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterLabelList) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterLabelList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterLabelList(val *ClusterLabelList) *NullableClusterLabelList {
	return &NullableClusterLabelList{value: val, isSet: true}
}

func (v NullableClusterLabelList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterLabelList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
