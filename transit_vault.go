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
	println("# ./transit_vault -cmd sign -key test123bacd -input main.go -signature main.go.signature -token s.0wtcFidsdcdscscsm -vaultaddress https://vault1:8200")
	println("")
	println("To verify a document/file")
	println("")
	println("# ./transit_vault -cmd verify -key test123bacd -input main.go -signature main.go.signature -token s.0wtcFidsdcdscscsm -vaultaddress https://vault1:8200")
	println("")
}

func main() {

	Cmd := flag.String("cmd", "nil", "sign/verify")
	VaultAddress := flag.String("vaultaddress", "", "Vault Address e.g. (https://xxx.xxx.xxx.xxx:8200)")
	InputFile := flag.String("input", "", "Input file to sign/verify")
	Signature := flag.String("signature", "", "Where to save/read digital signature")
	Key := flag.String("key", "", "Key to use to sign/verify")

	Token := flag.String("token", "nil", "Vault Token")

	flag.Parse()
	if flag.NFlag() < 6 {
		Usage()
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	switch *Cmd {

	case "sign":

		signDocument(*VaultAddress, *Token, *InputFile, *Key, *Signature)

	case "verify":

		verifyDocument(*VaultAddress, *Token, *InputFile, *Key, *Signature)

	default:
		fmt.Println("expected 'sign' or 'verify' subcommands")
		os.Exit(1)
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
	} else {
		fmt.Print("Sign is valid!!!! ")
		color.Green("OK")
	}
}
