package licserver

import (
	"fmt"
	"io/ioutil"
	"bufio"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

var Configs []Config

func init() {
	caddy.RegisterPlugin("licserver", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

type Config struct {
	Dbname string
}

type Licserver struct {
	Next    httpserver.Handler
	Configs []Config
}

// setup configures a new gzip middleware instance.
func setup(c *caddy.Controller) error {
	var err error
	Configs, err = parse(c)
	if err != nil {
		return err
	}

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Licserver{Next: next, Configs: Configs}
	})

	return nil
}

func (h Licserver) ServeHTTP(response http.ResponseWriter, request *http.Request) (int, error) {

	for _, c := range h.Configs {
		//LICENSE:EXPDATE:EMAIL:HWID = 0,1,2
		fmt.Println("License Server")
		if !CheckFileExist(c.Dbname) {
			fmt.Println("Database does not exist, creating new database.")
			_ = createFile(c.Dbname)
		}
		database, _ := readLines(c.Dbname)
		fmt.Println("Total Licenses:", len(database))

		request.ParseForm()
		license := request.FormValue("license")
		hwid := request.FormValue("hwid")
		//println(license, hwid)

		database, _ = readLines(c.Dbname)

		for _, table := range database {

			row := strings.Split(table, ":")

			t, err := time.Parse("2006-01-02", row[1])
			if err != nil {
				fmt.Println("ERROR: Error reading database")
			}

			t2, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

			if license == row[0] && t.After(t2) {
				if hwid == row[3] {
					fmt.Fprintf(response, "0") //Registed, Good licnese
				} else if row[3] == "NOTSET" {
					b, err := ioutil.ReadFile(c.Dbname)
					if err != nil {
						fmt.Println("READfromCHECK")
						//os.Exit(0)
					}

					str := string(b)
					edit := row[0] + ":" + row[1] + ":" + row[2] + ":" + hwid
					res := strings.Replace(str, table, edit, -1)

					err = ioutil.WriteFile(c.Dbname, []byte(res), 0644)
					if err != nil {
						fmt.Println("WRITEfromCHECK")
						//os.Exit(0)
					}

					fmt.Fprintf(response, "2") //Registed, Good licnese
				}
			} else if license == row[0] && !t.After(t2) {
				fmt.Fprintf(response, "1") //registerd but license experied
			}
		}
	}
	return h.Next.ServeHTTP(response, request)
}

func parse(c *caddy.Controller) ([]Config, error) {
	var configs []Config
	for c.Next() { // skip the directive name
		conf := Config{}

		for c.NextBlock() {
			switch c.Val() {
			case "dbname":
				if !c.NextArg() {
					return configs, c.ArgErr()
				}
				conf.Dbname = c.Val()
			default:
				return configs, c.ArgErr()
			}
		}
		configs = append(configs, conf)
	}
	return configs, nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createFile(pathFile string) error {
	file, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func randomString(n int) string {
	var letterRunes = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//LICENSE:EXPDATE:EMAIL:HWID = 0,1,2
func checkHandler(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	license := request.FormValue("license")
	hwid := request.FormValue("hwid")

	database, _ := readLines("db")
	for _, table := range database {

		row := strings.Split(table, ":")

		t, err := time.Parse("2006-01-02", row[1])
		if err != nil {
			fmt.Println("ERROR: Error reading database")
		}

		t2, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

		if license == row[0] && t.After(t2) {
			if hwid == row[3] {
				fmt.Fprintf(response, "0") //Registed, Good licnese
			} else if row[3] == "NOTSET" {
				b, err := ioutil.ReadFile("db")
				if err != nil {
					fmt.Println("READfromCHECK")
					os.Exit(0)
				}

				str := string(b)
				edit := row[0] + ":" + row[1] + ":" + row[2] + ":" + hwid
				res := strings.Replace(str, table, edit, -1)

				err = ioutil.WriteFile("db", []byte(res), 0644)
				if err != nil {
					fmt.Println("WRITEfromCHECK")
					os.Exit(0)
				}

				fmt.Fprintf(response, "2") //Registed, Good licnese
			}
		} else if license == row[0] && !t.After(t2) {
			fmt.Fprintf(response, "1") //registerd but license experied
		}
	}
}
