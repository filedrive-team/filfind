package ws

import (
	"context"
	"encoding/json"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/repo"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
)

// Note:
// Subscribed messages must be subscribed to receive

type Handler func(pkg *WrapPackageReq)

type ChannelHandler struct {
	Handler Handler
	Auth    bool // Authorized to access
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[string]*ClientSet

	// Unknown clients
	unknownClients *ClientSet

	// Inbound messages from the clients.
	broadcast chan *WrapPackageReq

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Registered handlers.
	handlers map[Channel]ChannelHandler

	repo *repo.Manager
}

func NewHub(repo *repo.Manager) *Hub {
	hub := &Hub{
		broadcast:      make(chan *WrapPackageReq),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[string]*ClientSet),
		unknownClients: NewClientSet(),
		handlers:       make(map[Channel]ChannelHandler),
		repo:           repo,
	}
	hub.RegisterHandler(MessageError, nil, false)
	hub.RegisterHandler(Echo, hub.echoHandler, false)
	hub.RegisterHandler(Login, hub.loginHandler, false)
	hub.RegisterHandler(Logout, hub.logoutHandler, true)
	hub.RegisterHandler(ChatSend, hub.sendHandler, true)
	hub.RegisterHandler(ChatReceive, nil, true)
	hub.RegisterHandler(ChatPartners, hub.partnersHandler, true)
	hub.RegisterHandler(ChatHistory, hub.historyHandler, true)
	hub.RegisterHandler(ChatPartnersStatus, hub.partnersStatusHandler, true)
	hub.RegisterHandler(ChatCheck, hub.checkHandler, true)
	hub.RegisterHandler(ChatCheckedStatus, nil, true)
	hub.RegisterHandler(ChatUncheckedTotal, hub.totalUncheckedHandler, true)
	hub.RegisterHandler(ChatUncheckedList, hub.uncheckedListHandler, true)

	return hub
}

func (h *Hub) RegisterHandler(channel Channel, handler Handler, auth bool) {
	h.handlers[channel] = ChannelHandler{
		Handler: handler,
		Auth:    auth,
	}
}

func (h *Hub) DeliverSystemMessage(recipient string, content string) error {
	systemUser, err := h.repo.QueryUserByEmail(settings.SystemUser)
	if err != nil {
		logger.WithError(err).Error("call QueryUserByEmail failed")
		return err
	}
	// check Recipient validity
	exist, err := h.repo.ExistUserByUid(recipient)
	if err != nil {
		logger.WithError(err).Error("call ExistUserByUid failed")
		return err
	} else if !exist {
		err = errors.New("not found the user")
		logger.WithField("Recipient", recipient).Error(err)
		return err
	}

	// If no relationship has been established between users, then establish a relationship
	has, err := h.repo.HasRelationByUser(recipient, systemUser.Uid.String())
	if err != nil {
		logger.WithError(err).Error("call HasRelationByUser failed")
		return err
	}
	if !has {
		if err = h.repo.CreateRelationship(systemUser.Uid.String(), recipient); err != nil {
			logger.WithError(err).Error("call CreateRelationship failed")
			return err
		}

		// notify them to update partners list each other
		h.deliverPartnersMsg(recipient)
	}

	// write message to database
	msgDb := &models.Message{
		Sender:    systemUser.Uid.String(),
		Recipient: recipient,
		Content:   content,
	}
	if err = h.repo.CreateMessage(msgDb); err != nil {
		logger.WithError(err).Error("call CreateMessage failed")
		return err
	}

	// notify the Recipient, update the number of unchecked messages
	h.deliverTotalUncheckedMsg(recipient)
	h.deliverUncheckedListMsg(recipient, systemUser.Uid.String())

	// deliver the message
	resMsg := &Message{
		Sender:    systemUser.Uid.String(),
		Recipient: recipient,
		Content:   content,
		Timestamp: &msgDb.CreatedAt,
	}
	h.sendTo(recipient, ChatReceive, resMsg)
	return nil
}

func (h *Hub) putUnknownClient(client *Client) {
	h.unknownClients.Add(client)
	logger.WithField("conns", h.unknownClients.Size()).
		Debug("unknown client connection")
}

func (h *Hub) registerClient(client *Client) {
	if set, ok := h.clients[client.uid]; ok {
		set.Add(client)
	} else {
		h.clients[client.uid] = NewClientSet(client)
		// notify its partners
		h.notifyPartnersStatus(client.uid, true)
	}
	logger.WithField("uid", client.uid).WithField("conns", h.clients[client.uid].Size()).
		Debug("accept client connection")
}

