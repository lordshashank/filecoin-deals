package dealutils

import (
	configurations "fildeal/src/config"
	"fildeal/src/types"
	"fmt"
	"os/exec"
	"strconv"
)

func InitiateDeal(fileName string, storageProvider string, pieceSize uint64, commpCid string, carFileSize uint64, flags types.DealFlags) error {

	// Payload cid is irrelevant for boost deal creation, putting dummy value
	payloadCid := "bafkreibtkdcncmofmavpdsar6msrmb2h4d7oetwtwtkz5cv3zsnwoyrrfq"
	command := "boost"
	var url string
	var verified string
	if flags.Testnet {
		url = configurations.LoadConfigurations().LighthouseDownloadURL + fileName
		verified = "true"
	} else {
		url = "http://localhost:8000/download/car?file_name=" + fileName + ".data"
		verified = "false"
	}
	args := []string{
		"deal",
		"--provider=" + storageProvider,
		// "--http-url=" +  "http://localhost:8000/download/car?file_name=" + fileName + ".data",
		"--http-url=" +  url,
		"--commp=" + commpCid,
		"--car-size=" + strconv.Itoa(int(carFileSize)),
		"--piece-size=" + strconv.Itoa(int(pieceSize)),
		"--payload-cid=" + payloadCid,
		"--duration=3542400",
		"--storage-price=0" ,
		"--verified=" + verified,
	}

	fmt.Println("Running command: ", command, args)
	dealResponse, err := exec.Command(command, args...).Output()

	if err != nil {
		fmt.Println(dealResponse)
		return fmt.Errorf("failed to initiate deal: %w", err)
	}

	fmt.Println("Deal initiated successfully for: " + fileName)


	fmt.Println("Deal Response: ", string(dealResponse))
	return nil
}
