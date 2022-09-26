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

func main() {

	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signVaultAddress := signCmd.String("vaultaddress", "", "Vault Address e.g. (https://xxx.xxx.xxx.xxx:8200)")
	signInputFile := signCmd.String("input", "", "Input file to sign")
	signOutputSignature := signCmd.String("signature", "", "Where to save digital signature")
	signKey := signCmd.String("key", "", "Key to use to sign")
	//signAlgorithm := signCmd.String("Algorithm", "sha2-256", "SHA2 algorithm to use")
	signToken := signCmd.String("token", "nil", "Vault Token")

	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)
	verifyVaultAddress := verifyCmd.String("vaultaddress", "", "Vault Address e.g. (https://xxx.xxx.xxx.xxx:8200)")
	verifyInputFile := verifyCmd.String("input", "", "Input file to verify")
	verifyInputSignature := verifyCmd.String("signature", "", "Where to load digital signature")
	verifyKey := verifyCmd.String("key", "", "Key to use to verify")
	//signAlgorithm := signCmd.String("Algorithm", "sha2-256", "SHA2 algorithm to use")
	verifyToken := verifyCmd.String("token", "nil", "Vault Token")

	if len(os.Args) < 2 {
		Usage()
		os.Exit(1)
	}

	switch os.Args[1] {

	case "sign":
		signCmd.Parse(os.Args[2:])
		if signCmd.NFlag() < 5 {
			Usage()
			os.Exit(1)
		}
		c, err := vault.NewClient(*signVaultAddress,
			vault.WithCaPath(""),
			vault.WithAuthToken(*signToken),
		)

		if err != nil {
			log.Fatal(err)
		}
		transit := c.Transit()

		//key := "test123bacd"
		content, err := os.ReadFile(*signInputFile)
		if err != nil {
			log.Fatal(err)
		}
		signResponse, err := transit.Sign(*signKey, &vault.TransitSignOptions{
			Plaintext: string(content[:]),
		})
		if err != nil {
			log.Fatalf("Error occurred during signing: %v", err)
		}
		//log.Println("Signature: ", signResponse.Data.Signature)
		os.WriteFile(*signOutputSignature, []byte(signResponse.Data.Signature), 0666)
		fmt.Printf("File %s correctly signed. Sign is in file %s\n", *signInputFile, *signOutputSignature)

	case "verify":
		verifyCmd.Parse(os.Args[2:])
		if verifyCmd.NFlag() < 5 {
			Usage()
			os.Exit(1)
		}

		c, err := vault.NewClient(*verifyVaultAddress,
			vault.WithCaPath(""),
			vault.WithAuthToken(*verifyToken),
		)

		if err != nil {
			log.Fatal(err)
		}
		transit := c.Transit()

		//key := "test123bacd"
		content, err := os.ReadFile(*verifyInputFile)
		if err != nil {
			log.Fatal(err)
		}
		signature, err := os.ReadFile(*verifyInputSignature)
		if err != nil {
			log.Fatal(err)
		}
		verifyResponse, err := transit.Verify(*verifyKey, &vault.TransitVerifyOptions{
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
		//log.Println("Valid: ", verifyResponse.Data.Valid)
		//os.WriteFile(*signOutputSignature, []byte(signResponse.Data.Signature), 0666)
		//fmt.Printf("File %s correctly signed. Sign is in file %s\n", *verifyInputFile, *verifyInputSignature)

	default:
		fmt.Println("expected 'sign' or 'verify' subcommands")
		os.Exit(1)
	}

	/*
		log.Println(c.Address())

		l, err := c.TransitWithMountPoint("transit").List()
		if err != nil {
			log.Fatal()
		}
		log.Println(l)
	*/
	//const rsa4096 = "rsa-4096"

	//fmt.Println(c.Token())

	/*
		err = transit.Create(key, &vault.TransitCreateOptions{
			Exportable: vault.BoolPtr(true),
			Type:       rsa4096,
		})
		if err != nil {
			log.Fatal(err)
		}
	*/
	/*
		res, err := transit.Read(key)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("%+v\n", res.Data)
		}
	*/
	/*
		exportRes, err := transit.Export(key, vault.TransitExportOptions{
			KeyType: "encryption-key",
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v+", exportRes.Data.Keys[1])
	*/

	//decryptResponse, err := transit.Decrypt(key, &vault.TransitDecryptOptions{
	//	Ciphertext: encryptResponse.Data.Ciphertext,
	//})
	//if err != nil {
	//	log.Fatalf("Error occurred during decryption: %v", err)
	//}
	//log.Println("Plaintext: ", decryptResponse.Data.Plaintext)

}
