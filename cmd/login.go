package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emdneto/otsgo/client"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

var res bool
var auth client.AuthYaml

func isTTY(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func promptEndpoint() (host, baseUri string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select default endpoint:")
	fmt.Println("  1) EU      https://eu.onetimesecret.com  (default)")
	fmt.Println("  2) US      https://onetimesecret.com")
	fmt.Println("  3) Custom")
	fmt.Print("Choice [1]: ")
	line, _ := reader.ReadString('\n')

	switch strings.TrimSpace(line) {
	case "2":
		return "https://onetimesecret.com", "https://onetimesecret.com/api"
	case "3":
		fmt.Print("Host URL (e.g. https://ots.example.com): ")
		h, _ := reader.ReadString('\n')
		h = strings.TrimRight(strings.TrimSpace(h), "/")
		if h == "" {
			fmt.Println("Empty host, falling back to EU default")
			return "https://eu.onetimesecret.com", "https://eu.onetimesecret.com/api"
		}
		return h, h + "/api"
	default:
		return "https://eu.onetimesecret.com", "https://eu.onetimesecret.com/api"
	}
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Perform basic http auth and store credentials",
	Run: func(cmd *cobra.Command, args []string) {

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		host := client.HOST
		baseUri := client.BASE_URI
		if isTTY(os.Stdin) {
			host, baseUri = promptEndpoint()
			client.HOST = host
			client.BASE_URI = baseUri
		}

		if len(username) != 0 && len(password) != 0 {
			fmt.Printf("WARNING! Your password will be stored unencrypted in %s\n", "~/.otsgo.yaml")
			fmt.Printf("\n")
			fmt.Printf("If you don't want to keep the credentials in the configuration file, use the environment variables like this: \n\nexport OTS_USER=demo; export OTS_TOKEN=demo\n\n")
			AuthInfo = client.Auth{
				Username: username,
				Password: password,
				Enabled:  true,
			}
		}

		res = client.Login(AuthInfo)
		if res {
			fmt.Printf("Login Succeeded\n")
			auth = client.AuthYaml{
				Username: AuthInfo.Username,
				Password: AuthInfo.Password,
				Host:     host,
				BaseUri:  baseUri,
			}
		} else {
			fmt.Printf("Login failed\n")
			return
		}

		yamlData, err := yaml.Marshal(&auth)
		if err != nil {
			fmt.Printf("Error while Marshaling. %v", err)
		}

		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		fileName := fmt.Sprintf("%s/.otsgo.yaml", dirname)
		err = os.WriteFile(fileName, yamlData, 0600)
		if err != nil {
			panic("Unable to write data into the file")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().StringP("username", "u", "", "OTS Username")
	loginCmd.PersistentFlags().StringP("password", "p", "", "OTS Token")
	loginCmd.PersistentFlags().BoolP("password-stdin", "", false, "Take the API Token from stdin")

}
