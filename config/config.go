package i2pgateconfig

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ipfs/go-ipfs-util"
	"github.com/mitchellh/go-homedir"
	ma "github.com/multiformats/go-multiaddr"

	//serialize "github.com/ipfs/go-ipfs-config/serialize"
	serialize "github.com/ipfs/go-ipfs-config/serialize"
)

// Config is a struct very alike the one used to configure IPFS which is used
//to create, load, and access i2p configurations.
type Config struct {
	SAMHost                            string
	SAMPort                            string
	TunName                            string
	AddressRPC                         string
	AddressHTTP                        string
	EncryptLeaseSet                    bool
	EncryptedLeaseSetKey               string
	EncryptedLeaseSetPrivateKey        string
	EncryptedLeaseSetPrivateSigningKey string
	InAllowZeroHop                     bool
	OutAllowZeroHop                    bool
	InLength                           int
	OutLength                          int
	InQuantity                         int
	OutQuantity                        int
	InVariance                         int
	OutVariance                        int
	InBackupQuantity                   int
	OutBackupQuantity                  int
	UseCompression                     bool
	FastRecieve                        bool
	ReduceIdle                         bool
	ReduceIdleTime                     int
	ReduceIdleQuantity                 int
	CloseIdle                          bool
	CloseIdleTime                      int
	AccessListType                     string
	AccessList                         []string
	MessageReliability                 string
	OnlyI2P                            bool
	ListenerBase32                     string
	ListenerBase64                     string
}

func (c *Config) accesslisttype() string {
	if c.AccessListType == "whitelist" {
		return "i2cp.enableAccessList=true"
	} else if c.AccessListType == "blacklist" {
		return "i2cp.enableBlackList=true"
	} else if c.AccessListType == "none" {
		return ""
	}
	return ""
}

func (c *Config) accesslist() string {
	if c.AccessListType != "" && len(c.AccessList) > 0 {
		r := ""
		for _, s := range c.AccessList {
			r += s + ","
		}
		if r != "" {
			return "i2cp.accessList=" + strings.TrimSuffix(r, ",")
		}
	}
	return ""
}

// Print returns and prints a formatted list of configured tunnel settings.
func (c *Config) Print() []string {
	confstring := []string{
		"inbound.length=" + strconv.Itoa(c.InLength),
		"outbound.length=" + strconv.Itoa(c.OutLength),
		"inbound.lengthVariance=" + strconv.Itoa(c.InVariance),
		"outbound.lengthVariance=" + strconv.Itoa(c.OutVariance),
		"inbound.backupQuantity=" + strconv.Itoa(c.InBackupQuantity),
		"outbound.backupQuantity=" + strconv.Itoa(c.OutBackupQuantity),
		"inbound.quantity=" + strconv.Itoa(c.InQuantity),
		"outbound.quantity=" + strconv.Itoa(c.OutQuantity),
		"inbound.allowZeroHop=" + strconv.FormatBool(c.InAllowZeroHop),
		"outbound.allowZeroHop=" + strconv.FormatBool(c.OutAllowZeroHop),
		"i2cp.encryptLeaseSet=" + strconv.FormatBool(c.EncryptLeaseSet),
		"i2cp.gzip=" + strconv.FormatBool(c.UseCompression),
		"i2cp.reduceOnIdle=" + strconv.FormatBool(c.ReduceIdle),
		"i2cp.reduceIdleTime=" + strconv.Itoa(c.ReduceIdleTime),
		"i2cp.reduceQuantity=" + strconv.Itoa(c.ReduceIdleQuantity),
		"i2cp.closeOnIdle=" + strconv.FormatBool(c.CloseIdle),
		"i2cp.closeIdleTime=" + strconv.Itoa(c.CloseIdleTime),
		c.accesslisttype(),
		c.accesslist(),
	}

	log.Println(confstring)
	return confstring
}

func (c *Config) SAMAddr() string {
	return c.SAMHost + ":" + c.SAMPort
}

func (c *Config) SAMMultiaddr() (ma.Multiaddr, error) {
	return ma.NewMultiaddr(c.SAMAddr())
}