func (h *Hub) closeClient(client *Client) {
	if set, ok := h.clients[client.uid]; ok {
		set.Remove(client)
		if set.Size() == 0 {
			delete(h.clients, client.uid)
			// notify its partners
			h.notifyPartnersStatus(client.uid, false)
		}
		logger.WithField("uid", client.uid).WithField("conns", set.Size()).
			Debug("disconnect client connection")
	} else if h.unknownClients.Contains(client) {
		h.unknownClients.Remove(client)
		logger.WithField("unknown_conns", h.unknownClients.Size()).
			Debug("disconnect unknown client connection")
	} else {
		logger.Errorf("close client connection that havn't handle")
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.putUnknownClient(client)
		case client := <-h.unregister:
			h.closeClient(client)
		case pkg := <-h.broadcast:
			if handle, ok := h.handlers[pkg.Channel]; ok {
				if handle.Auth {
					// Require authorized access
					if pkg.Sender.uid != "" {
						if handle.Handler != nil {
							handle.Handler(pkg)
						}
					} else {
						h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrForbidden, pkg.Raw))
					}
				} else {
					// Not require authorized access
					if handle.Handler != nil {
						handle.Handler(pkg)
					}
				}
			} else {
				h.sendMessage(pkg.Sender, MessageError, newNotSupportMessageError(pkg.Raw))
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) sendMessage(client *Client, channel Channel, body interface{}) {
	pkgBytes, _ := encodePackageResp(channel, body)
	h.send(client, pkgBytes)
}

func (h *Hub) sendTo(uid string, channel Channel, body interface{}) {
	h.sendToExclude(uid, channel, body, nil)
}

func (h *Hub) sendToExclude(uid string, channel Channel, body interface{}, exclude *Client) {
	pkgBytes, _ := encodePackageResp(channel, body)
	if set, ok := h.clients[uid]; ok {
		set.Map(func(client *Client) {
			if client == exclude {
				return
			}
			if client.channelSet.Contains(channel) {
				h.send(client, pkgBytes)
			}
		})
	}
}

func (h *Hub) send(client *Client, pkgBytes []byte) {
	select {
	case client.send <- pkgBytes:
	default:
		// If too many messages are not sent, than close the connection.
		// Possible reason: the network may be unstable or attacked.
		h.closeClient(client)
	}
}

func (h *Hub) notifyPartnersStatus(clientUid string, online bool) {
	uid, err := uuid.FromString(clientUid)
	if err != nil {
		logger.WithError(err).WithField("clientUid", clientUid).Error("call uuid.FromString failed")
		return
	}
	err = h.repo.UpsertUserOnline(&models.UserStatus{
		Uid:    uid,
		Online: online,
	})
	if err != nil {
		logger.WithError(err).Error("call UpsertUserOnline failed")
		return
	}
	// notify its partners
	partners, err := h.repo.PartnersByUid(clientUid)
	if err != nil {
		logger.WithError(err).Error("call PartnersByUid failed")
		return
	}
	for _, partner := range partners {
		h.sendTo(partner.Uid, ChatPartnersStatus, []*repo.PartnerStatus{
			{
				Uid:    clientUid,
				Online: online,
			},
		})
	}
}

func (h *Hub) echoHandler(pkg *WrapPackageReq) {
	pkgBytes, _ := encodePackageResp(pkg.Channel, string(pkg.Raw))
	h.send(pkg.Sender, pkgBytes)
}

func (h *Hub) loginHandler(pkg *WrapPackageReq) {
	login := &LoginReq{}
	if err := json.Unmarshal(pkg.Body, login); err != nil {
		logger.WithError(err).WithField("package", pkg).Error("call json.Unmarshal to Login failed")
		h.sendMessage(pkg.Sender, MessageError, newNotSupportMessageError(pkg.Raw))
		return
	}
	uid, err := CheckAuthority(login.Token, h.repo.GetTokenVerify())
	if err != nil {
		h.sendMessage(pkg.Sender, MessageError, newMessageError(err.Error(), pkg.Raw))
		return
	}
	if h.unknownClients.Contains(pkg.Sender) {
		h.unknownClients.Remove(pkg.Sender)
		pkg.Sender.uid = uid
		h.registerClient(pkg.Sender)
		// deliver the message
		h.sendMessage(pkg.Sender, Login, &LoginResp{Msg: "ok"})
	} else if cliSet, ok := h.clients[uid]; ok && cliSet.Size() > 0 {
		// deliver the message
		h.sendMessage(pkg.Sender, Login, &LoginResp{Msg: "ok"})
	} else {
		h.sendMessage(pkg.Sender, MessageError, newMessageError("The previous user has logged in", pkg.Raw))
	}
}

