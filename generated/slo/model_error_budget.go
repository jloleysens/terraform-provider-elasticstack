/*
SLOs

OpenAPI schema for SLOs endpoints

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package slo

import (
	"encoding/json"
)

// checks if the ErrorBudget type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ErrorBudget{}

// ErrorBudget struct for ErrorBudget
type ErrorBudget struct {
	// The initial error budget, as 1 - objective
	Initial float64 `json:"initial"`
	// The error budget consummed, as a percentage of the initial value.
	Consumed float64 `json:"consumed"`
	// The error budget remaining, as a percentage of the initial value.
	Remaining float64 `json:"remaining"`
	// Only for SLO defined with occurrences budgeting method and calendar aligned time window.
	IsEstimated bool `json:"isEstimated"`
}

// NewErrorBudget instantiates a new ErrorBudget object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewErrorBudget(initial float64, consumed float64, remaining float64, isEstimated bool) *ErrorBudget {
	this := ErrorBudget{}
	this.Initial = initial
	this.Consumed = consumed
	this.Remaining = remaining
	this.IsEstimated = isEstimated
	return &this
}

// NewErrorBudgetWithDefaults instantiates a new ErrorBudget object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewErrorBudgetWithDefaults() *ErrorBudget {
	this := ErrorBudget{}
	return &this
}

// GetInitial returns the Initial field value
func (o *ErrorBudget) GetInitial() float64 {
	if o == nil {
		var ret float64
		return ret
	}

	return o.Initial
}

// GetInitialOk returns a tuple with the Initial field value
// and a boolean to check if the value has been set.
func (o *ErrorBudget) GetInitialOk() (*float64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Initial, true
}

// SetInitial sets field value
func (o *ErrorBudget) SetInitial(v float64) {
	o.Initial = v
}

// GetConsumed returns the Consumed field value
func (o *ErrorBudget) GetConsumed() float64 {
	if o == nil {
		var ret float64
		return ret
	}

	return o.Consumed
}

// GetConsumedOk returns a tuple with the Consumed field value
// and a boolean to check if the value has been set.
func (o *ErrorBudget) GetConsumedOk() (*float64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Consumed, true
}

// SetConsumed sets field value
func (o *ErrorBudget) SetConsumed(v float64) {
	o.Consumed = v
}

// GetRemaining returns the Remaining field value
func (o *ErrorBudget) GetRemaining() float64 {
	if o == nil {
		var ret float64
		return ret
	}

	return o.Remaining
}

// GetRemainingOk returns a tuple with the Remaining field value
// and a boolean to check if the value has been set.
func (o *ErrorBudget) GetRemainingOk() (*float64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Remaining, true
}

// SetRemaining sets field value
func (o *ErrorBudget) SetRemaining(v float64) {
	o.Remaining = v
}

// GetIsEstimated returns the IsEstimated field value
func (o *ErrorBudget) GetIsEstimated() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.IsEstimated
}

// GetIsEstimatedOk returns a tuple with the IsEstimated field value
// and a boolean to check if the value has been set.
func (o *ErrorBudget) GetIsEstimatedOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IsEstimated, true
}

// SetIsEstimated sets field value
func (o *ErrorBudget) SetIsEstimated(v bool) {
	o.IsEstimated = v
}

func (o ErrorBudget) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ErrorBudget) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["initial"] = o.Initial
	toSerialize["consumed"] = o.Consumed
	toSerialize["remaining"] = o.Remaining
	toSerialize["isEstimated"] = o.IsEstimated
	return toSerialize, nil
}

type NullableErrorBudget struct {
	value *ErrorBudget
	isSet bool
}

func (v NullableErrorBudget) Get() *ErrorBudget {
	return v.value
}

func (v *NullableErrorBudget) Set(val *ErrorBudget) {
	v.value = val
	v.isSet = true
}

func (v NullableErrorBudget) IsSet() bool {
	return v.isSet
}

func (v *NullableErrorBudget) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableErrorBudget(val *ErrorBudget) *NullableErrorBudget {
	return &NullableErrorBudget{value: val, isSet: true}
}

func (v NullableErrorBudget) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableErrorBudget) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
