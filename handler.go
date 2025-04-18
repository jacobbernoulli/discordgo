package discordgo

const (
	applicationCommandPermissionsUpdateEventType = "APPLICATION_COMMAND_PERMISSIONS_UPDATE"
	autoModerationActionExecutionEventType       = "AUTO_MODERATION_ACTION_EXECUTION"
	autoModerationRuleCreateEventType            = "AUTO_MODERATION_RULE_CREATE"
	autoModerationRuleDeleteEventType            = "AUTO_MODERATION_RULE_DELETE"
	autoModerationRuleUpdateEventType            = "AUTO_MODERATION_RULE_UPDATE"
	channelCreateEventType                       = "CHANNEL_CREATE"
	channelDeleteEventType                       = "CHANNEL_DELETE"
	channelPinsUpdateEventType                   = "CHANNEL_PINS_UPDATE"
	channelUpdateEventType                       = "CHANNEL_UPDATE"
	connectEventType                             = "__CONNECT__"
	disconnectEventType                          = "__DISCONNECT__"
	entitlementCreateEventType                   = "ENTITLEMENT_CREATE"
	entitlementDeleteEventType                   = "ENTITLEMENT_DELETE"
	entitlementUpdateEventType                   = "ENTITLEMENT_UPDATE"
	eventEventType                               = "__EVENT__"
	guildAuditLogEntryCreateEventType            = "GUILD_AUDIT_LOG_ENTRY_CREATE"
	guildBanAddEventType                         = "GUILD_BAN_ADD"
	guildBanRemoveEventType                      = "GUILD_BAN_REMOVE"
	guildCreateEventType                         = "GUILD_CREATE"
	guildDeleteEventType                         = "GUILD_DELETE"
	guildEmojisUpdateEventType                   = "GUILD_EMOJIS_UPDATE"
	guildIntegrationsUpdateEventType             = "GUILD_INTEGRATIONS_UPDATE"
	guildMemberAddEventType                      = "GUILD_MEMBER_ADD"
	guildMemberRemoveEventType                   = "GUILD_MEMBER_REMOVE"
	guildMemberUpdateEventType                   = "GUILD_MEMBER_UPDATE"
	guildMembersChunkEventType                   = "GUILD_MEMBERS_CHUNK"
	guildRoleCreateEventType                     = "GUILD_ROLE_CREATE"
	guildRoleDeleteEventType                     = "GUILD_ROLE_DELETE"
	guildRoleUpdateEventType                     = "GUILD_ROLE_UPDATE"
	guildScheduledEventCreateEventType           = "GUILD_SCHEDULED_EVENT_CREATE"
	guildScheduledEventDeleteEventType           = "GUILD_SCHEDULED_EVENT_DELETE"
	guildScheduledEventUpdateEventType           = "GUILD_SCHEDULED_EVENT_UPDATE"
	guildScheduledEventUserAddEventType          = "GUILD_SCHEDULED_EVENT_USER_ADD"
	guildScheduledEventUserRemoveEventType       = "GUILD_SCHEDULED_EVENT_USER_REMOVE"
	guildUpdateEventType                         = "GUILD_UPDATE"
	integrationCreateEventType                   = "INTEGRATION_CREATE"
	integrationDeleteEventType                   = "INTEGRATION_DELETE"
	integrationUpdateEventType                   = "INTEGRATION_UPDATE"
	interactionCreateEventType                   = "INTERACTION_CREATE"
	inviteCreateEventType                        = "INVITE_CREATE"
	inviteDeleteEventType                        = "INVITE_DELETE"
	messageCreateEventType                       = "MESSAGE_CREATE"
	messageDeleteEventType                       = "MESSAGE_DELETE"
	messageDeleteBulkEventType                   = "MESSAGE_DELETE_BULK"
	messagePollVoteAddEventType                  = "MESSAGE_POLL_VOTE_ADD"
	messagePollVoteRemoveEventType               = "MESSAGE_POLL_VOTE_REMOVE"
	messageReactionAddEventType                  = "MESSAGE_REACTION_ADD"
	messageReactionRemoveEventType               = "MESSAGE_REACTION_REMOVE"
	messageReactionRemoveAllEventType            = "MESSAGE_REACTION_REMOVE_ALL"
	messageUpdateEventType                       = "MESSAGE_UPDATE"
	presenceUpdateEventType                      = "PRESENCE_UPDATE"
	presencesReplaceEventType                    = "PRESENCES_REPLACE"
	rateLimitEventType                           = "__RATE_LIMIT__"
	readyEventType                               = "READY"
	resumedEventType                             = "RESUMED"
	stageInstanceEventCreateEventType            = "STAGE_INSTANCE_EVENT_CREATE"
	stageInstanceEventDeleteEventType            = "STAGE_INSTANCE_EVENT_DELETE"
	stageInstanceEventUpdateEventType            = "STAGE_INSTANCE_EVENT_UPDATE"
	subscriptionCreateEventType                  = "SUBSCRIPTION_CREATE"
	subscriptionDeleteEventType                  = "SUBSCRIPTION_DELETE"
	subscriptionUpdateEventType                  = "SUBSCRIPTION_UPDATE"
	threadCreateEventType                        = "THREAD_CREATE"
	threadDeleteEventType                        = "THREAD_DELETE"
	threadListSyncEventType                      = "THREAD_LIST_SYNC"
	threadMemberUpdateEventType                  = "THREAD_MEMBER_UPDATE"
	threadMembersUpdateEventType                 = "THREAD_MEMBERS_UPDATE"
	threadUpdateEventType                        = "THREAD_UPDATE"
	typingStartEventType                         = "TYPING_START"
	userUpdateEventType                          = "USER_UPDATE"
	voiceServerUpdateEventType                   = "VOICE_SERVER_UPDATE"
	voiceStateUpdateEventType                    = "VOICE_STATE_UPDATE"
	webhooksUpdateEventType                      = "WEBHOOKS_UPDATE"
)

type EventHandler interface {
	Type() string
	Handle(*Session, interface{})
}

type EventInterfaceProvider interface {
	Type() string
	New() interface{}
}

const interfaceEventType = "__INTERFACE__"

type interfaceEventHandler func(*Session, interface{})

