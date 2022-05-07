package cmd

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/datalink"
	"github.com/NubeDev/bacnet/helpers/data"
	ip2bytes "github.com/NubeDev/bacnet/helpers/ipbytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// readMultiCmd represents the readMultiCmd command
var readMultiCmd = &cobra.Command{
	Use:   "multi",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: readMulti,
}

func readMulti(cmd *cobra.Command, args []string) {
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
		Type:     btypes.DeviceType,
		Instance: btypes.ObjectInstance(deviceID),
	}

	//get max adpu len
	rp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.DeviceType,
				Instance: btypes.ObjectInstance(deviceID),
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropMaxAPDU,
					ArrayIndex: bacnet.ArrayAll,
				},
			},
		},
	}

	dest := btypes.Device{
		ID:   object,
		Addr: addr,
	}
	// get the device MaxApdu
	out, err := c.ReadProperty(dest, rp)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, dest.MaxApdu = data.ToUint32(out)

	fmt.Println("MaxApdu", dest.MaxApdu)

	rp = btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     8,
				Instance: btypes.ObjectInstance(deviceID),
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropObjectList,
					ArrayIndex: bacnet.ArrayAll,
				},
			},
		},
	}

	// get the device object list
	out, err = c.ReadProperty(dest, rp)
	if err != nil {
		log.Fatal(err)
		return
	}

	ids, ok := out.Object.Properties[0].Data.([]interface{})
	if !ok {
		fmt.Println("unable to get object list")
		return
	}
	fmt.Println(len(ids), "LEN")
	rpm := btypes.MultiplePropertyData{}

	rpm.Objects = []btypes.Object{
		btypes.Object{
			ID: btypes.ObjectID{
				Type:     8,
				Instance: btypes.ObjectInstance(deviceID),
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropObjectList,
					ArrayIndex: bacnet.ArrayAll,
				},
			},
		},
		btypes.Object{
			ID: btypes.ObjectID{
				Type:     8,
				Instance: btypes.ObjectInstance(deviceID),
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropObjectList,
					ArrayIndex: bacnet.ArrayAll,
				},
			},
		},
	}

	rpmRes, err := c.ReadMultiProperty(dest, rpm)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(rpmRes)

	//for i, objs := range ids {
	//	id, ok := objs.(btypes.ObjectID)
	//	if !ok {
	//		log.Printf("unable to read object id %v\n", objs)
	//		return
	//	}
	//
	//	if id.Type == btypes.AnalogOutput {
	//
	//		props := []btypes.Property{
	//			btypes.Property{
	//				Type:       btypes.PropObjectName,
	//				ArrayIndex: bacnet.ArrayAll,
	//			},
	//			btypes.Property{
	//				Type:       btypes.PropDescription,
	//				ArrayIndex: bacnet.ArrayAll,
	//			},
	//		}
	//		rpm.Objects[i].Properties = append(rpm.Objects[i].Properties, props...)
	//	}
	//
	//}
	//fmt.Println(rpm)
	//rpmRes, err := c.ReadMultiProperty(dest, rpm)
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(rpmRes)
}

func init() {
	RootCmd.AddCommand(readMultiCmd)
	readMultiCmd.PersistentFlags().IntVarP(&deviceID, "device", "d", 1234, "device id")
	readMultiCmd.Flags().StringVarP(&deviceIP, "address", "", "192.168.15.202", "device ip")
	readMultiCmd.Flags().IntVarP(&devicePort, "dport", "", 47808, "device port")
	readMultiCmd.Flags().IntVarP(&networkNumber, "network", "", 0, "bacnet network number")
	readMultiCmd.Flags().IntVarP(&deviceHardwareMac, "mstp", "", 0, "device hardware mstp addr")

}
