package get_geo_by_ip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway"
	"github.com/yasonofriychuk/school-21-workshop/internal/gateway/geo_ip/utils"
	"net/http"
	"net/url"
)

var (
	queryParams = map[string][]string{
		"output": {"json"},
		"lang":   {"ru"},
	}
)

type Method struct {
	queryTemplate string
	client        gateway.HTTPGetter
}

func New(httpClient gateway.HTTPGetter) *Method {
	return &Method{
		queryTemplate: fmt.Sprintf("%s/%%s?%s", utils.BaseURL, url.Values(queryParams).Encode()),
		client:        httpClient,
	}
}

func (m *Method) GetGeoByIP(_ context.Context, ipAdr string) (Coords, error) {
	resp, err := m.client.Get(fmt.Sprintf(m.queryTemplate, ipAdr))
	if err != nil {
		return Coords{}, fmt.Errorf("m.client.Get: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coords{}, fmt.Errorf("status code not 200: code: %d", resp.StatusCode)
	}

	var body Response
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return Coords{}, fmt.Errorf("json.NewDecoder: Decode: %w", err)
	}

	return Coords{
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
	}, nil
}