func (eh interfaceEventHandler) Type() string {
	return interfaceEventType
}

func (eh interfaceEventHandler) Handle(s *Session, i interface{}) {
	eh(s, i)
}

var registeredInterfaceProviders = map[string]EventInterfaceProvider{}

func registerInterfaceProvider(eh EventInterfaceProvider) {
	if _, ok := registeredInterfaceProviders[eh.Type()]; ok {
		return
	}
	registeredInterfaceProviders[eh.Type()] = eh
}

type eventHandlerInstance struct {
	eventHandler EventHandler
}

func (s *Session) addEventHandler(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.handlers == nil {
		s.handlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.handlers[eventHandler.Type()] = append(s.handlers[eventHandler.Type()], ehi)

	return func() {
		s.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

func (s *Session) addEventHandlerOnce(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.onceHandlers == nil {
		s.onceHandlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.onceHandlers[eventHandler.Type()] = append(s.onceHandlers[eventHandler.Type()], ehi)

	return func() {
		s.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

func (s *Session) AddHandler(handler interface{}) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		s.log(LogError, "Invalid handler type, handler will never be called")
		return func() {}
	}

	return s.addEventHandler(eh)
}

func (s *Session) AddHandlerOnce(handler interface{}) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		s.log(LogError, "Invalid handler type, handler will never be called")
		return func() {}
	}

	return s.addEventHandlerOnce(eh)
}

func (s *Session) removeEventHandlerInstance(t string, ehi *eventHandlerInstance) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	handlers := s.handlers[t]
	for i := range handlers {
		if handlers[i] == ehi {
			s.handlers[t] = append(handlers[:i], handlers[i+1:]...)
		}
	}

	onceHandlers := s.onceHandlers[t]
	for i := range onceHandlers {
		if onceHandlers[i] == ehi {
			s.onceHandlers[t] = append(onceHandlers[:i], onceHandlers[i+1:]...)
		}
	}
}

func (s *Session) handle(t string, i interface{}) {
	for _, eh := range s.handlers[t] {
		if s.SyncEvents {
			eh.eventHandler.Handle(s, i)
		} else {
			go eh.eventHandler.Handle(s, i)
		}
	}

	if len(s.onceHandlers[t]) > 0 {
		for _, eh := range s.onceHandlers[t] {
			if s.SyncEvents {
				eh.eventHandler.Handle(s, i)
			} else {
				go eh.eventHandler.Handle(s, i)
			}
		}
		s.onceHandlers[t] = nil
	}
}

func (s *Session) handleEvent(t string, i interface{}) {
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()
	s.onInterface(i)
	s.handle(interfaceEventType, i)
	s.handle(t, i)
}

func setGuildIds(g *Guild) {
	for _, c := range g.Channels {
		c.GuildID = g.ID
	}

	for _, m := range g.Members {
		m.GuildID = g.ID
	}

	for _, vs := range g.VoiceStates {
		vs.GuildID = g.ID
	}
}

func (s *Session) onInterface(i interface{}) {
	switch t := i.(type) {
	case *Ready:
		for _, g := range t.Guilds {
			setGuildIds(g)
		}
		s.onReady(t)
	case *GuildCreate:
		setGuildIds(t.Guild)
	case *GuildUpdate:
		setGuildIds(t.Guild)
	case *VoiceServerUpdate:
		go s.onVoiceServerUpdate(t)
	case *VoiceStateUpdate:
		go s.onVoiceStateUpdate(t)
	}
	err := s.State.OnInterface(s, i)
	if err != nil {
		s.log(LogDebug, "error dispatching internal event, %s", err)
	}
}

func (s *Session) onReady(r *Ready) {
	s.sessionID = r.SessionID
}

type applicationCommandPermissionsUpdateEventHandler func(*Session, *ApplicationCommandPermissionsUpdate)

func (eh applicationCommandPermissionsUpdateEventHandler) Type() string {
	return applicationCommandPermissionsUpdateEventType
}

func (eh applicationCommandPermissionsUpdateEventHandler) New() interface{} {
	return &ApplicationCommandPermissionsUpdate{}
}

func (eh applicationCommandPermissionsUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ApplicationCommandPermissionsUpdate); ok {
		eh(s, t)
	}
}

type autoModerationActionExecutionEventHandler func(*Session, *AutoModerationActionExecution)

func (eh autoModerationActionExecutionEventHandler) Type() string {
	return autoModerationActionExecutionEventType
}

func (eh autoModerationActionExecutionEventHandler) New() interface{} {
	return &AutoModerationActionExecution{}
}

func (eh autoModerationActionExecutionEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*AutoModerationActionExecution); ok {
		eh(s, t)
	}
}

type autoModerationRuleCreateEventHandler func(*Session, *AutoModerationRuleCreate)

func (eh autoModerationRuleCreateEventHandler) Type() string {
	return autoModerationRuleCreateEventType
}

func (eh autoModerationRuleCreateEventHandler) New() interface{} {
	return &AutoModerationRuleCreate{}
}

func (eh autoModerationRuleCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*AutoModerationRuleCreate); ok {
		eh(s, t)
	}
}

type autoModerationRuleDeleteEventHandler func(*Session, *AutoModerationRuleDelete)

func (eh autoModerationRuleDeleteEventHandler) Type() string {
	return autoModerationRuleDeleteEventType
}

func (eh autoModerationRuleDeleteEventHandler) New() interface{} {
	return &AutoModerationRuleDelete{}
}

func (eh autoModerationRuleDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*AutoModerationRuleDelete); ok {
		eh(s, t)
	}
}

type autoModerationRuleUpdateEventHandler func(*Session, *AutoModerationRuleUpdate)

func (eh autoModerationRuleUpdateEventHandler) Type() string {
	return autoModerationRuleUpdateEventType
}

func (eh autoModerationRuleUpdateEventHandler) New() interface{} {
	return &AutoModerationRuleUpdate{}
}

func (eh autoModerationRuleUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*AutoModerationRuleUpdate); ok {
		eh(s, t)
	}
}

type channelCreateEventHandler func(*Session, *ChannelCreate)

func (eh channelCreateEventHandler) Type() string {
	return channelCreateEventType
}

