package linear_graphql

import (
	"encoding/json"
	"reflect"
	"strings"
)

// marshalOmitNil builds a JSON object from a struct, skipping fields whose
// values are nil pointers or nil slices. This is needed because genqlient
// generates structs without omitempty tags, causing Linear's API to reject
// null values for optional input fields.
func marshalOmitNil(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	m := make(map[string]any, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		// Strip options like ",omitempty" from the tag.
		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}
		fv := rv.Field(i)
		switch fv.Kind() {
		case reflect.Ptr, reflect.Interface:
			if fv.IsNil() {
				continue
			}
		case reflect.Slice, reflect.Map:
			if fv.IsNil() {
				continue
			}
		}
		m[tag] = fv.Interface()
	}
	return json.Marshal(m)
}

// MarshalJSON overrides for all input and filter types used by the CLI.
// These are needed because genqlient v0.8.1 does not add omitempty to
// schema input types, causing Linear's API to reject requests with null
// values for optional fields.

func (v IssueUpdateInput) MarshalJSON() ([]byte, error) { return marshalOmitNil(v) }
func (v IssueCreateInput) MarshalJSON() ([]byte, error) { return marshalOmitNil(v) }
func (v IssueFilter) MarshalJSON() ([]byte, error)      { return marshalOmitNil(v) }
func (v IssueCollectionFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}
func (v TeamFilter) MarshalJSON() ([]byte, error)        { return marshalOmitNil(v) }
func (v IDComparator) MarshalJSON() ([]byte, error)      { return marshalOmitNil(v) }
func (v DateComparator) MarshalJSON() ([]byte, error)    { return marshalOmitNil(v) }
func (v StringComparator) MarshalJSON() ([]byte, error)  { return marshalOmitNil(v) }
func (v NumberComparator) MarshalJSON() ([]byte, error)  { return marshalOmitNil(v) }
func (v BooleanComparator) MarshalJSON() ([]byte, error) { return marshalOmitNil(v) }
func (v EstimateComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableDateComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableStringComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableNumberComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableUserFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableCycleFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableProjectFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableIssueFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v NullableTeamFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}
func (v UserFilter) MarshalJSON() ([]byte, error)    { return marshalOmitNil(v) }
func (v CycleFilter) MarshalJSON() ([]byte, error)   { return marshalOmitNil(v) }
func (v ProjectFilter) MarshalJSON() ([]byte, error) { return marshalOmitNil(v) }
func (v WorkflowStateFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v IssueLabelFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v IssueLabelCollectionFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}
func (v CommentFilter) MarshalJSON() ([]byte, error) { return marshalOmitNil(v) }
func (v CommentCollectionFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v UserCollectionFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v TeamCollectionFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v ProjectCollectionFilter) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v IssueIDComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v ContentComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v RelationExistsComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v SlaStatusComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v SourceTypeComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}

func (v SubTypeComparator) MarshalJSON() ([]byte, error) {
	return marshalOmitNil(v)
}
