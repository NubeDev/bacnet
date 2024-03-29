// Code generated by "stringer -type=ErrorClass"; DO NOT EDIT.

package bacerr

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DeviceError-0]
	_ = x[ObjectError-1]
	_ = x[PropertyError-2]
	_ = x[ResourcesError-3]
	_ = x[SecurityError-4]
	_ = x[ServicesError-5]
	_ = x[VTError-6]
	_ = x[CommunicationError-7]
}

const _ErrorClass_name = "DeviceErrorObjectErrorPropertyErrorResourcesErrorSecurityErrorServicesErrorVTErrorCommunicationError"

var _ErrorClass_index = [...]uint8{0, 11, 22, 35, 49, 62, 75, 82, 100}

func (i ErrorClass) String() string {
	if i >= ErrorClass(len(_ErrorClass_index)-1) {
		return "ErrorClass(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrorClass_name[_ErrorClass_index[i]:_ErrorClass_index[i+1]]
}
