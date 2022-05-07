package cmd

import (
	"encoding/json"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/datalink"
	"github.com/spf13/cobra"
	"log"
	"os"
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

	dataLink, err := datalink.NewUDPDataLink(Interface, Port)
	if err != nil {
		log.Fatal(err)
	}
	c := bacnet.NewClient(dataLink, 0)
	defer c.Close()
	go c.Run()
	wh := &bacnet.WhoIsBuilder{
		GlobalBroadcast: true,
		NetworkNumber:   0,
	}
	wh.Low = startRange
	wh.High = endRange
	ids, err := c.WhoIs(wh)
	if err != nil {
		log.Fatal(err)
	}

	ioWriter := os.Stdout
	// Check to see if a file was passed to us
	if len(outputFilename) > 0 {
		ioWriter, err = os.Create(outputFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer ioWriter.Close()
	}
	// Pretty Print!
	w := json.NewEncoder(ioWriter)
	w.SetIndent("", "    ")
	w.Encode(ids)

}

func init() {
	RootCmd.AddCommand(whoIsCmd)
	whoIsCmd.Flags().IntVarP(&startRange, "start", "s", -1, "Start range of discovery")
	whoIsCmd.Flags().IntVarP(&endRange, "end", "e", int(0xBAC0), "End range of discovery")
	whoIsCmd.Flags().StringVarP(&outputFilename, "out", "o", "", "Output results into the given filename in json structure.")
}
