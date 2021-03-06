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

	serialize "gx/ipfs/QmTbcMKv6GU3fxhnNcbzYChdox9Fdd7VpucM3PQ7UWjX3D/go-ipfs-config/serialize"
)

// Config is a struct very alike the one used to configure IPFS which is used
//to create, load, and access i2p configurations.
type Config struct {
	SAMHost                            string
	SAMPort                            string
	TunName                            string
	AddressRPC                         string
	AddressHTTP                        string
	AddressSwarm                       string
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
	ListenerSK                         string
	ListenerBase32RPC                  string
	ListenerBase64RPC                  string
	ListenerSKRPC                      string
	ListenerBase32Swarm                string
	ListenerBase64Swarm                string
	ListenerSKSwarm                    string
	I2PBootstrapAddresses              []string
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

// TargetHTTP returns the string representation of the IPFS gateway in host:port
// form(not a multiaddr)
func (c *Config) TargetHTTP() string {
	h, _ := c.HTTPHost()
	p, _ := c.HTTPPort()
	return h + ":" + p
}

// HTTPHost returns the string representation of the IPFS gateway host only
func (c *Config) HTTPHost() (string, error) {
	temp, err := c.MaTargetHTTP()
	if err != nil {
		log.Println("Failed to get HTTP Multiaddr", err.Error())
		return "", err
	}
	log.Println("Getting value for HTTP host IPv4 protocol")
	return temp.ValueForProtocol(ma.P_IP4)
}

// HTTPPort returns the string representation of the IPFS gateway port only
func (c *Config) HTTPPort() (string, error) {
	temp, err := c.MaTargetHTTP()
	if err != nil {
		log.Println("Failed to get HTTP Multiaddr", err.Error())
		return "", err
	}
	log.Println("Value for HTTP port IPv4 protocol")
	return temp.ValueForProtocol(ma.P_TCP)
}

// TargetRPC returns the string representation of the RPC gateway in host:port
// form(not a multiaddr)
func (c *Config) TargetRPC() string {
	h, _ := c.RPCHost()
	p, _ := c.RPCPort()
	return h + ":" + p
}

// RPCHost returns the string representation of the RPC gateway host only
func (c *Config) RPCHost() (string, error) {
	temp, err := c.MaTargetRPC()
	if err != nil {
		log.Println("Failed to get RPC Multiaddr", err.Error())
		return "", err
	}
	log.Println("Value for RPC host IPv4 protocol")
	return temp.ValueForProtocol(ma.P_IP4)
}

// RPCPort returns the string representation of the RPC gateway port only
func (c *Config) RPCPort() (string, error) {
	temp, err := c.MaTargetRPC()
	if err != nil {
		log.Println("Failed to get RPC Multiaddr", err.Error())
		return "", err
	}
	log.Println("Value for RPC port IPv4 protocol")
	return temp.ValueForProtocol(ma.P_TCP)
}

func (c *Config) TargetSwarm() string {
	h, _ := c.SwarmHost()
	p, _ := c.SwarmPort()
	return h + "" + p
}

func (c *Config) SwarmHost() (string, error) {
	temp, err := c.MaTargetSwarm()
	if err != nil {
		log.Println("Failed to get Swarm Multiaddr", err.Error())
		return "", err
	}
	log.Println("Value for Swarm host IPv4 protocol")
	return temp.ValueForProtocol(ma.P_IP4)
}

func (c *Config) SwarmPort() (string, error) {
	temp, err := c.MaTargetSwarm()
	if err != nil {
		log.Println("Failed to get RPC Multiaddr", err.Error())
		return "", err
	}
	log.Println("Value for Swarm port IPv4 protocol")
	return temp.ValueForProtocol(ma.P_TCP)
}

func (c *Config) MaTargetHTTP() (ma.Multiaddr, error) {
	log.Println("Detected HTTP address:", c.AddressHTTP)
	return ma.NewMultiaddr(Unquote(c.AddressHTTP))
}

func (c *Config) MaTargetRPC() (ma.Multiaddr, error) {
	log.Println("Detected RPC address:", c.AddressRPC)
	return ma.NewMultiaddr(Unquote(c.AddressRPC))
}

func (c *Config) MaTargetSwarm() (ma.Multiaddr, error) {
	log.Println("Detected Swarm Address:", c.AddressSwarm)
	return ma.NewMultiaddr(Unquote(c.AddressSwarm))
}

func (c *Config) HostSAM() string {
	m, _ := c.SAMMultiaddr()
	at, _ := m.ValueForProtocol(ma.P_IP4)
	log.Println("SAM Host:", at)
	return at
}

func (c *Config) PortSAM() string {
	m, _ := c.SAMMultiaddr()
	at, _ := m.ValueForProtocol(ma.P_TCP)
	log.Println("SAM Port:", at)
	return at
}

func (c *Config) SAMAddr() string {
	rt := strings.Replace(c.SAMHost+c.SAMPort, "//", "/", -1)
	return rt
}

func (c *Config) SAMMultiaddr() (ma.Multiaddr, error) {
	return ma.NewMultiaddr(c.SAMAddr())
}

func (c *Config) BootstrapAddresses() []string {
	return c.I2PBootstrapAddresses
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
		AddressSwarm:                       "/ip4/127.0.0.1/tcp/4001/",
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
		ReduceIdleTime:                     2000000,
		ReduceIdleQuantity:                 1,
		CloseIdle:                          false,
		AccessListType:                     "none",
		AccessList:                         []string{""},
		OnlyI2P:                            false,
		ListenerBase32:                     "",
		ListenerBase64:                     "",
		ListenerSK:                         "",
		ListenerBase32RPC:                  "",
		ListenerBase64RPC:                  "",
		ListenerSKRPC:                      "",
		ListenerBase32Swarm:                "",
		ListenerBase64Swarm:                "",
		ListenerSKSwarm:                    "",
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

// Load reads a config file, or if one does not exist, initializes one.
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
	log.Println("i2p Gateway tunnel configuration found in: ", filename)
	var cfg Config
	err := serialize.ReadConfigFile(filename, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes a config file.
func (cfg *Config) Save(ipfs_path string) (*Config, error) {
	var filename string
	var err error
	if filename, err = Filename(ipfs_path); err != nil {
		return nil, err
	}

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

//
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

func AddressSwarm(addr string, cfg interface{}) error {
	cfg.(*Config).AddressSwarm = addr
	return nil
}

func AddressesBootstrap(addr []string, cfg interface{}) error {
	for _, v := range addr {
		cfg.(*Config).I2PBootstrapAddresses = append(cfg.(*Config).I2PBootstrapAddresses, v)
	}
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

func ListenerBase32RPC(addr string, cfg interface{}) error {
	cfg.(*Config).ListenerBase32RPC = addr
	return nil
}

func ListenerBase64RPC(addr string, cfg interface{}) error {
	cfg.(*Config).ListenerBase64RPC = addr
	return nil
}

func ListenerBase32Swarm(addr string, cfg interface{}) error {
	cfg.(*Config).ListenerBase32Swarm = addr
	return nil
}

func ListenerBase64Swarm(addr string, cfg interface{}) error {
	cfg.(*Config).ListenerBase64Swarm = addr
	return nil
}

func Unquote(s string) string {
	return strings.Replace(s, "\"", "", -1)
}
