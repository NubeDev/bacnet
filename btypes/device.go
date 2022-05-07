package btypes

type Enumerated uint32

type IAm struct {
	ID           ObjectID
	MaxApdu      uint32
	Segmentation Enumerated
	Vendor       uint32
	Addr         Address
}

type Device struct {
	ID           ObjectID
	MaxApdu      uint32
	Segmentation Enumerated
	Vendor       uint32
	Addr         Address
	Objects      ObjectMap
	IsTypeMSTP   bool
}

// ObjectSlice returns all the objects in the device as a slice (not thread-safe)
func (dev *Device) ObjectSlice() []Object {
	var objs []Object
	for _, objMap := range dev.Objects {
		for _, o := range objMap {
			objs = append(objs, o)
		}
	}
	return objs
}
