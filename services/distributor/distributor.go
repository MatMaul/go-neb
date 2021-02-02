package distributor

import (
	"fmt"
	"strings"

	"github.com/matrix-org/go-neb/types"
	mevt "maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const ServiceType = "distributor"

type item struct {
	Name      string `json:"name"`
	MatrixURL string `json:"matrix_url"`
	Reaction  string `json:"reaction"`
}

type Service struct {
	types.DefaultService

	Items map[string][]item `json:"items"`
}

// Commands supported:
//    !echo some message
// Responds with a notice of "some message".
func (e *Service) Commands(cli types.MatrixClient) []types.Command {
	var cmds []types.Command

	for itemType, items := range e.Items {
		itemsMap := make(map[string]item)
		for _, item := range items {
			itemsMap[item.Name] = item
		}

		cmds = append(cmds, types.Command{

			Path: []string{itemType},
			Command: func(roomID id.RoomID, userID id.UserID, args []string, eventID id.EventID) ([]interface{}, error) {
				itemName := strings.TrimSpace(strings.Join(args, " "))
				item, exists := itemsMap[itemName]
				if exists {
					return []interface{}{
						&mevt.ReactionEventContent{
							RelatesTo: mevt.RelatesTo{
								Type:    mevt.RelAnnotation,
								EventID: id.EventID(eventID),
								Key:     item.Reaction,
							},
						},
						&mevt.MessageEventContent{
							MsgType: mevt.MsgImage,
							URL:     id.ContentURIString(item.MatrixURL),
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
