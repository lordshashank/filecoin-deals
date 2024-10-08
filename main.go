package main

import (
	"fmt"
	"io"
	"os"

	configurations "fildeal/src/config"
	"fildeal/src/deal"
	"fildeal/src/index"
	mkpiece "fildeal/src/mkpiece"
	"fildeal/src/server"
	"fildeal/src/types"
	"fildeal/src/utils"
)

func main() {

    // Detailed usage information
    usage := `Usage: fildeal <command> [arguments]
    Commands:
    cmp <parentFile> <childFile>       Compare two files and find the offset of the child file in the parent file.
    generate <files...>                Generate a data segment piece from the given files and output it to stdout.
    splitpiece <file> <outputDir>      Split the specified file into pieces and save them in the output directory.
    initiate <inputFolder> <miner>     Initiate a deal with the specified input folder and miner.
    boost-index <file>                 Parse and index a file similar to Boost.

    Examples:
    fildeal cmp a.car b.car
    fildeal generate a.car b.car c.car > out.dat
    fildeal splitpiece input.car outputDir
    fildeal initiate inputFolder miner [--server]
    fildeal boost-index file.car
    `

    // Check for --help flag
    if len(os.Args) < 2 || os.Args[1] == "--help" {
        fmt.Println(usage)
        return
    }

    command := os.Args[1]

    switch command {
    case "cmp":
        if len(os.Args) != 4 {
            fmt.Println("Usage: fildeal cmp <parentFile> <childFile>")
            return
        }
        fileAPath := os.Args[2]
        fileBPath := os.Args[3]
        offset, err := utils.FindOffset(fileAPath, fileBPath)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Printf("Child file starts at offset %d in parent\n", offset)

    case "generate":
        if len(os.Args) < 3 {
            fmt.Println("Usage: fildeal generate <files...> > out.dat")
            return
        }
        readers := make([]io.ReadSeeker, 0)
        for _, arg := range os.Args[2:] {
            r, err := os.Open(arg)
            if err != nil {
                panic(err)
            }
            readers = append(readers, r)
        }
        out := mkpiece.MakeDataSegmentPiece(readers)
        _, err := io.Copy(os.Stdout, out)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        for _, reader := range readers {
            // Check if the reader is fully consumed
            if _, err := reader.Read(make([]byte, 1)); err != io.EOF {
                panic(err)
            }
        }

    case "splitpiece":
        if len(os.Args) != 3 {
            fmt.Println("Usage: fildeal splitpiece <file> <outputDir>")
            return
        }
        filePath := os.Args[2]
        outputDir := os.Args[3]
        err := mkpiece.SplitPiece(filePath, outputDir)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

    case "initiate":
        if len(os.Args) < 4 {
            fmt.Println("Usage: fildeal initiate <inputFolder> <miner> [--server] [--testnet]")
            return
        }
        inputFolder := os.Args[2]
        miner := os.Args[3]
    
        // Initialize flags
        flags := types.DealFlags{}
    
        // Parse additional flags
        for _, arg := range os.Args[4:] {
            switch arg {
            case "--testnet":
                flags.Testnet = true
            case "--server":
                flags.Server = true
            default:
                fmt.Printf("Unknown flag: %s\n", arg)
                return
            }
        }
    
        err := deal.MakeDeal(inputFolder, miner, flags)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        if flags.Server {
            config := configurations.Configurations{Port: configurations.LoadConfigurations().Port} // Example configuration
            handler := server.SetupRouter()
            server.StartServer(config, handler)
        }

    case "boost-index":
        if len(os.Args) != 3 {
            fmt.Println("Usage: fildeal boost-index <file>")
            return
        }
        filePath := os.Args[2]
        err := index.BoostIndex(filePath)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

    default:
        fmt.Println("Unknown command:", command)
		fmt.Println(usage)
    }
}



