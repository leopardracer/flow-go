package data_providers

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/engine/access/rest/http/request"
	"github.com/onflow/flow-go/engine/access/rest/websockets/models"
	"github.com/onflow/flow-go/engine/access/state_stream"
	"github.com/onflow/flow-go/engine/access/state_stream/backend"
	"github.com/onflow/flow-go/engine/access/subscription"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/counters"
)

// eventsArguments contains the arguments a user passes to subscribe to events
type eventsArguments struct {
	StartBlockID      flow.Identifier          // ID of the block to start subscription from
	StartBlockHeight  uint64                   // Height of the block to start subscription from
	Filter            state_stream.EventFilter // Filter applied to events for a given subscription
	HeartbeatInterval *uint64                  // Maximum number of blocks message won't be sent. Nil if not set
}

// EventsDataProvider is responsible for providing events
type EventsDataProvider struct {
	*baseDataProvider

	logger         zerolog.Logger
	stateStreamApi state_stream.API

	heartbeatInterval uint64
}

var _ DataProvider = (*EventsDataProvider)(nil)

// NewEventsDataProvider creates a new instance of EventsDataProvider.
func NewEventsDataProvider(
	ctx context.Context,
	logger zerolog.Logger,
	stateStreamApi state_stream.API,
	subscriptionID string,
	topic string,
	arguments models.Arguments,
	send chan<- interface{},
	chain flow.Chain,
	eventFilterConfig state_stream.EventFilterConfig,
	heartbeatInterval uint64,
) (*EventsDataProvider, error) {
	if stateStreamApi == nil {
		return nil, fmt.Errorf("this access node does not support streaming events")
	}

	p := &EventsDataProvider{
		logger:            logger.With().Str("component", "events-data-provider").Logger(),
		stateStreamApi:    stateStreamApi,
		heartbeatInterval: heartbeatInterval,
	}

	// Initialize arguments passed to the provider.
	eventArgs, err := parseEventsArguments(arguments, chain, eventFilterConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid arguments for events data provider: %w", err)
	}
	if eventArgs.HeartbeatInterval != nil {
		p.heartbeatInterval = *eventArgs.HeartbeatInterval
	}

	subCtx, cancel := context.WithCancel(ctx)

	p.baseDataProvider = newBaseDataProvider(
		subscriptionID,
		topic,
		arguments,
		cancel,
		send,
		p.createSubscription(subCtx, eventArgs), // Set up a subscription to events based on arguments.
	)

	return p, nil
}

// Run starts processing the subscription for events and handles responses.
//
// Expected errors during normal operations:
//   - context.Canceled: if the operation is canceled, during an unsubscribe action.
func (p *EventsDataProvider) Run() error {
	return subscription.HandleSubscription(p.subscription, p.handleResponse())
}

// handleResponse processes events and sends the formatted response.
//
// No errors are expected during normal operations.
func (p *EventsDataProvider) handleResponse() func(eventsResponse *backend.EventsResponse) error {
	blocksSinceLastMessage := uint64(0)
	messageIndex := counters.NewMonotonicCounter(0)

	return func(eventsResponse *backend.EventsResponse) error {
		// check if there are any events in the response. if not, do not send a message unless the last
		// response was more than HeartbeatInterval blocks ago
		if len(eventsResponse.Events) == 0 {
			blocksSinceLastMessage++
			if blocksSinceLastMessage < p.heartbeatInterval {
				return nil
			}
		}
		blocksSinceLastMessage = 0

		index := messageIndex.Value()
		if ok := messageIndex.Set(messageIndex.Value() + 1); !ok {
			return fmt.Errorf("message index already incremented to: %d", messageIndex.Value())
		}

		var eventsPayload models.EventResponse
		eventsPayload.Build(eventsResponse, index)

		var response models.BaseDataProvidersResponse
		response.Build(p.ID(), p.Topic(), &eventsPayload)

		p.send <- &response

		return nil
	}
}

// createSubscription creates a new subscription using the specified input arguments.
func (p *EventsDataProvider) createSubscription(ctx context.Context, args eventsArguments) subscription.Subscription {
	if args.StartBlockID != flow.ZeroID {
		return p.stateStreamApi.SubscribeEventsFromStartBlockID(ctx, args.StartBlockID, args.Filter)
	}

	if args.StartBlockHeight != request.EmptyHeight {
		return p.stateStreamApi.SubscribeEventsFromStartHeight(ctx, args.StartBlockHeight, args.Filter)
	}

	return p.stateStreamApi.SubscribeEventsFromLatest(ctx, args.Filter)
}

// parseEventsArguments validates and initializes the events arguments.
func parseEventsArguments(
	arguments models.Arguments,
	chain flow.Chain,
	eventFilterConfig state_stream.EventFilterConfig,
) (eventsArguments, error) {
	allowedFields := map[string]struct{}{
		"start_block_id":     {},
		"start_block_height": {},
		"event_types":        {},
		"addresses":          {},
		"contracts":          {},
		"heartbeat_interval": {},
	}
	err := ensureAllowedFields(arguments, allowedFields)
	if err != nil {
		return eventsArguments{}, err
	}

	var args eventsArguments

	// Parse block arguments
	startBlockID, startBlockHeight, err := parseStartBlock(arguments)
	if err != nil {
		return eventsArguments{}, err
	}
	args.StartBlockID = startBlockID
	args.StartBlockHeight = startBlockHeight

	// Parse 'heartbeat_interval' argument
	heartbeatInterval, err := extractHeartbeatInterval(arguments)
	if err != nil {
		return eventsArguments{}, err
	}
	args.HeartbeatInterval = heartbeatInterval

	// Parse 'event_types' as a JSON array
	eventTypes, err := extractArrayOfStrings(arguments, "event_types", false)
	if err != nil {
		return eventsArguments{}, err
	}

	// Parse 'addresses' as []string{}
	addresses, err := extractArrayOfStrings(arguments, "addresses", false)
	if err != nil {
		return eventsArguments{}, err
	}

	// Parse 'contracts' as []string{}
	contracts, err := extractArrayOfStrings(arguments, "contracts", false)
	if err != nil {
		return eventsArguments{}, err
	}

	// Initialize the event filter with the parsed arguments
	filter, err := state_stream.NewEventFilter(eventFilterConfig, chain, eventTypes, addresses, contracts)
	if err != nil {
		return eventsArguments{}, fmt.Errorf("error creating event filter: %w", err)
	}
	args.Filter = filter

	return args, nil
}
