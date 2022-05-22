package cmd

import (
	"fmt"
	"github.com/NubeDev/bacnet"
	pprint "github.com/NubeDev/bacnet/helpers/print"
	"github.com/NubeDev/bacnet/local"
	"github.com/spf13/cobra"
)

// Flags
var startRange int
var endRange int

var outputFilename string

// whoIsCmd represents the whoIs command
var whoIsCmd = &cobra.Command{
	Use:   "whois",
	Short: "BACnet device discovery",
	Long: `whoIs does a bacnet network discovery to find devices in the network
 given the provided range.`,
	Run: main,
}

func main(cmd *cobra.Command, args []string) {

	client, err := local.New(&local.Local{Interface: Interface, Port: Port})
	if err != nil {
		fmt.Println("ERR-client", err)
		return
	}
	defer client.ClientClose()
	go client.ClientRun()

	wi := &bacnet.WhoIsOpts{
		High:            endRange,
		Low:             startRange,
		GlobalBroadcast: true,
		NetworkNumber:   uint16(networkNumber),
	}

	whoIs, err := client.Whois(wi)
	if err != nil {
		fmt.Println("ERR-whoIs", err)
		return
	}

	pprint.PrintJOSN(whoIs)
}

func init() {
	RootCmd.AddCommand(whoIsCmd)
	whoIsCmd.Flags().IntVarP(&startRange, "start", "s", -1, "Start range of discovery")
	whoIsCmd.Flags().IntVarP(&endRange, "end", "e", int(0xBAC0), "End range of discovery")
	whoIsCmd.Flags().IntVarP(&networkNumber, "network", "", 0, "network number")
	whoIsCmd.Flags().StringVarP(&outputFilename, "out", "o", "", "Output results into the given filename in json structure.")
}
