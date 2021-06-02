package distributor

//
//import (
//	"context"
//	"fmt"
//	api "github.com/keptn/go-utils/pkg/api/utils"
//	"testing"
//	"time"
//)
//
//func Test_Watcher(t *testing.T) {
//	eventHandler := api.NewAuthenticatedEventHandler("http://35.195.89.83.nip.io/api", "zzUF69pYugbHbkPGIVnzcXeInButNDGbZ0jDCFjU4sqyU", "x-token", nil, "http")
//
//	watcher := api.NewEventWatcher(eventHandler,
//		api.WithEventFilter(api.EventFilter{ // use custom filter
//			Project:      "sockshop",
//			KeptnContext: "df"}),
//		api.WithInterval(time.NewTicker(5*time.Second)), // fetch every 5 seconds
//		api.WithStartTime(time.Now()),                   // start fetching events newer than this timestamp
//		api.WithTimeout(time.Second*15)) // stop fetching events after 15 secs
//
//
//		// start watcher and consume events
//	allEvents, _ := watcher.Watch(context.Background())
//	for events := range allEvents {
//		for _, e := range events {
//			fmt.Println(*e.Type)
//		}
//	}
//}
