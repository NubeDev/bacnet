package local

import (
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/helpers/data"
)

type Point struct {
	ObjectID      btypes.ObjectInstance
	ObjectType    btypes.ObjectType
	WriteValue    interface{}
	WriteNull     bool
	WritePriority uint8
}

/*
***** READS *****
 */

//PointReadFloat64 use this when wanting to read point values for an AI, AV, AO
func (device *Device) PointReadFloat64(pnt *Point) (float64, error) {
	if device.isPointFloat(pnt) {

	}
	obj := &Object{
		ObjectID:   pnt.ObjectID,
		ObjectType: pnt.ObjectType,
		Prop:       btypes.PropPresentValue,
		ArrayIndex: bacnet.ArrayAll,
	}

	read, err := device.Read(obj)
	if err != nil {
		return 0, err
	}
	return device.toFloat(read), nil
}

//PointReadBool use this when wanting to read point values for an BI, BV, BO
func (device *Device) PointReadBool(pnt *Point) (bool, error) {
	if !device.isPointBool(pnt) {

	}
	obj := &Object{
		ObjectID:   pnt.ObjectID,
		ObjectType: pnt.ObjectType,
		Prop:       btypes.PropPresentValue,
		ArrayIndex: bacnet.ArrayAll,
	}

	read, err := device.Read(obj)
	if err != nil {
		return false, err
	}
	return device.toBool(read), nil
}

func (device *Device) PointReleaseOverride(pnt *Point) (bool, error) {
	if !device.isPointWriteable(pnt) {
		//TODO add errors
	}
	obj := &Object{
		ObjectID:   pnt.ObjectID,
		ObjectType: pnt.ObjectType,
		Prop:       btypes.PropPresentValue,
		ArrayIndex: bacnet.ArrayAll,
	}

	read, err := device.Read(obj)
	if err != nil {
		return false, err
	}
	return device.toBool(read), nil
}

/*
***** WRITES *****
 */

//PointWriteAnalogue use this when wanting to write a new value for an AV, AO
func (device *Device) PointWriteAnalogue(pnt *Point, writeValue float32) error {
	if device.isPointFloat(pnt) {

	}
	write := &Write{
		ObjectID:   pnt.ObjectID,
		ObjectType: pnt.ObjectType,
		Prop:       btypes.PropPresentValue,
		WriteValue: writeValue,
	}
	err := device.Write(write)
	if err != nil {
		return err
	}
	return nil
}

//PointWriteBool use this when wanting to write a new value for an BV, AO
func (device *Device) PointWriteBool(pnt *Point, writeValue uint32) error {
	if device.isPointFloat(pnt) {

	}
	write := &Write{
		ObjectID:   pnt.ObjectID,
		ObjectType: pnt.ObjectType,
		Prop:       btypes.PropPresentValue,
		WriteValue: writeValue,
	}
	err := device.Write(write)
	if err != nil {
		return err
	}
	return nil
}

/*
***** HELPERS *****
 */

func (device *Device) toFloat(d btypes.PropertyData) float64 {
	_, out := data.ToFloat64(d)
	return out
}

func (device *Device) ToBitString(d btypes.PropertyData) *btypes.BitString {
	_, out := data.ToBitString(d)
	return out
}

func (device *Device) toUint32(d btypes.PropertyData) uint32 {
	_, out := data.ToUint32(d)
	return out
}

func (device *Device) toInt(d btypes.PropertyData) int {
	_, out := data.ToInt(d)
	return out
}

func (device *Device) toBool(d btypes.PropertyData) bool {
	_, out := data.ToBool(d)
	return out
}

func (device *Device) toStr(d btypes.PropertyData) string {
	_, out := data.ToStr(d)
	return out
}

func (device *Device) isPointWriteable(pnt *Point) (ok bool) {
	if pnt.ObjectType != btypes.BinaryOutput {
		return true
	}
	if pnt.ObjectType != btypes.BinaryValue {
		return true
	}
	if pnt.ObjectType != btypes.AnalogOutput {
		return true
	}
	if pnt.ObjectType != btypes.AnalogOutput {
		return true
	}
	if pnt.ObjectType != btypes.MultiStateOutput {
		return true
	}
	if pnt.ObjectType != btypes.MultiStateValue {
		return true
	}
	return false
}

func (device *Device) isPointFloat(pnt *Point) (ok bool) {
	if pnt.ObjectType == btypes.AnalogInput {
		return true
	}
	if pnt.ObjectType == btypes.AnalogOutput {
		return true
	}
	if pnt.ObjectType == btypes.AnalogValue {
		return true
	}
	return false
}

func (device *Device) isPointBool(pnt *Point) (ok bool) {
	if pnt.ObjectType == btypes.BinaryInput {
		return true
	}
	if pnt.ObjectType == btypes.BinaryOutput {
		return true
	}
	if pnt.ObjectType == btypes.BinaryValue {
		return true
	}
	return false
}