func (eh channelCreateEventHandler) New() interface{} {
	return &ChannelCreate{}
}

func (eh channelCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelCreate); ok {
		eh(s, t)
	}
}

type channelDeleteEventHandler func(*Session, *ChannelDelete)

func (eh channelDeleteEventHandler) Type() string {
	return channelDeleteEventType
}

func (eh channelDeleteEventHandler) New() interface{} {
	return &ChannelDelete{}
}

func (eh channelDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelDelete); ok {
		eh(s, t)
	}
}

type channelPinsUpdateEventHandler func(*Session, *ChannelPinsUpdate)

func (eh channelPinsUpdateEventHandler) Type() string {
	return channelPinsUpdateEventType
}

func (eh channelPinsUpdateEventHandler) New() interface{} {
	return &ChannelPinsUpdate{}
}

func (eh channelPinsUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelPinsUpdate); ok {
		eh(s, t)
	}
}

type channelUpdateEventHandler func(*Session, *ChannelUpdate)

func (eh channelUpdateEventHandler) Type() string {
	return channelUpdateEventType
}

func (eh channelUpdateEventHandler) New() interface{} {
	return &ChannelUpdate{}
}

func (eh channelUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelUpdate); ok {
		eh(s, t)
	}
}

type connectEventHandler func(*Session, *Connect)

func (eh connectEventHandler) Type() string {
	return connectEventType
}

func (eh connectEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Connect); ok {
		eh(s, t)
	}
}

type disconnectEventHandler func(*Session, *Disconnect)

func (eh disconnectEventHandler) Type() string {
	return disconnectEventType
}

func (eh disconnectEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Disconnect); ok {
		eh(s, t)
	}
}

type entitlementCreateEventHandler func(*Session, *EntitlementCreate)

func (eh entitlementCreateEventHandler) Type() string {
	return entitlementCreateEventType
}

func (eh entitlementCreateEventHandler) New() interface{} {
	return &EntitlementCreate{}
}

func (eh entitlementCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*EntitlementCreate); ok {
		eh(s, t)
	}
}

type entitlementDeleteEventHandler func(*Session, *EntitlementDelete)

func (eh entitlementDeleteEventHandler) Type() string {
	return entitlementDeleteEventType
}

func (eh entitlementDeleteEventHandler) New() interface{} {
	return &EntitlementDelete{}
}

func (eh entitlementDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*EntitlementDelete); ok {
		eh(s, t)
	}
}

type entitlementUpdateEventHandler func(*Session, *EntitlementUpdate)

func (eh entitlementUpdateEventHandler) Type() string {
	return entitlementUpdateEventType
}

func (eh entitlementUpdateEventHandler) New() interface{} {
	return &EntitlementUpdate{}
}

func (eh entitlementUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*EntitlementUpdate); ok {
		eh(s, t)
	}
}

type eventEventHandler func(*Session, *Event)

func (eh eventEventHandler) Type() string {
	return eventEventType
}

func (eh eventEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Event); ok {
		eh(s, t)
	}
}

type guildAuditLogEntryCreateEventHandler func(*Session, *GuildAuditLogEntryCreate)

func (eh guildAuditLogEntryCreateEventHandler) Type() string {
	return guildAuditLogEntryCreateEventType
}

func (eh guildAuditLogEntryCreateEventHandler) New() interface{} {
	return &GuildAuditLogEntryCreate{}
}

func (eh guildAuditLogEntryCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildAuditLogEntryCreate); ok {
		eh(s, t)
	}
}

type guildBanAddEventHandler func(*Session, *GuildBanAdd)

func (eh guildBanAddEventHandler) Type() string {
	return guildBanAddEventType
}

func (eh guildBanAddEventHandler) New() interface{} {
	return &GuildBanAdd{}
}

func (eh guildBanAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildBanAdd); ok {
		eh(s, t)
	}
}

type guildBanRemoveEventHandler func(*Session, *GuildBanRemove)

func (eh guildBanRemoveEventHandler) Type() string {
	return guildBanRemoveEventType
}

func (eh guildBanRemoveEventHandler) New() interface{} {
	return &GuildBanRemove{}
}

func (eh guildBanRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildBanRemove); ok {
		eh(s, t)
	}
}

type guildCreateEventHandler func(*Session, *GuildCreate)

func (eh guildCreateEventHandler) Type() string {
	return guildCreateEventType
}

func (eh guildCreateEventHandler) New() interface{} {
	return &GuildCreate{}
}

func (eh guildCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildCreate); ok {
		eh(s, t)
	}
}

type guildDeleteEventHandler func(*Session, *GuildDelete)

func (eh guildDeleteEventHandler) Type() string {
	return guildDeleteEventType
}

func (eh guildDeleteEventHandler) New() interface{} {
	return &GuildDelete{}
}

func (eh guildDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildDelete); ok {
		eh(s, t)
	}
}

type guildEmojisUpdateEventHandler func(*Session, *GuildEmojisUpdate)

func (eh guildEmojisUpdateEventHandler) Type() string {
	return guildEmojisUpdateEventType
}

func (eh guildEmojisUpdateEventHandler) New() interface{} {
	return &GuildEmojisUpdate{}
}

func (eh guildEmojisUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildEmojisUpdate); ok {
		eh(s, t)
	}
}

type guildIntegrationsUpdateEventHandler func(*Session, *GuildIntegrationsUpdate)

func (eh guildIntegrationsUpdateEventHandler) Type() string {
	return guildIntegrationsUpdateEventType
}

func (eh guildIntegrationsUpdateEventHandler) New() interface{} {
	return &GuildIntegrationsUpdate{}
}

func (eh guildIntegrationsUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildIntegrationsUpdate); ok {
		eh(s, t)
	}
}

type guildMemberAddEventHandler func(*Session, *GuildMemberAdd)

func (eh guildMemberAddEventHandler) Type() string {
	return guildMemberAddEventType
}

func (eh guildMemberAddEventHandler) New() interface{} {
	return &GuildMemberAdd{}
}

func (eh guildMemberAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMemberAdd); ok {
		eh(s, t)
	}
}

type guildMemberRemoveEventHandler func(*Session, *GuildMemberRemove)

