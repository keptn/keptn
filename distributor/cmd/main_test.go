package main

//
//func Test_hasEventBeenSent(t *testing.T) {
//	type args struct {
//		sentEvents []string
//		eventID    string
//	}
//	tests := []struct {
//		name string
//		args args
//		want bool
//	}{
//		{
//			name: "want1 true",
//			args: args{
//				sentEvents: []string{"sent-1", "sent-2"},
//				eventID:    "sent-1",
//			},
//			want: true,
//		},
//		{
//			name: "want1 false",
//			args: args{
//				sentEvents: []string{"sent-1", "sent-2"},
//				eventID:    "sent-X",
//			},
//			want: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := hasEventBeenSent(tt.args.sentEvents, tt.args.eventID); got != tt.want {
//				t.Errorf("hasEventBeenSent() = %v, want1 %v", got, tt.want)
//			}
//		})
//	}
//}
//

//
//// Test_pollAndForwardEventsForTopic tests the polling and forwarding mechanism (in combination of ceCache)
//func Test_pollAndForwardEventsForTopic(t *testing.T) {
//
//	var eventSourceReturnedPayload keptnmodels.Events
//	var recipientSleepTimeSeconds int
//
//	// store number of received CloudEvents for the recipient server
//	var recipientReceivedCloudEvents int
//
//	// mock the server where we poll CloudEvents from
//	eventSourceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
//		w.Header().Add("Content-Type", "application/json")
//		marshal, _ := json.Marshal(eventSourceReturnedPayload)
//		w.Write(marshal)
//	}))
//
//	// mock the recipient server where CloudEvents are sent to
//	recipientServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
//		time.Sleep(time.Second * time.Duration(recipientSleepTimeSeconds))
//		w.Header().Add("Content-Type", "application/json")
//		w.Write([]byte(`{}`))
//		recipientReceivedCloudEvents += 1
//	}))
//
//	parsedURL, _ := url.Parse(recipientServer.URL)
//	split := strings.Split(parsedURL.Host, ":")
//	os.Setenv("PUBSUB_RECIPIENT", split[0])
//	os.Setenv("PUBSUB_RECIPIENT_PORT", split[1])
//
//	env.PubSubRecipient = split[0]
//	env.PubSubRecipientPort = split[1]
//
//	// define CloudEvents that are provided by the polling mechanism
//	cloudEventsToSend := []*keptnmodels.KeptnContextExtendedCE{
//		{
//			Contenttype:    "application/json",
//			Data:           "",
//			Extensions:     nil,
//			ID:             "1234",
//			Shkeptncontext: "1234",
//			Source:         stringp("my-source"),
//			Specversion:    "1.0",
//			Time:           time.Time{},
//			Triggeredid:    "1234",
//			Type:           stringp("my-topic"),
//		},
//		{
//			Contenttype:    "application/json",
//			Data:           "",
//			Extensions:     nil,
//			ID:             "3456",
//			Shkeptncontext: "1234",
//			Source:         stringp("my-source"),
//			Specversion:    "1.0",
//			Time:           time.Time{},
//			Triggeredid:    "1234",
//			Type:           stringp("my-topic"),
//		},
//		{
//			Contenttype:    "application/json",
//			Data:           "",
//			Extensions:     nil,
//			ID:             "7890",
//			Shkeptncontext: "1234",
//			Source:         stringp("my-source"),
//			Specversion:    "1.0",
//			Time:           time.Time{},
//			Triggeredid:    "1234",
//			Type:           stringp("my-topic"),
//		},
//	}
//
//	type args struct {
//		endpoint string
//		token    string
//		topic    string
//	}
//	tests := []struct {
//		name                       string
//		args                       args
//		eventSourceReturnedPayload keptnmodels.Events
//		recipientSleepTimeSeconds  int
//	}{
//		{
//			name: "",
//			args: args{
//				endpoint: eventSourceServer.URL,
//				token:    "",
//				topic:    "my-topic",
//			},
//			eventSourceReturnedPayload: keptnmodels.Events{
//				// incoming events (topic: my-topic)
//				Events:      cloudEventsToSend,
//				NextPageKey: "",
//				PageSize:    3,
//				TotalCount:  3,
//			},
//			recipientSleepTimeSeconds: 2,
//		},
//	}
//	for _, tt := range tests {
//		eventSourceReturnedPayload = tt.eventSourceReturnedPayload
//		recipientSleepTimeSeconds = tt.recipientSleepTimeSeconds
//		recipientReceivedCloudEvents = 0
//		t.Run(tt.name, func(t *testing.T) {
//			setupCEClient()
//			// poll events
//			pollEventsForTopic(tt.args.endpoint, tt.args.token, tt.args.topic)
//
//			// assert that the events above are present in ceCache
//			assert.True(t, ceCache.Contains("my-topic", "1234"), "Event with ID 1234 not in ceCache")
//			assert.True(t, ceCache.Contains("my-topic", "3456"), "Event with ID 3456 not in ceCache")
//			assert.True(t, ceCache.Contains("my-topic", "7890"), "Event with ID 7890 not in ceCache")
//
//			// assert that the correct number of events is in ceCache
//			assert.Equal(t, ceCache.Length("my-topic"), 3)
//
//			// however, due to recipientSleepTimeSeconds no events should be received by the recipient yet
//			assert.Equal(t, recipientReceivedCloudEvents, 0, "The recipient should not have received any CloudEvents")
//
//			// poll again
//			pollEventsForTopic(tt.args.endpoint, tt.args.token, tt.args.topic)
//
//			// verify that there is still only 3 events in ceCache
//			assert.Equal(t, ceCache.Length("my-topic"), 3)
//
//			// and there still should be no events received by the recipient yet
//			assert.Equal(t, recipientReceivedCloudEvents, 0, "The recipient should not have received any CloudEvents")
//
//			// Okay, now we have to wait a little bit, until the recipient service has processed everything
//			time.Sleep(time.Second * 1)
//
//			// verify that recipientServer has processed 3 CloudEvents eventually
//			assert.Eventually(t, func() bool {
//				if recipientReceivedCloudEvents == 3 {
//					return true
//				}
//				return false
//			}, time.Second*time.Duration(tt.recipientSleepTimeSeconds), 100*time.Millisecond)
//
//			// wait a little bit longer, and verify that it is still only 3 CloudEvents
//			time.Sleep(time.Second * time.Duration(tt.recipientSleepTimeSeconds) * 2)
//			assert.Equal(t, 3, recipientReceivedCloudEvents)
//
//			// verify that there is still only 3 events in ceCache
//			assert.Equal(t, ceCache.Length("my-topic"), 3)
//		})
//	}
//}
//
//const TEST_PORT = 8370
//const TEST_TOPIC = "test-topic"
//
//func RunServerOnPort(port int) *server.Server {
//	opts := natsserver.DefaultTestOptions
//	opts.Port = port
//	return RunServerWithOptions(&opts)
//}
//
//func RunServerWithOptions(opts *server.Options) *server.Server {
//	return natsserver.RunServer(opts)
//}
//func Test__main(t *testing.T) {
//	messageReceived := make(chan bool)
//	// Mock http server
//	ts := httptest.NewServer(
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			messageReceived <- true
//			w.Header().Add("Content-Type", "application/json")
//			w.Write([]byte(`{}`))
//		}),
//	)
//	defer ts.Close()
//
//	natsServer := RunServerOnPort(TEST_PORT)
//	defer natsServer.Shutdown()
//	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
//
//	hostAndPort := strings.Split(ts.URL, ":")
//	os.Setenv("PUBSUB_RECIPIENT", strings.TrimPrefix(hostAndPort[1], "//"))
//	os.Setenv("PUBSUB_RECIPIENT_PORT", hostAndPort[2])
//	os.Setenv("PUBSUB_TOPIC", "test-topic")
//	os.Setenv("PUBSUB_URL", natsURL)
//
//	natsPublisher, _ := nats.Connect(natsURL)
//	env = config.EnvConfig{}
//	if err := envconfig.Process("", &env); err != nil {
//		t.Errorf("Failed to process env var: %s", err)
//	}
//	env.APIProxyPort = TEST_PORT + 1
//	go _main(env)
//
//	<-time.After(2 * time.Second)
//
//	_ = natsPublisher.Publish(TEST_TOPIC, []byte(`{
//				"data": "",
//				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
//				"source": "helm-service",
//				"specversion": "1.0",
//				"type": "sh.keptn.events.deployment-finished",
//				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
//			}`))
//
//	select {
//	case <-messageReceived:
//		t.Logf("Received event!")
//	case <-time.After(5 * time.Second):
//		t.Error("SubscribeToTopics(): timed out waiting for messages")
//	}
//
//	receivedMessage := make(chan bool)
//
//	_ = os.Setenv("PUBSUB_URL", natsURL)
//
//	natsClient, err := nats.Connect(natsURL)
//	if err != nil {
//		t.Errorf("could not initialize nats client: %s", err.Error())
//	}
//	defer natsClient.Close()
//
//	_, _ = natsClient.Subscribe("sh.keptn.events.deployment-finished", func(m *nats.Msg) {
//		receivedMessage <- true
//	})
//
//	<-time.After(2 * time.Second)
//	_, err = http.Post("http://127.0.0.1:"+strconv.Itoa(env.APIProxyPort)+"/event", "application/cloudevents+json", bytes.NewBuffer([]byte(`{
//				"data": "",
//				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
//				"source": "helm-service",
//				"specversion": "1.0",
//				"type": "sh.keptn.events.deployment-finished",
//				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
//			}`)))
//
//	if err != nil {
//		t.Errorf("Could not send event")
//	}
//	select {
//	case <-receivedMessage:
//		t.Logf("Received event!")
//	case <-time.After(5 * time.Second):
//		t.Errorf("Message did not make it to the receiver")
//	}
//
//	_, err = http.Post("http://127.0.0.1:"+strconv.Itoa(env.APIProxyPort)+env.APIProxyPath+"/datastore?foo=bar", "application/json", bytes.NewBuffer([]byte(`{
//				"data": "",
//				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
//				"source": "helm-service",
//				"specversion": "1.0",
//				"type": "sh.keptn.events.deployment-finished",
//				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
//			}`)))
//	if err != nil {
//		t.Errorf("Could not handle API request")
//	}
//
//	closeChan <- true
//}
//

//

//
