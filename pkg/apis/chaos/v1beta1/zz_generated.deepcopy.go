// +build !ignore_autogenerated

/*
Copyright 2019 Datadog.

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
// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DependencyFailureInjection) DeepCopyInto(out *DependencyFailureInjection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DependencyFailureInjection.
func (in *DependencyFailureInjection) DeepCopy() *DependencyFailureInjection {
	if in == nil {
		return nil
	}
	out := new(DependencyFailureInjection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DependencyFailureInjection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DependencyFailureInjectionList) DeepCopyInto(out *DependencyFailureInjectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DependencyFailureInjection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DependencyFailureInjectionList.
func (in *DependencyFailureInjectionList) DeepCopy() *DependencyFailureInjectionList {
	if in == nil {
		return nil
	}
	out := new(DependencyFailureInjectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DependencyFailureInjectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DependencyFailureInjectionSpec) DeepCopyInto(out *DependencyFailureInjectionSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DependencyFailureInjectionSpec.
func (in *DependencyFailureInjectionSpec) DeepCopy() *DependencyFailureInjectionSpec {
	if in == nil {
		return nil
	}
	out := new(DependencyFailureInjectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DependencyFailureInjectionStatus) DeepCopyInto(out *DependencyFailureInjectionStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DependencyFailureInjectionStatus.
func (in *DependencyFailureInjectionStatus) DeepCopy() *DependencyFailureInjectionStatus {
	if in == nil {
		return nil
	}
	out := new(DependencyFailureInjectionStatus)
	in.DeepCopyInto(out)
	return out
}