func (eh guildMemberRemoveEventHandler) Type() string {
	return guildMemberRemoveEventType
}

func (eh guildMemberRemoveEventHandler) New() interface{} {
	return &GuildMemberRemove{}
}

func (eh guildMemberRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMemberRemove); ok {
		eh(s, t)
	}
}

type guildMemberUpdateEventHandler func(*Session, *GuildMemberUpdate)

func (eh guildMemberUpdateEventHandler) Type() string {
	return guildMemberUpdateEventType
}

func (eh guildMemberUpdateEventHandler) New() interface{} {
	return &GuildMemberUpdate{}
}

func (eh guildMemberUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMemberUpdate); ok {
		eh(s, t)
	}
}

type guildMembersChunkEventHandler func(*Session, *GuildMembersChunk)

func (eh guildMembersChunkEventHandler) Type() string {
	return guildMembersChunkEventType
}

func (eh guildMembersChunkEventHandler) New() interface{} {
	return &GuildMembersChunk{}
}

func (eh guildMembersChunkEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMembersChunk); ok {
		eh(s, t)
	}
}

type guildRoleCreateEventHandler func(*Session, *GuildRoleCreate)

func (eh guildRoleCreateEventHandler) Type() string {
	return guildRoleCreateEventType
}

func (eh guildRoleCreateEventHandler) New() interface{} {
	return &GuildRoleCreate{}
}

func (eh guildRoleCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildRoleCreate); ok {
		eh(s, t)
	}
}

type guildRoleDeleteEventHandler func(*Session, *GuildRoleDelete)

func (eh guildRoleDeleteEventHandler) Type() string {
	return guildRoleDeleteEventType
}

func (eh guildRoleDeleteEventHandler) New() interface{} {
	return &GuildRoleDelete{}
}

func (eh guildRoleDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildRoleDelete); ok {
		eh(s, t)
	}
}

type guildRoleUpdateEventHandler func(*Session, *GuildRoleUpdate)

func (eh guildRoleUpdateEventHandler) Type() string {
	return guildRoleUpdateEventType
}

func (eh guildRoleUpdateEventHandler) New() interface{} {
	return &GuildRoleUpdate{}
}

func (eh guildRoleUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildRoleUpdate); ok {
		eh(s, t)
	}
}

type guildScheduledEventCreateEventHandler func(*Session, *GuildScheduledEventCreate)

func (eh guildScheduledEventCreateEventHandler) Type() string {
	return guildScheduledEventCreateEventType
}

func (eh guildScheduledEventCreateEventHandler) New() interface{} {
	return &GuildScheduledEventCreate{}
}

func (eh guildScheduledEventCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildScheduledEventCreate); ok {
		eh(s, t)
	}
}

type guildScheduledEventDeleteEventHandler func(*Session, *GuildScheduledEventDelete)

func (eh guildScheduledEventDeleteEventHandler) Type() string {
	return guildScheduledEventDeleteEventType
}

func (eh guildScheduledEventDeleteEventHandler) New() interface{} {
	return &GuildScheduledEventDelete{}
}

func (eh guildScheduledEventDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildScheduledEventDelete); ok {
		eh(s, t)
	}
}

type guildScheduledEventUpdateEventHandler func(*Session, *GuildScheduledEventUpdate)

func (eh guildScheduledEventUpdateEventHandler) Type() string {
	return guildScheduledEventUpdateEventType
}

func (eh guildScheduledEventUpdateEventHandler) New() interface{} {
	return &GuildScheduledEventUpdate{}
}

func (eh guildScheduledEventUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildScheduledEventUpdate); ok {
		eh(s, t)
	}
}

type guildScheduledEventUserAddEventHandler func(*Session, *GuildScheduledEventUserAdd)

func (eh guildScheduledEventUserAddEventHandler) Type() string {
	return guildScheduledEventUserAddEventType
}

func (eh guildScheduledEventUserAddEventHandler) New() interface{} {
	return &GuildScheduledEventUserAdd{}
}

func (eh guildScheduledEventUserAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildScheduledEventUserAdd); ok {
		eh(s, t)
	}
}

type guildScheduledEventUserRemoveEventHandler func(*Session, *GuildScheduledEventUserRemove)

func (eh guildScheduledEventUserRemoveEventHandler) Type() string {
	return guildScheduledEventUserRemoveEventType
}

func (eh guildScheduledEventUserRemoveEventHandler) New() interface{} {
	return &GuildScheduledEventUserRemove{}
}

func (eh guildScheduledEventUserRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildScheduledEventUserRemove); ok {
		eh(s, t)
	}
}

type guildUpdateEventHandler func(*Session, *GuildUpdate)

func (eh guildUpdateEventHandler) Type() string {
	return guildUpdateEventType
}

func (eh guildUpdateEventHandler) New() interface{} {
	return &GuildUpdate{}
}

func (eh guildUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildUpdate); ok {
		eh(s, t)
	}
}

type integrationCreateEventHandler func(*Session, *IntegrationCreate)

func (eh integrationCreateEventHandler) Type() string {
	return integrationCreateEventType
}

func (eh integrationCreateEventHandler) New() interface{} {
	return &IntegrationCreate{}
}

func (eh integrationCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*IntegrationCreate); ok {
		eh(s, t)
	}
}

type integrationDeleteEventHandler func(*Session, *IntegrationDelete)

func (eh integrationDeleteEventHandler) Type() string {
	return integrationDeleteEventType
}

func (eh integrationDeleteEventHandler) New() interface{} {
	return &IntegrationDelete{}
}

func (eh integrationDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*IntegrationDelete); ok {
		eh(s, t)
	}
}

type integrationUpdateEventHandler func(*Session, *IntegrationUpdate)

func (eh integrationUpdateEventHandler) Type() string {
	return integrationUpdateEventType
}

func (eh integrationUpdateEventHandler) New() interface{} {
	return &IntegrationUpdate{}
}

func (eh integrationUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*IntegrationUpdate); ok {
		eh(s, t)
	}
}

type interactionCreateEventHandler func(*Session, *InteractionCreate)

func (eh interactionCreateEventHandler) Type() string {
	return interactionCreateEventType
}

