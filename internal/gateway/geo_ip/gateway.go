// DOCS: https://ipwhois.io/documentation

package geo_ip

import (
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/geo_ip/methods/get_geo_by_ip"
)

type (
	getGeoByIpMethod struct{ *get_geo_by_ip.Method }
)

type Gateway struct {
	getGeoByIpMethod
}

func New(client gateway.HTTPGetter) *Gateway {
	return &Gateway{
		getGeoByIpMethod{get_geo_by_ip.New(client)},
	}
}
