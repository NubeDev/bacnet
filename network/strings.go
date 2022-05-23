package network

//ReadString to read a string like objectName
func (device *Device) ReadString(obj *Object) (string, error) {
	read, err := device.Read(obj)
	if err != nil {
		return "", err
	}
	return device.toStr(read), nil
}
