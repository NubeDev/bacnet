package cmd

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	ip2bytes "github.com/NubeDev/bacnet/helpers/ipbytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
)

// Flags
var (
	networkNumber     int
	deviceID          int
	deviceIP          string
	devicePort        int
	deviceHardwareMac int
	objectID          int
	objectType        int
	arrayIndex        uint32
	propertyType      string
	listProperties    bool
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Prints out a device's object's property",
	Long: `
 Given a device's object instance and selected property, we print the value
 stored there. There are some autocomplete features to try and minimize the
 amount of arguments that need to be passed, but do take into consideration
 this discovery process may cause longer reads.
	`,
	Run: readProp,
}

func readProp(cmd *cobra.Command, args []string) {
	if listProperties {
		btypes.PrintAllProperties()
		return
	}
	cb := &bacnet.ClientBuilder{
		Interface: Interface,
		Port:      Port,
	}
	c, _ := bacnet.NewClient(cb)
	defer c.Close()
	go c.ClientRun()

	ip, err := ip2bytes.New(deviceIP, uint16(devicePort))
	if err != nil {
		return
	}

	addr := btypes.Address{
		Net: uint16(networkNumber),
		Mac: ip,
		Adr: []uint8{uint8(deviceHardwareMac)},
	}
	object := btypes.ObjectID{
		Type:     8,
		Instance: btypes.ObjectInstance(deviceID),
	}

	dest := btypes.Device{
		ID:   object,
		Addr: addr,
	}

	var propInt btypes.PropertyType
	// Check to see if an int was passed
	if i, err := strconv.Atoi(propertyType); err == nil {
		propInt = btypes.PropertyType(uint32(i))
	} else {
		propInt, err = btypes.Get(propertyType)
	}

	if btypes.IsDeviceProperty(propInt) {
		objectType = 8
	}

	if err != nil {
		log.Fatal(err)
	}

	rp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.ObjectType(objectType),
				Instance: btypes.ObjectInstance(objectID),
			},
			Properties: []btypes.Property{
				{
					Type:       propInt,
					ArrayIndex: arrayIndex,
				},
			},
		},
	}
	out, err := c.ReadProperty(dest, rp)
	if err != nil {
		if rp.Object.Properties[0].Type == btypes.PropObjectList {
			log.Error("Note: PropObjectList reads may need to be broken up into multiple reads due to length. Read index 0 for array length")
		}
		log.Fatal(err)
	}
	if len(out.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return
	}
	fmt.Println(out.Object.Properties[0].Data)
}
func init() {
	// Descriptions are kept separate for legibility purposes.
	propertyTypeDescr := `type of read that will be done. Support both the
	property type as an integer or as a string. e.g. PropObjectName or 77 are both
	support. Run --list to see available properties.`
	listPropertiesDescr := `list all string versions of properties that are
	support by property flag`

	RootCmd.AddCommand(readCmd)

	// Pass flags to children
	readCmd.PersistentFlags().IntVarP(&deviceID, "device", "d", 1234, "device id")
	readCmd.Flags().StringVarP(&deviceIP, "address", "", "192.168.15.202", "device ip")
	readCmd.Flags().IntVarP(&devicePort, "dport", "", 47808, "device port")
	readCmd.Flags().IntVarP(&networkNumber, "network", "", 0, "bacnet network number")
	readCmd.Flags().IntVarP(&deviceHardwareMac, "mstp", "", 0, "device hardware mstp addr")
	readCmd.Flags().IntVarP(&objectID, "objectID", "o", 1234, "object ID")
	readCmd.Flags().IntVarP(&objectType, "objectType", "j", 8, "object type")
	readCmd.Flags().StringVarP(&propertyType, "property", "t",
		btypes.ObjectNameStr, propertyTypeDescr)

	readCmd.Flags().Uint32Var(&arrayIndex, "index", bacnet.ArrayAll, "Which position to return.")

	readCmd.PersistentFlags().BoolVarP(&listProperties, "list", "l", false,
		listPropertiesDescr)
}