const (
	// DefaultPathName is the default config dir name
	DefaultPathName = ".ipfs"
	// DefaultPathRoot is the path to the default config dir location.
	DefaultPathRoot = "~/" + DefaultPathName
	// DefaultConfigFile is the filename of the configuration file
	DefaultConfigFile = "i2pconfig"
	// EnvDir is the environment variable used to change the path root.
	EnvDir = "IPFS_PATH"
)

func Init(out io.Writer) (*Config, error) {
	cfg := &Config{
		SAMHost:                            "/ip4/127.0.0.1/",
		SAMPort:                            "/tcp/7656/",
		TunName:                            "ipfs",
		AddressRPC:                         "/ip4/127.0.0.1/tcp/4001/",
		AddressHTTP:                        "/ip4/127.0.0.1/tcp/5001/",
		EncryptLeaseSet:                    false,
		EncryptedLeaseSetKey:               "",
		EncryptedLeaseSetPrivateKey:        "",
		EncryptedLeaseSetPrivateSigningKey: "",
		InAllowZeroHop:                     false,
		OutAllowZeroHop:                    false,
		InLength:                           3,
		OutLength:                          3,
		InQuantity:                         3,
		OutQuantity:                        3,
		InVariance:                         0,
		OutVariance:                        0,
		InBackupQuantity:                   1,
		OutBackupQuantity:                  1,
		UseCompression:                     true,
		FastRecieve:                        true,
		ReduceIdle:                         true,
		ReduceIdleQuantity:                 1,
		CloseIdle:                          false,
		AccessListType:                     "none",
		AccessList:                         []string{""},
		OnlyI2P:                            false,
		ListenerBase32:                     "",
		ListenerBase64:                     "",
	}
	return cfg, nil
}

// ConfigAt loads an i2p gateway plugin from the IPFS_PATH directory. It's a
// file intended to be as similar to the IPFS config as possible.
func ConfigAt(ipfs_path string) (*Config, error) {
	var final error
	if filename, final := Filename(ipfs_path); final == nil {
		return Load(filename)
	}
	return nil, final
}

// Filename returns the correct path to the config file for consumption by other
// parts of the application
func Filename(ipfs_path string) (string, error) {
	return Path(ipfs_path, DefaultConfigFile)
}

func Load(filename string) (*Config, error) {
	// if nothing is there, generate a 'safe(paranoid)' default config and
	// inform the user thusly
	if !util.FileExists(filename) {
		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		log.Println("i2p Gateway tunnel configuration initialized in: ", filename)
		cfg, err := Init(f)
		if err != nil {
			return nil, err
		}
		return cfg, serialize.WriteConfigFile(filename, cfg)
	}

	var cfg Config
	err := serialize.ReadConfigFile(filename, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) Save(ipfs_path string) (*Config, error) {
	var filename string
	var err error
	if filename, err = Filename(ipfs_path); err != nil {
		return nil, err
	}

	// if nothing is there, generate a 'safe(paranoid)' default config and
	// inform the user thusly
	if util.FileExists(filename) {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		log.Println("i2p Gateway tunnel configuration saved in: ", filename)
		return cfg, serialize.WriteConfigFile(filename, cfg)
	}
	return Load(filename)

}

func Path(configroot, extension string) (string, error) {
	if len(configroot) == 0 {
		dir, err := PathRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, extension), nil

	}
	return filepath.Join(configroot, extension), nil
}

// PathRoot returns the default configuration root directory
func PathRoot() (string, error) {
	dir := os.Getenv(EnvDir)
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(DefaultPathRoot)
	}
	return dir, err
}

func AddressRPC(addr string, cfg interface{}) error {
	cfg.(*Config).AddressRPC = addr
	return nil
}

func AddressHTTP(addr string, cfg interface{}) error {
	cfg.(*Config).AddressHTTP = addr
	return nil
}

func ListenerBase32(addr string, cfg interface{}) error {
	cfg.(*Config).ListenerBase32 = addr
	return nil
}

func ListenerBase64(addr string, cfg interface{}) error {
	cfg.(*Config).ListenerBase64 = addr
	return nil
}
