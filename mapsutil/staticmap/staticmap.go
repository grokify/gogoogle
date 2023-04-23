// staticmap assists with https://developers.google.com/maps/documentation/maps-static/start
package staticmap

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/grokify/mogo/data/location"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/mime/mimeutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"google.golang.org/genproto/googleapis/type/latlng"
)

const (
	DefaultSize = "500x500"
	DefaultZoom = 5
	APIURL      = "https://maps.googleapis.com/maps/api/staticmap"

	ParamCenter  = "center"
	ParamKey     = "key"
	ParamMarkers = "markers"
	ParamSize    = "size"
	ParamZoom    = "zoom"
)

type StaticMap struct {
	Center          *latlng.LatLng
	Zoom            uint
	Height          uint
	Width           uint
	MarkersList     []Markers
	LatLngPrecision uint
}

func (sm *StaticMap) SetDefaultUSAKHI() {
	sm.Center = &location.USCenterAKHI
	sm.Height = 400
	sm.Width = 800
	sm.Zoom = 3
}

func (sm *StaticMap) SetDefaultUSContiguous() {
	sm.Center = &location.USCenterContiguous
	sm.Height = 400
	sm.Width = 800
	sm.Zoom = 4
}

func (sm *StaticMap) SetDefaultWorld() {
	sm.Center = &location.USHIEUCenter
	sm.Height = 400
	sm.Width = 800
	sm.Zoom = 1
}

func (sm *StaticMap) SetDefaultEU() {
	sm.Center = &location.EUCenter2020
	sm.Height = 600
	sm.Width = 800
	sm.Zoom = 4
}

func (sm *StaticMap) Size() string {
	if sm.Height > 0 && sm.Width > 0 {
		return strconv.Itoa(int(sm.Width)) + "x" + strconv.Itoa(int(sm.Height))
	}
	return DefaultSize
}

func (sm *StaticMap) URL(key string) string {
	params := url.Values{}
	params.Add(ParamCenter, location.LatLngString(sm.Center, ",", int(sm.LatLngPrecision)))
	params.Add(ParamSize, sm.Size())
	if sm.Zoom > 0 {
		params.Add(ParamZoom, strconv.Itoa(int(sm.Zoom)))
	} else {
		params.Add(ParamZoom, strconv.Itoa(DefaultZoom))
	}
	if len(strings.TrimSpace(key)) > 0 {
		params.Add(ParamKey, key)
	}
	for _, m := range sm.MarkersList {
		params.Add(ParamMarkers, m.String(sm.LatLngPrecision))
	}
	qry := params.Encode()
	u := APIURL + "?" + qry
	return u
}

func (sm *StaticMap) WriteFilePNG(filename, key string) error {
	u := sm.URL(key)
	r, err := http.Get(u) // #nosec G107
	if err != nil {
		return err
	}
	ct := r.Header.Get(httputilmore.HeaderContentType)
	if !mimeutil.IsType(ct, httputilmore.ContentTypeImagePNG) {
		return errorsutil.Wrapf(mimeutil.ErrUnknownMediaType, "mediaType [%s]", ct)
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	_, err = f.ReadFrom(r.Body)
	if err != nil {
		errClose := f.Close()
		if errClose != nil {
			return errorsutil.Wrapf(err, errClose.Error())
		}
		return err
	}
	return f.Close()
}
