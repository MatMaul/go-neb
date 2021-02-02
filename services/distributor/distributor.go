package distributor

import (
	"fmt"
	"strings"

	"github.com/matrix-org/go-neb/types"
	mevt "maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const ServiceType = "distributor"

type Service struct {
	types.DefaultService

	Items map[string]map[string]string `json:"items"`
}

// Commands supported:
//    !echo some message
// Responds with a notice of "some message".
func (e *Service) Commands(cli types.MatrixClient) []types.Command {
	var cmds []types.Command

	for itemType, items := range e.Items {
		cmds = append(cmds, types.Command{

			Path: []string{itemType},
			Command: func(roomID id.RoomID, userID id.UserID, args []string, eventID id.EventID) ([]interface{}, error) {
				itemName := strings.TrimSpace(strings.Join(args, " "))
				itemURL, exists := items[itemName]
				if exists {
					return []interface{}{
						&mevt.ReactionEventContent{
							RelatesTo: mevt.RelatesTo{
								Type:    mevt.RelAnnotation,
								EventID: id.EventID(eventID),
								Key:     fmt.Sprintf("<img data-mx-emoticon src=\"%s\">", itemURL),
							},
						},
						&mevt.MessageEventContent{
							MsgType: mevt.MsgImage,
							URL:     id.ContentURIString(itemURL),
						},
					}, nil
				} else {
					return []interface{}{
						&mevt.MessageEventContent{
							MsgType: mevt.MsgText,
							Body:    fmt.Sprintf("Sorry! This distributor doesn't have this %s unfortunately :(", itemType),
						},
					}, nil
				}
			},
		})
	}

	return cmds
}

func init() {
	types.RegisterService(func(serviceID string, serviceUserID id.UserID, webhookEndpointURL string) types.Service {
		return &Service{
			DefaultService: types.NewDefaultService(serviceID, serviceUserID, ServiceType),
		}
	})
}
