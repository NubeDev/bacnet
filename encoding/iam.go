package encoding

import (
	"github.com/NubeDev/bacnet/btypes"
)

func (enc *Encoder) IAm(id btypes.IAm) error {
	enc.AppData(id.ID, false)
	enc.AppData(id.MaxApdu, false)
	enc.AppData(id.Segmentation, false)
	enc.AppData(id.Vendor, false)
	return enc.Error()
}

func (d *Decoder) IAm(id *btypes.IAm) error {
	objID, err := d.AppData()
	if err != nil {
		return err
	}
	if i, ok := objID.(btypes.ObjectID); ok {
		id.ID = i
	}
	maxapdu, _ := d.AppData()
	if m, ok := maxapdu.(uint32); ok {
		id.MaxApdu = m
	}
	segmentation, _ := d.AppData()
	if m, ok := segmentation.(uint32); ok {
		id.Segmentation = btypes.Enumerated(m)
	}
	vendor, err := d.AppData()
	if v, ok := vendor.(uint32); ok {
		id.Vendor = v
	}
	return d.Error()
}
