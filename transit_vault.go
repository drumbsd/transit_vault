package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	vault "github.com/drumbsd/vaultgo"
	color "github.com/fatih/color"
)

func Usage() {
	println("transit_vault v0.0.1")
	println("")
	println("Example:")
	println("")
	println("To sign a document/file:")
	println("")
	println("# ./transit_vault sign -key test123bacd -input main.go -signature main.go.signature -token s.0wtcFidsdcdscscsm -vaultaddress https://vault1:8200")
	println("")
	println("To verify a document/file")
	println("")
	println("# ./transit_vault verify -key test123bacd -input main.go -signature main.go.signature -token s.0wtcFidsdcdscscsm -vaultaddress https://vault1:8200")
	println("")
}

var (
	signCmd      = flag.NewFlagSet("sign", flag.ExitOnError)
	verifyCmd    = flag.NewFlagSet("verify", flag.ExitOnError)
	VaultAddress string
	InputFile    string
	Signature    string
	Key          string
	Token        string
)

func setupCommonFlags() {
	for _, fs := range []*flag.FlagSet{signCmd, verifyCmd} {
		fs.StringVar(&VaultAddress, "vaultaddress", "", "Vault Address e.g. (https://xxx.xxx.xxx.xxx:8200)")
		fs.StringVar(&InputFile, "input", "", "Input file to sign/verify")
		fs.StringVar(&Signature, "signature", "", "Where to save/read digital signature")
		fs.StringVar(&Key, "key", "", "Key to use to sign/verify")
		fs.StringVar(&Token, "token", "nil", "Vault Token")

	}
}
func main() {

	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	setupCommonFlags()
	switch os.Args[1] {
	case "sign":
		signCmd.Parse(os.Args[2:])
		if signCmd.NFlag() < 5 {
			Usage()
			os.Exit(1)
		}
		signDocument(VaultAddress, Token, InputFile, Key, Signature)
	case "verify":
		verifyCmd.Parse(os.Args[2:])
		if verifyCmd.NFlag() < 5 {
			Usage()
			os.Exit(1)
		}
		verifyDocument(VaultAddress, Token, InputFile, Key, Signature)

	default:
		log.Fatalf("[ERROR] unknown subcommand '%s', see help for more details.", os.Args[1])
	}
}

func signDocument(VaultAddress string, Token string, InputFile string, Key string, Signature string) {
	c, err := vault.NewClient(VaultAddress,
		vault.WithCaPath(""),
		vault.WithAuthToken(Token),
	)

	if err != nil {
		log.Fatal(err)
	}
	transit := c.Transit()

	content, err := os.ReadFile(InputFile)
	if err != nil {
		log.Fatal(err)
	}
	signResponse, err := transit.Sign(Key, &vault.TransitSignOptions{
		Plaintext: string(content[:]),
	})
	if err != nil {
		log.Fatalf("Error occurred during signing: %v", err)
	}

	os.WriteFile(Signature, []byte(signResponse.Data.Signature), 0666)
	fmt.Printf("File %s correctly signed. Sign is in file %s\n", InputFile, Signature)
}

func verifyDocument(VaultAddress string, Token string, InputFile string, Key string, InputSignature string) {
	c, err := vault.NewClient(VaultAddress,
		vault.WithCaPath(""),
		vault.WithAuthToken(Token),
	)

	if err != nil {
		log.Fatal(err)
	}
	transit := c.Transit()

	content, err := os.ReadFile(InputFile)
	if err != nil {
		log.Fatal(err)
	}
	signature, err := os.ReadFile(InputSignature)
	if err != nil {
		log.Fatal(err)
	}
	verifyResponse, err := transit.Verify(Key, &vault.TransitVerifyOptions{
		Plaintext: string(content[:]),
		Signature: string(signature[:]),
	})
	if err != nil {
		log.Fatalf("Error occurred during signing: %v", err)
	}
	if !verifyResponse.Data.Valid {
		fmt.Print("Sign is not valid!!!! ")
		color.Red("KO")
		os.Exit(1)
	} else {
		fmt.Print("Sign is valid!!!! ")
		color.Green("OK")
	}
}