func (eh interactionCreateEventHandler) New() interface{} {
	return &InteractionCreate{}
}

func (eh interactionCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*InteractionCreate); ok {
		eh(s, t)
	}
}

type inviteCreateEventHandler func(*Session, *InviteCreate)

func (eh inviteCreateEventHandler) Type() string {
	return inviteCreateEventType
}

func (eh inviteCreateEventHandler) New() interface{} {
	return &InviteCreate{}
}

func (eh inviteCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*InviteCreate); ok {
		eh(s, t)
	}
}

type inviteDeleteEventHandler func(*Session, *InviteDelete)

func (eh inviteDeleteEventHandler) Type() string {
	return inviteDeleteEventType
}

func (eh inviteDeleteEventHandler) New() interface{} {
	return &InviteDelete{}
}

func (eh inviteDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*InviteDelete); ok {
		eh(s, t)
	}
}

type messageCreateEventHandler func(*Session, *MessageCreate)

func (eh messageCreateEventHandler) Type() string {
	return messageCreateEventType
}

func (eh messageCreateEventHandler) New() interface{} {
	return &MessageCreate{}
}

func (eh messageCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageCreate); ok {
		eh(s, t)
	}
}

type messageDeleteEventHandler func(*Session, *MessageDelete)

func (eh messageDeleteEventHandler) Type() string {
	return messageDeleteEventType
}

func (eh messageDeleteEventHandler) New() interface{} {
	return &MessageDelete{}
}

func (eh messageDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageDelete); ok {
		eh(s, t)
	}
}

type messageDeleteBulkEventHandler func(*Session, *MessageDeleteBulk)

func (eh messageDeleteBulkEventHandler) Type() string {
	return messageDeleteBulkEventType
}

func (eh messageDeleteBulkEventHandler) New() interface{} {
	return &MessageDeleteBulk{}
}

func (eh messageDeleteBulkEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageDeleteBulk); ok {
		eh(s, t)
	}
}

type messagePollVoteAddEventHandler func(*Session, *MessagePollVoteAdd)

func (eh messagePollVoteAddEventHandler) Type() string {
	return messagePollVoteAddEventType
}

func (eh messagePollVoteAddEventHandler) New() interface{} {
	return &MessagePollVoteAdd{}
}

func (eh messagePollVoteAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessagePollVoteAdd); ok {
		eh(s, t)
	}
}

type messagePollVoteRemoveEventHandler func(*Session, *MessagePollVoteRemove)

func (eh messagePollVoteRemoveEventHandler) Type() string {
	return messagePollVoteRemoveEventType
}

func (eh messagePollVoteRemoveEventHandler) New() interface{} {
	return &MessagePollVoteRemove{}
}

func (eh messagePollVoteRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessagePollVoteRemove); ok {
		eh(s, t)
	}
}

type messageReactionAddEventHandler func(*Session, *MessageReactionAdd)

func (eh messageReactionAddEventHandler) Type() string {
	return messageReactionAddEventType
}

func (eh messageReactionAddEventHandler) New() interface{} {
	return &MessageReactionAdd{}
}

func (eh messageReactionAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageReactionAdd); ok {
		eh(s, t)
	}
}

type messageReactionRemoveEventHandler func(*Session, *MessageReactionRemove)

func (eh messageReactionRemoveEventHandler) Type() string {
	return messageReactionRemoveEventType
}

func (eh messageReactionRemoveEventHandler) New() interface{} {
	return &MessageReactionRemove{}
}

func (eh messageReactionRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageReactionRemove); ok {
		eh(s, t)
	}
}

type messageReactionRemoveAllEventHandler func(*Session, *MessageReactionRemoveAll)

func (eh messageReactionRemoveAllEventHandler) Type() string {
	return messageReactionRemoveAllEventType
}

func (eh messageReactionRemoveAllEventHandler) New() interface{} {
	return &MessageReactionRemoveAll{}
}

func (eh messageReactionRemoveAllEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageReactionRemoveAll); ok {
		eh(s, t)
	}
}

type messageUpdateEventHandler func(*Session, *MessageUpdate)

func (eh messageUpdateEventHandler) Type() string {
	return messageUpdateEventType
}

func (eh messageUpdateEventHandler) New() interface{} {
	return &MessageUpdate{}
}

func (eh messageUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageUpdate); ok {
		eh(s, t)
	}
}

type presenceUpdateEventHandler func(*Session, *PresenceUpdate)

func (eh presenceUpdateEventHandler) Type() string {
	return presenceUpdateEventType
}

func (eh presenceUpdateEventHandler) New() interface{} {
	return &PresenceUpdate{}
}

func (eh presenceUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*PresenceUpdate); ok {
		eh(s, t)
	}
}

type presencesReplaceEventHandler func(*Session, *PresencesReplace)

func (eh presencesReplaceEventHandler) Type() string {
	return presencesReplaceEventType
}

func (eh presencesReplaceEventHandler) New() interface{} {
	return &PresencesReplace{}
}

func (eh presencesReplaceEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*PresencesReplace); ok {
		eh(s, t)
	}
}

type rateLimitEventHandler func(*Session, *RateLimit)

func (eh rateLimitEventHandler) Type() string {
	return rateLimitEventType
}

func (eh rateLimitEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*RateLimit); ok {
		eh(s, t)
	}
}

type readyEventHandler func(*Session, *Ready)

func (eh readyEventHandler) Type() string {
	return readyEventType
}

func (eh readyEventHandler) New() interface{} {
	return &Ready{}
}

func (eh readyEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Ready); ok {
		eh(s, t)
	}
}

type resumedEventHandler func(*Session, *Resumed)

func (eh resumedEventHandler) Type() string {
	return resumedEventType
}

func (eh resumedEventHandler) New() interface{} {
	return &Resumed{}
}

func (eh resumedEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Resumed); ok {
		eh(s, t)
	}
}

type stageInstanceEventCreateEventHandler func(*Session, *StageInstanceEventCreate)

func (eh stageInstanceEventCreateEventHandler) Type() string {
	return stageInstanceEventCreateEventType
}

func (eh stageInstanceEventCreateEventHandler) New() interface{} {
	return &StageInstanceEventCreate{}
}

