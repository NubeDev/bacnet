package cmd

import (
	"fmt"
	"github.com/NubeDev/bacnet/datalink"
	"strconv"

	"github.com/spf13/viper"

	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// Flags
var (
	deviceID       int
	objectID       int
	objectType     int
	arrayIndex     uint32
	propertyType   string
	listProperties bool
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

	dataLink, err := datalink.NewUDPDataLink(viper.GetString("interface"), viper.GetInt("port"))
	if err != nil {
		log.Fatal(err)
	}
	c := bacnet.NewClient(dataLink, 0)
	defer c.Close()
	go c.Run()
	//raddr := net.TCPAddr{IP: net.IPv4(151, 101, 1, 69), Port: int(80)}
	// We need the actual address of the device first.
	//resp, err := c.WhoIs(deviceID, deviceID)
	//if err != nil {
	//	log.Fatal(err)
	//}
	////
	////if len(resp) == 0 {
	////	log.Fatal("Device id was not found on the network.")
	////}
	//
	//dest := resp[0]
	//dest.Addr.Adr = []byte{4}
	//dest.Addr.Adr = []uint8{4}
	//addr := btypes.Address{
	//	Net: 4,
	//	Mac: []uint8{192, 168, 15, 20, 186, 192},
	//	Adr: []uint8{4},
	//	Len: 1,
	//}

	//TODO to get build mac address for network IP:port
	dlink, _ := datalink.NewUDPDataLink(Interface, 47899)
	fmt.Println(dlink.GetMyAddress())
	fmt.Println(Interface)

	addr := btypes.Address{
		Net: 4,
		Mac: []uint8{192, 168, 15, 20, 186, 192},
		Adr: []uint8{4},
		Len: 1,
	}
	object := btypes.ObjectID{
		Type:     8,
		Instance: 1103,
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
				btypes.Property{
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
	readCmd.Flags().IntVarP(&objectID, "objectID", "o", 1234, "object ID")
	readCmd.Flags().IntVarP(&objectType, "objectType", "j", 8, "object type")
	readCmd.Flags().StringVarP(&propertyType, "property", "t",
		btypes.ObjectNameStr, propertyTypeDescr)

	readCmd.Flags().Uint32Var(&arrayIndex, "index", bacnet.ArrayAll, "Which position to return.")

	readCmd.PersistentFlags().BoolVarP(&listProperties, "list", "l", false,
		listPropertiesDescr)
}
