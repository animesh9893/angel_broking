package angel_broking

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

// API endpoints giving by brokers
const (
	URILogin            string = "rest/auth/angelbroking/user/v1/loginByPassword"
	URIUserSessionRenew string = "rest/auth/angelbroking/jwt/v1/generateTokens"
	URIUserProfile      string = "rest/secure/angelbroking/user/v1/getProfile"
	URILogout           string = "rest/secure/angelbroking/user/v1/logout"
	URIGetOrderBook     string = "rest/secure/angelbroking/order/v1/getOrderBook"
	URIPlaceOrder       string = "rest/secure/angelbroking/order/v1/placeOrder"
	URIModifyOrder      string = "rest/secure/angelbroking/order/v1/modifyOrder"
	URICancelOrder      string = "rest/secure/angelbroking/order/v1/cancelOrder"
	URIGetHoldings      string = "rest/secure/angelbroking/portfolio/v1/getHolding"
	URIGetPositions     string = "rest/secure/angelbroking/order/v1/getPosition"
	URIGetTradeBook     string = "rest/secure/angelbroking/order/v1/getTradeBook"
	URILTP              string = "rest/secure/angelbroking/order/v1/getLtpData"
	URIRMS              string = "rest/secure/angelbroking/user/v1/getRMS"
	URIConvertPosition  string = "rest/secure/angelbroking/order/v1/convertPosition"
)

// mapping function
func structToMap(obj interface{}, tagName string) map[string]interface{} {
	var values reflect.Value
	switch obj.(type) {
	case OrderParams:
		{
			con := obj.(OrderParams)
			values = reflect.ValueOf(&con).Elem()
		}

	case ModifyOrderParams:
		{
			con := obj.(ModifyOrderParams)
			values = reflect.ValueOf(&con).Elem()
		}
	case LTPParams:
		{
			con := obj.(LTPParams)
			values = reflect.ValueOf(&con).Elem()
		}
	case ConvertPositionParams:
		{
			con := obj.(ConvertPositionParams)
			values = reflect.ValueOf(&con).Elem()
		}
	}

	tags := reflect.TypeOf(obj)
	params := make(map[string]interface{})
	for i := 0; i < values.NumField(); i++ {
		params[tags.Field(i).Tag.Get(tagName)] = values.Field(i).Interface()
	}

	return params
}

// return ip and mac address -> local,public,mac,error
func getIpAndMac() (string, string, string, error) {

	var localIp, currentNetworkHardwareName string

	localIp, err := getLocalIP()

	if err != nil {
		return "", "", "", err
	}

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {

		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				if strings.Contains(addr.String(), localIp) {
					currentNetworkHardwareName = interf.Name
				}
			}
		}
	}

	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		return "", "", "", err
	}

	macAddress := netInterface.HardwareAddr

	_, err = net.ParseMAC(macAddress.String())

	if err != nil {
		return "", "", "", err
	}

	publicIp, err := getPublicIp()
	if err != nil {
		return "", "", "", err
	}

	return localIp, publicIp, macAddress.String(), nil

}

// return local ip and error -> local ip,error
func getLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("I think network is Down")
}

func getPublicIp() (string, error) {
	// extracted by other website
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return "", err
	}

	content, _ := ioutil.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(content), nil
}