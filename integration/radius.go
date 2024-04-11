package integration

import (
	"context"
	"fmt"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func (r *RadiusEndpointType) Authentication(username, password string) (err error) {
	packet := radius.New(radius.CodeAccessRequest, []byte(r.Secret))
	rfc2865.UserName_SetString(packet, username)
	rfc2865.UserPassword_SetString(packet, password)
	response, err := radius.Exchange(context.Background(), packet, fmt.Sprintf("%s:%s", r.Hostname, r.Secret))
	if err != nil {
		return
	} else {

		if response.Code != radius.CodeAccessAccept {
			err = fmt.Errorf("%s authentication failed", username)
			return
		}

	}

	return
}
