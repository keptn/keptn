package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	models "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

// AddWatchFlag ads the --watch flag to the command
func AddWatchFlag(cmd *cobra.Command) *bool {
	return cmd.Flags().BoolP("watch", "w", false, "Print event stream")
}

// AddWatchTimeFlag adds the --watch-time flag to the command
func AddWatchTimeFlag(cmd *cobra.Command) *int {
	return cmd.Flags().Int("watch-time", math.MaxInt32, "Timeout (in seconds) used for the --watch flag")
}

// AddOutputFormatFlag adds the --output flag to the command
func AddOutputFormatFlag(cmd *cobra.Command) *string {
	return cmd.Flags().StringP("output", "o", "",
		"Output format for the --watch flag. One of: json|yaml")
}

// NewDefaultWatcher creates a preconfigured EventWatcher used in cli commands
func NewDefaultWatcher(eventHandler apiutils.EventHandlerInterface, filter apiutils.EventFilter, timeOut time.Duration) *apiutils.EventWatcher {
	watcher := apiutils.NewEventWatcher(
		eventHandler,
		apiutils.WithEventFilter(filter),
		apiutils.WithInterval(time.NewTicker(5*time.Second)),
		apiutils.WithStartTime(time.Time{}), // this makes sure that we also capture old events
		apiutils.WithTimeout(timeOut),
	)
	return watcher
}

// Watcher is the interface to an EventWatcher
type Watcher interface {
	Watch(ctx context.Context) (<-chan []*models.KeptnContextExtendedCE, context.CancelFunc)
}

// PrintEventWatcher uses the given watcher type and prints its result to the given writer in the given format
// Note that this function is blocking until the watcher is canceled
func PrintEventWatcher(context context.Context, watcher Watcher, format string, writer io.Writer) {
	eventChan, _ := watcher.Watch(context)
	for events := range eventChan {
		for _, e := range events {
			PrintEvents(writer, format, *e)
		}
	}
}

// PrintEvents can be used to print events (and structs in general) to the given write
// either in YAML or JSON format
func PrintEvents(writer io.Writer, format string, content interface{}) {
	if format == "yaml" {
		PrintAsYAML(writer, content)
	} else { //default
		PrintAsJSON(writer, content)
	}
}

// PrintAsYAML prints events in YAML format to std::out
func PrintAsYAML(writer io.Writer, events interface{}) {
	eventsYAML, _ := yaml.Marshal(events)
	fmt.Fprintf(writer, "%s\n", string(eventsYAML))

}

// PrintAsJSON prints events in JSON format to std::out
func PrintAsJSON(writer io.Writer, events interface{}) {
	eventsJSON, _ := json.MarshalIndent(events, "", "    ")
	fmt.Fprintf(writer, "%s\n", string(eventsJSON))
}

func LookupHostname(hostname string, lookupFn resolveFunc, sleepFn sleepFunc) bool {
	if strings.HasSuffix(hostname, "xip.io") {
		logging.PrintLog("Skipping lookup of xip.io domain", logging.InfoLevel)
		return true
	}
	// first, try to resolve the domain (and retry it)
	for retries := 0; retries < 3; sleepFn(5 * time.Second) {
		_, err := lookupFn(hostname)
		if err != nil {
			logging.PrintLog("Failed to resolve hostname "+hostname, logging.InfoLevel)
			logging.PrintLog("Retrying...", logging.InfoLevel)
			retries++
		} else {
			return true
		}
	}

	return false
}

func CheckEndpointStatus(endPoint string) error {
	if checkEndPointStatusMock {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxHTTPTimeout)
	defer cancel()

	req, err := http.NewRequest("HEAD", endPoint, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
