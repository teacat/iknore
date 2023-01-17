package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var C Context = Context{
	Port:            "8080",
	DBUsername:      "root",
	DBPassword:      "root",
	DBHost:          "127.0.0.1",
	DBPort:          "32769",
	DBCharset:       "utf8mb4",
	DBDatabase:      "iknore",
	ProgramTimeZone: "Asia/Taipei",
}

type Context struct {
	Port            string
	DBUsername      string
	DBPassword      string
	DBHost          string
	DBPort          string
	DBCharset       string
	DBDatabase      string
	ProgramTimeZone string
}

// Config
type Config struct {
	Types         map[string]*ConfigType   `yaml:"types"`
	FormatOptions map[string]*FormatOption `yaml:"format_options"`
	Placeholder   string                   `yaml:"placeholder"`
}

// ConfigType
type ConfigType struct {
	Covers           []string                 `yaml:"covers"`
	Placeholder      string                   `yaml:"placeholder"`
	BackgroundColors []string                 `yaml:"background_colors"`
	Original         *ConfigTypeOriginal      `yaml:"original"`
	Sizes            []string                 `yaml:"sizes"`
	Formats          []string                 `yaml:"formats"`
	FormatOptions    map[string]*FormatOption `yaml:"format_options"`
}

// ConfigTypeOriginal
type ConfigTypeOriginal struct {
	Size        string `yaml:"size"`
	Format      string `yaml:"format"`
	Compression string `yaml:"compression"`
	Quality     int    `yaml:"quality"`
}

// FormatOption
type FormatOption struct {
	RelatedSizes []string `yaml:"releated_sizes"`
	Compression  string   `yaml:"compression"`
	Quality      int      `yaml:"quality"`
}

// LoadConfig
func LoadConfig() *Config {
	var c *Config
	b, err := os.ReadFile("./config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		log.Fatalln(err)
	}
	return c
}

// InitPlaceholders
func (c *Config) InitPlaceholders() map[string][]byte {
	p := make(map[string][]byte)

	//
	globalB, err := os.ReadFile(c.Placeholder)
	if err != nil {
		log.Fatalln(err)
	}
	p["*"+filepath.Ext(c.Placeholder)] = globalB
	//
	for k, v := range c.Types {

		if v.Placeholder == "" {
			p[k+filepath.Ext(v.Placeholder)] = globalB
			continue
		}
		b, err := os.ReadFile(v.Placeholder)
		if err != nil {
			log.Fatalln(err)
		}
		p[k+filepath.Ext(v.Placeholder)] = b
	}
	return p
}

func (c *Config) InitSizeAliases() map[string]map[string][2]int {
	p := make(map[string]map[string][2]int)
	//
	for k, v := range c.Types {
		p[k] = make(map[string][2]int)
		//
		for _, v := range v.Sizes {
			if !strings.Contains(v, "(") {
				continue
			}
			name := strings.Split(strings.Split(v, "(")[1], ")")[0] // 400x400 ([small])
			size := strings.Split(strings.Split(v, " ")[0], "x")    // [400x400] (small)
			//
			var w, h int
			if size[0] != "" {
				w, _ = strconv.Atoi(size[0])
			}
			if size[1] != "" {
				h, _ = strconv.Atoi(size[1])
			}
			p[k][name] = [2]int{w, h}
		}
	}
	log.Println(p)
	return p
}