func (eh stageInstanceEventCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*StageInstanceEventCreate); ok {
		eh(s, t)
	}
}

type stageInstanceEventDeleteEventHandler func(*Session, *StageInstanceEventDelete)

func (eh stageInstanceEventDeleteEventHandler) Type() string {
	return stageInstanceEventDeleteEventType
}

func (eh stageInstanceEventDeleteEventHandler) New() interface{} {
	return &StageInstanceEventDelete{}
}

func (eh stageInstanceEventDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*StageInstanceEventDelete); ok {
		eh(s, t)
	}
}

type stageInstanceEventUpdateEventHandler func(*Session, *StageInstanceEventUpdate)

func (eh stageInstanceEventUpdateEventHandler) Type() string {
	return stageInstanceEventUpdateEventType
}

func (eh stageInstanceEventUpdateEventHandler) New() interface{} {
	return &StageInstanceEventUpdate{}
}

func (eh stageInstanceEventUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*StageInstanceEventUpdate); ok {
		eh(s, t)
	}
}

type subscriptionCreateEventHandler func(*Session, *SubscriptionCreate)

func (eh subscriptionCreateEventHandler) Type() string {
	return subscriptionCreateEventType
}

func (eh subscriptionCreateEventHandler) New() interface{} {
	return &SubscriptionCreate{}
}

func (eh subscriptionCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*SubscriptionCreate); ok {
		eh(s, t)
	}
}

type subscriptionDeleteEventHandler func(*Session, *SubscriptionDelete)

func (eh subscriptionDeleteEventHandler) Type() string {
	return subscriptionDeleteEventType
}

func (eh subscriptionDeleteEventHandler) New() interface{} {
	return &SubscriptionDelete{}
}

func (eh subscriptionDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*SubscriptionDelete); ok {
		eh(s, t)
	}
}

type subscriptionUpdateEventHandler func(*Session, *SubscriptionUpdate)

func (eh subscriptionUpdateEventHandler) Type() string {
	return subscriptionUpdateEventType
}

func (eh subscriptionUpdateEventHandler) New() interface{} {
	return &SubscriptionUpdate{}
}

func (eh subscriptionUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*SubscriptionUpdate); ok {
		eh(s, t)
	}
}

type threadCreateEventHandler func(*Session, *ThreadCreate)

func (eh threadCreateEventHandler) Type() string {
	return threadCreateEventType
}

func (eh threadCreateEventHandler) New() interface{} {
	return &ThreadCreate{}
}

func (eh threadCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ThreadCreate); ok {
		eh(s, t)
	}
}

type threadDeleteEventHandler func(*Session, *ThreadDelete)

func (eh threadDeleteEventHandler) Type() string {
	return threadDeleteEventType
}

func (eh threadDeleteEventHandler) New() interface{} {
	return &ThreadDelete{}
}

func (eh threadDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ThreadDelete); ok {
		eh(s, t)
	}
}

type threadListSyncEventHandler func(*Session, *ThreadListSync)

func (eh threadListSyncEventHandler) Type() string {
	return threadListSyncEventType
}

func (eh threadListSyncEventHandler) New() interface{} {
	return &ThreadListSync{}
}

func (eh threadListSyncEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ThreadListSync); ok {
		eh(s, t)
	}
}

type threadMemberUpdateEventHandler func(*Session, *ThreadMemberUpdate)

func (eh threadMemberUpdateEventHandler) Type() string {
	return threadMemberUpdateEventType
}

func (eh threadMemberUpdateEventHandler) New() interface{} {
	return &ThreadMemberUpdate{}
}

func (eh threadMemberUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ThreadMemberUpdate); ok {
		eh(s, t)
	}
}

type threadMembersUpdateEventHandler func(*Session, *ThreadMembersUpdate)

func (eh threadMembersUpdateEventHandler) Type() string {
	return threadMembersUpdateEventType
}

func (eh threadMembersUpdateEventHandler) New() interface{} {
	return &ThreadMembersUpdate{}
}

func (eh threadMembersUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ThreadMembersUpdate); ok {
		eh(s, t)
	}
}

type threadUpdateEventHandler func(*Session, *ThreadUpdate)

func (eh threadUpdateEventHandler) Type() string {
	return threadUpdateEventType
}

func (eh threadUpdateEventHandler) New() interface{} {
	return &ThreadUpdate{}
}

func (eh threadUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ThreadUpdate); ok {
		eh(s, t)
	}
}

type typingStartEventHandler func(*Session, *TypingStart)

func (eh typingStartEventHandler) Type() string {
	return typingStartEventType
}

func (eh typingStartEventHandler) New() interface{} {
	return &TypingStart{}
}

func (eh typingStartEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*TypingStart); ok {
		eh(s, t)
	}
}

type userUpdateEventHandler func(*Session, *UserUpdate)

func (eh userUpdateEventHandler) Type() string {
	return userUpdateEventType
}

func (eh userUpdateEventHandler) New() interface{} {
	return &UserUpdate{}
}

func (eh userUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*UserUpdate); ok {
		eh(s, t)
	}
}

type voiceServerUpdateEventHandler func(*Session, *VoiceServerUpdate)

func (eh voiceServerUpdateEventHandler) Type() string {
	return voiceServerUpdateEventType
}

func (eh voiceServerUpdateEventHandler) New() interface{} {
	return &VoiceServerUpdate{}
}

func (eh voiceServerUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*VoiceServerUpdate); ok {
		eh(s, t)
	}
}

type voiceStateUpdateEventHandler func(*Session, *VoiceStateUpdate)

func (eh voiceStateUpdateEventHandler) Type() string {
	return voiceStateUpdateEventType
}

func (eh voiceStateUpdateEventHandler) New() interface{} {
	return &VoiceStateUpdate{}
}

func (eh voiceStateUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*VoiceStateUpdate); ok {
		eh(s, t)
	}
}

type webhooksUpdateEventHandler func(*Session, *WebhooksUpdate)

func (eh webhooksUpdateEventHandler) Type() string {
	return webhooksUpdateEventType
}

func (eh webhooksUpdateEventHandler) New() interface{} {
	return &WebhooksUpdate{}
}

func (eh webhooksUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*WebhooksUpdate); ok {
		eh(s, t)
	}
}