func (h *Hub) logoutHandler(pkg *WrapPackageReq) {
	client := pkg.Sender
	if set, ok := h.clients[client.uid]; ok {
		set.Remove(client)
		if set.Size() == 0 {
			delete(h.clients, client.uid)
			// notify its partners
			h.notifyPartnersStatus(client.uid, false)
		}
		logger.WithField("uid", client.uid).WithField("conns", set.Size()).
			Debug("disconnect client connection")
	}
	// clear login status and add to unregistered clients
	client.uid = ""
	h.putUnknownClient(client)
}

func (h *Hub) sendHandler(pkg *WrapPackageReq) {
	msg := &Message{}
	if err := json.Unmarshal(pkg.Body, msg); err != nil {
		logger.WithError(err).WithField("package", pkg).Error("call json.Unmarshal to Message failed")
		h.sendMessage(pkg.Sender, MessageError, newNotSupportMessageError(pkg.Raw))
		return
	}
	// check Recipient validity
	exist, err := h.repo.ExistUserByUid(msg.Recipient)
	if err != nil {
		logger.WithError(err).Error("call ExistUserByUid failed")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrInternalServer, pkg.Raw))
		return
	} else if !exist {
		logger.WithField("Recipient", msg.Recipient).Error("not found user")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrUserNotFound, pkg.Raw))
		return
	}

	// is empty message?
	if msg.Content == "" {
		logger.Error("content of message is empty")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrMessageEmpty, pkg.Raw))
		return
	}

	// If no relationship has been established between users, then establish a relationship
	has, err := h.repo.HasRelationByUser(msg.Recipient, pkg.Sender.uid)
	if err != nil {
		logger.WithError(err).Error("call HasRelationByUser failed")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrInternalServer, pkg.Raw))
		return
	}
	if !has {
		if err = h.repo.CreateRelationship(pkg.Sender.uid, msg.Recipient); err != nil {
			logger.WithError(err).Error("call CreateRelationship failed")
			h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrInternalServer, pkg.Raw))
			return
		}

		// notify them to update partners list each other
		h.deliverPartnersMsg(msg.Sender)
		h.deliverPartnersMsg(msg.Recipient)
	}

	// write message to database
	msgDb := &models.Message{
		Sender:    pkg.Sender.uid,
		Recipient: msg.Recipient,
		Content:   msg.Content,
	}
	if err = h.repo.CreateMessage(msgDb); err != nil {
		logger.WithError(err).Error("call CreateMessage failed")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrInternalServer, pkg.Raw))
		return
	}

	// notify the Recipient, update the number of unchecked messages
	h.deliverTotalUncheckedMsg(msg.Recipient)
	h.deliverUncheckedListMsg(msg.Recipient, pkg.Sender.uid)

	// deliver the message
	resMsg := &Message{
		Sender:    pkg.Sender.uid,
		Recipient: msg.Recipient,
		Content:   msg.Content,
		Timestamp: &msgDb.CreatedAt,
	}
	h.sendTo(msg.Recipient, ChatReceive, resMsg)
	h.sendTo(pkg.Sender.uid, ChatReceive, resMsg)
}

func (h *Hub) partnersHandler(pkg *WrapPackageReq) {
	h.deliverPartnersMsg(pkg.Sender.uid)
}

func (h *Hub) deliverPartnersMsg(uid string) {
	partners, err := h.repo.PartnersByUid(uid)
	if err != nil {
		logger.WithError(err).Error("call PartnersByUid failed")
		return
	}
	// deliver the message
	h.sendTo(uid, ChatPartners, partners)
}