func handlerForInterface(handler interface{}) EventHandler {
	switch v := handler.(type) {
	case func(*Session, interface{}):
		return interfaceEventHandler(v)
	case func(*Session, *ApplicationCommandPermissionsUpdate):
		return applicationCommandPermissionsUpdateEventHandler(v)
	case func(*Session, *AutoModerationActionExecution):
		return autoModerationActionExecutionEventHandler(v)
	case func(*Session, *AutoModerationRuleCreate):
		return autoModerationRuleCreateEventHandler(v)
	case func(*Session, *AutoModerationRuleDelete):
		return autoModerationRuleDeleteEventHandler(v)
	case func(*Session, *AutoModerationRuleUpdate):
		return autoModerationRuleUpdateEventHandler(v)
	case func(*Session, *ChannelCreate):
		return channelCreateEventHandler(v)
	case func(*Session, *ChannelDelete):
		return channelDeleteEventHandler(v)
	case func(*Session, *ChannelPinsUpdate):
		return channelPinsUpdateEventHandler(v)
	case func(*Session, *ChannelUpdate):
		return channelUpdateEventHandler(v)
	case func(*Session, *Connect):
		return connectEventHandler(v)
	case func(*Session, *Disconnect):
		return disconnectEventHandler(v)
	case func(*Session, *EntitlementCreate):
		return entitlementCreateEventHandler(v)
	case func(*Session, *EntitlementDelete):
		return entitlementDeleteEventHandler(v)
	case func(*Session, *EntitlementUpdate):
		return entitlementUpdateEventHandler(v)
	case func(*Session, *Event):
		return eventEventHandler(v)
	case func(*Session, *GuildAuditLogEntryCreate):
		return guildAuditLogEntryCreateEventHandler(v)
	case func(*Session, *GuildBanAdd):
		return guildBanAddEventHandler(v)
	case func(*Session, *GuildBanRemove):
		return guildBanRemoveEventHandler(v)
	case func(*Session, *GuildCreate):
		return guildCreateEventHandler(v)
	case func(*Session, *GuildDelete):
		return guildDeleteEventHandler(v)
	case func(*Session, *GuildEmojisUpdate):
		return guildEmojisUpdateEventHandler(v)
	case func(*Session, *GuildIntegrationsUpdate):
		return guildIntegrationsUpdateEventHandler(v)
	case func(*Session, *GuildMemberAdd):
		return guildMemberAddEventHandler(v)
	case func(*Session, *GuildMemberRemove):
		return guildMemberRemoveEventHandler(v)
	case func(*Session, *GuildMemberUpdate):
		return guildMemberUpdateEventHandler(v)
	case func(*Session, *GuildMembersChunk):
		return guildMembersChunkEventHandler(v)
	case func(*Session, *GuildRoleCreate):
		return guildRoleCreateEventHandler(v)
	case func(*Session, *GuildRoleDelete):
		return guildRoleDeleteEventHandler(v)
	case func(*Session, *GuildRoleUpdate):
		return guildRoleUpdateEventHandler(v)
	case func(*Session, *GuildScheduledEventCreate):
		return guildScheduledEventCreateEventHandler(v)
	case func(*Session, *GuildScheduledEventDelete):
		return guildScheduledEventDeleteEventHandler(v)
	case func(*Session, *GuildScheduledEventUpdate):
		return guildScheduledEventUpdateEventHandler(v)
	case func(*Session, *GuildScheduledEventUserAdd):
		return guildScheduledEventUserAddEventHandler(v)
	case func(*Session, *GuildScheduledEventUserRemove):
		return guildScheduledEventUserRemoveEventHandler(v)
	case func(*Session, *GuildUpdate):
		return guildUpdateEventHandler(v)
	case func(*Session, *IntegrationCreate):
		return integrationCreateEventHandler(v)
	case func(*Session, *IntegrationDelete):
		return integrationDeleteEventHandler(v)
	case func(*Session, *IntegrationUpdate):
		return integrationUpdateEventHandler(v)
	case func(*Session, *InteractionCreate):
		return interactionCreateEventHandler(v)
	case func(*Session, *InviteCreate):
		return inviteCreateEventHandler(v)
	case func(*Session, *InviteDelete):
		return inviteDeleteEventHandler(v)
	case func(*Session, *MessageCreate):
		return messageCreateEventHandler(v)
	case func(*Session, *MessageDelete):
		return messageDeleteEventHandler(v)
	case func(*Session, *MessageDeleteBulk):
		return messageDeleteBulkEventHandler(v)
	case func(*Session, *MessagePollVoteAdd):
		return messagePollVoteAddEventHandler(v)
	case func(*Session, *MessagePollVoteRemove):
		return messagePollVoteRemoveEventHandler(v)
	case func(*Session, *MessageReactionAdd):
		return messageReactionAddEventHandler(v)
	case func(*Session, *MessageReactionRemove):
		return messageReactionRemoveEventHandler(v)
	case func(*Session, *MessageReactionRemoveAll):
		return messageReactionRemoveAllEventHandler(v)
	case func(*Session, *MessageUpdate):
		return messageUpdateEventHandler(v)
	case func(*Session, *PresenceUpdate):
		return presenceUpdateEventHandler(v)
	case func(*Session, *PresencesReplace):
		return presencesReplaceEventHandler(v)
	case func(*Session, *RateLimit):
		return rateLimitEventHandler(v)
	case func(*Session, *Ready):
		return readyEventHandler(v)
	case func(*Session, *Resumed):
		return resumedEventHandler(v)
	case func(*Session, *StageInstanceEventCreate):
		return stageInstanceEventCreateEventHandler(v)
	case func(*Session, *StageInstanceEventDelete):
		return stageInstanceEventDeleteEventHandler(v)
	case func(*Session, *StageInstanceEventUpdate):
		return stageInstanceEventUpdateEventHandler(v)
	case func(*Session, *SubscriptionCreate):
		return subscriptionCreateEventHandler(v)
	case func(*Session, *SubscriptionDelete):
		return subscriptionDeleteEventHandler(v)
	case func(*Session, *SubscriptionUpdate):
		return subscriptionUpdateEventHandler(v)
	case func(*Session, *ThreadCreate):
		return threadCreateEventHandler(v)
	case func(*Session, *ThreadDelete):
		return threadDeleteEventHandler(v)
	case func(*Session, *ThreadListSync):
		return threadListSyncEventHandler(v)
	case func(*Session, *ThreadMemberUpdate):
		return threadMemberUpdateEventHandler(v)
	case func(*Session, *ThreadMembersUpdate):
		return threadMembersUpdateEventHandler(v)
	case func(*Session, *ThreadUpdate):
		return threadUpdateEventHandler(v)
	case func(*Session, *TypingStart):
		return typingStartEventHandler(v)
	case func(*Session, *UserUpdate):
		return userUpdateEventHandler(v)
	case func(*Session, *VoiceServerUpdate):
		return voiceServerUpdateEventHandler(v)
	case func(*Session, *VoiceStateUpdate):
		return voiceStateUpdateEventHandler(v)
	case func(*Session, *WebhooksUpdate):
		return webhooksUpdateEventHandler(v)
	}

	return nil
}

func init() {
	registerInterfaceProvider(applicationCommandPermissionsUpdateEventHandler(nil))
	registerInterfaceProvider(autoModerationActionExecutionEventHandler(nil))
	registerInterfaceProvider(autoModerationRuleCreateEventHandler(nil))
	registerInterfaceProvider(autoModerationRuleDeleteEventHandler(nil))
	registerInterfaceProvider(autoModerationRuleUpdateEventHandler(nil))
	registerInterfaceProvider(channelCreateEventHandler(nil))
	registerInterfaceProvider(channelDeleteEventHandler(nil))
	registerInterfaceProvider(channelPinsUpdateEventHandler(nil))
	registerInterfaceProvider(channelUpdateEventHandler(nil))
	registerInterfaceProvider(entitlementCreateEventHandler(nil))
	registerInterfaceProvider(entitlementDeleteEventHandler(nil))
	registerInterfaceProvider(entitlementUpdateEventHandler(nil))
	registerInterfaceProvider(guildAuditLogEntryCreateEventHandler(nil))
	registerInterfaceProvider(guildBanAddEventHandler(nil))
	registerInterfaceProvider(guildBanRemoveEventHandler(nil))
	registerInterfaceProvider(guildCreateEventHandler(nil))
	registerInterfaceProvider(guildDeleteEventHandler(nil))
	registerInterfaceProvider(guildEmojisUpdateEventHandler(nil))
	registerInterfaceProvider(guildIntegrationsUpdateEventHandler(nil))
	registerInterfaceProvider(guildMemberAddEventHandler(nil))
	registerInterfaceProvider(guildMemberRemoveEventHandler(nil))
	registerInterfaceProvider(guildMemberUpdateEventHandler(nil))
	registerInterfaceProvider(guildMembersChunkEventHandler(nil))
	registerInterfaceProvider(guildRoleCreateEventHandler(nil))
	registerInterfaceProvider(guildRoleDeleteEventHandler(nil))
	registerInterfaceProvider(guildRoleUpdateEventHandler(nil))
	registerInterfaceProvider(guildScheduledEventCreateEventHandler(nil))
	registerInterfaceProvider(guildScheduledEventDeleteEventHandler(nil))
	registerInterfaceProvider(guildScheduledEventUpdateEventHandler(nil))
	registerInterfaceProvider(guildScheduledEventUserAddEventHandler(nil))
	registerInterfaceProvider(guildScheduledEventUserRemoveEventHandler(nil))
	registerInterfaceProvider(guildUpdateEventHandler(nil))
	registerInterfaceProvider(integrationCreateEventHandler(nil))
	registerInterfaceProvider(integrationDeleteEventHandler(nil))
	registerInterfaceProvider(integrationUpdateEventHandler(nil))
	registerInterfaceProvider(interactionCreateEventHandler(nil))
	registerInterfaceProvider(inviteCreateEventHandler(nil))
	registerInterfaceProvider(inviteDeleteEventHandler(nil))
	registerInterfaceProvider(messageCreateEventHandler(nil))
	registerInterfaceProvider(messageDeleteEventHandler(nil))
	registerInterfaceProvider(messageDeleteBulkEventHandler(nil))
	registerInterfaceProvider(messagePollVoteAddEventHandler(nil))
	registerInterfaceProvider(messagePollVoteRemoveEventHandler(nil))
	registerInterfaceProvider(messageReactionAddEventHandler(nil))
	registerInterfaceProvider(messageReactionRemoveEventHandler(nil))
	registerInterfaceProvider(messageReactionRemoveAllEventHandler(nil))
	registerInterfaceProvider(messageUpdateEventHandler(nil))
	registerInterfaceProvider(presenceUpdateEventHandler(nil))
	registerInterfaceProvider(presencesReplaceEventHandler(nil))
	registerInterfaceProvider(readyEventHandler(nil))
	registerInterfaceProvider(resumedEventHandler(nil))
	registerInterfaceProvider(stageInstanceEventCreateEventHandler(nil))
	registerInterfaceProvider(stageInstanceEventDeleteEventHandler(nil))
	registerInterfaceProvider(stageInstanceEventUpdateEventHandler(nil))
	registerInterfaceProvider(subscriptionCreateEventHandler(nil))
	registerInterfaceProvider(subscriptionDeleteEventHandler(nil))
	registerInterfaceProvider(subscriptionUpdateEventHandler(nil))
	registerInterfaceProvider(threadCreateEventHandler(nil))
	registerInterfaceProvider(threadDeleteEventHandler(nil))
	registerInterfaceProvider(threadListSyncEventHandler(nil))
	registerInterfaceProvider(threadMemberUpdateEventHandler(nil))
	registerInterfaceProvider(threadMembersUpdateEventHandler(nil))
	registerInterfaceProvider(threadUpdateEventHandler(nil))
	registerInterfaceProvider(typingStartEventHandler(nil))
	registerInterfaceProvider(userUpdateEventHandler(nil))
	registerInterfaceProvider(voiceServerUpdateEventHandler(nil))
	registerInterfaceProvider(voiceStateUpdateEventHandler(nil))
	registerInterfaceProvider(webhooksUpdateEventHandler(nil))
}