func (h *Hub) historyHandler(pkg *WrapPackageReq) {
	hmReq := &HistoryMessageReq{}
	if err := json.Unmarshal(pkg.Body, hmReq); err != nil {
		logger.WithError(err).WithField("package", pkg).Error("call json.Unmarshal to HistoryMessageReq failed")
		h.sendMessage(pkg.Sender, MessageError, newNotSupportMessageError(pkg.Raw))
		return
	}
	if hmReq.Limit < 1 || hmReq.Limit > 100 {
		err := errors.New("limit param is out of limit[1,100]")
		logger.WithField("HistoryMessageReq", hmReq).Error(err)
		h.sendMessage(pkg.Sender, MessageError, newMessageError(err.Error(), pkg.Raw))
		return
	}

	list, err := h.repo.MessagesByUser(pkg.Sender.uid, hmReq.Partner, hmReq.Before, hmReq.Limit)
	if err != nil {
		logger.WithError(err).Error("call MessagesByUser failed")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrInternalServer, pkg.Raw))
		return
	}
	// deliver the message
	results := make([]*Message, 0, len(list))
	for _, msg := range list {
		results = append(results, &Message{
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
			Type:      msg.Type,
			Content:   msg.Content,
			Timestamp: &msg.CreatedAt,
			Checked:   msg.Checked,
		})
	}
	h.sendMessage(pkg.Sender, ChatHistory, results)
}

func (h *Hub) partnersStatusHandler(pkg *WrapPackageReq) {
	results, err := h.repo.PartnersStatusByUid(pkg.Sender.uid)
	if err != nil {
		logger.WithError(err).Error("call PartnersStatusByUid failed")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrInternalServer, pkg.Raw))
		return
	}
	// deliver the message
	h.sendMessage(pkg.Sender, ChatPartnersStatus, results)
}

func (h *Hub) checkHandler(pkg *WrapPackageReq) {
	cmReq := &CheckMessageReq{}
	if err := json.Unmarshal(pkg.Body, cmReq); err != nil {
		logger.WithError(err).WithField("package", pkg).Error("call json.Unmarshal to CheckMessageReq failed")
		h.sendMessage(pkg.Sender, MessageError, newNotSupportMessageError(pkg.Raw))
		return
	}
	// checked message from the partner
	err := h.repo.UpdateMessagesChecked(cmReq.Partner, pkg.Sender.uid)
	if err != nil {
		logger.WithError(err).Error("call UpdateMessagesChecked failed")
		h.sendMessage(pkg.Sender, MessageError, newMessageError(ErrInternalServer, pkg.Raw))
		return
	}

	// notify the sender, update the number of unchecked messages
	h.deliverTotalUncheckedMsg(pkg.Sender.uid)
	h.deliverZeroUncheckedListMsg(pkg.Sender.uid, cmReq.Partner)

	// notify the partner, checked the message
	resMsg := &CheckedMessage{
		Partner: pkg.Sender.uid,
	}
	h.sendTo(cmReq.Partner, ChatCheckedStatus, resMsg)
}

func (h *Hub) totalUncheckedHandler(pkg *WrapPackageReq) {
	h.deliverTotalUncheckedMsg(pkg.Sender.uid)
}

func (h *Hub) deliverTotalUncheckedMsg(recipientUid string) {
	total, err := h.repo.UncheckedMessagesByRecipient(recipientUid)
	if err != nil {
		logger.WithError(err).Error("call UncheckedMessagesByRecipient failed")
		return
	}
	// deliver the message
	h.sendTo(recipientUid, ChatUncheckedTotal, &TotalUncheckedMessage{Total: total})
}

func (h *Hub) uncheckedListHandler(pkg *WrapPackageReq) {
	list, err := h.repo.UncheckedMessageGroupListByRecipient(pkg.Sender.uid)
	if err != nil {
		logger.WithError(err).Error("call UncheckedMessageGroupListByRecipient failed")
		return
	}
	if len(list) == 0 {
		return
	}
	// deliver the message
	h.sendTo(pkg.Sender.uid, ChatUncheckedList, list)
}

func (h *Hub) deliverUncheckedListMsg(recipientUid, senderUid string) {
	list, err := h.repo.UncheckedMessagesByPair(recipientUid, senderUid)
	if err != nil {
		logger.WithError(err).Error("call UncheckedMessagesByPair failed")
		return
	}
	if len(list) == 0 {
		return
	}
	// deliver the message
	h.sendTo(recipientUid, ChatUncheckedList, list)
}

func (h *Hub) deliverZeroUncheckedListMsg(senderUid, partner string) {
	list := []*repo.UncheckedMessageItem{
		{
			Partner: partner,
			Number:  0,
		},
	}
	// deliver the message
	h.sendTo(senderUid, ChatUncheckedList, list)
}
