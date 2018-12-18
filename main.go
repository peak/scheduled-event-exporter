package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	time2 "time"
)

const (
	namespace = "aws_events"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of Aws successful.",
		nil, nil,
	)
	scheduledEventsChecks = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduled_events_status"),
		"Amazon EC2 can schedule events for your instances related to hardware issues, software updates, or system maintenance.",
		[]string{"instance_id"}, nil,
	)
)

type Exporter struct {
	client *ec2.EC2
}

// NewExporter returns an initialized Exporter.
func NewExporter(awsRegion *string) (*Exporter, error) {
	// Load session from shared config
	config := aws.Config{Region: awsRegion}
	sess := session.New(&config)
	/*session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))*/

	// Create new EC2 client
	svc := ec2.New(sess)
	if sess == nil {
		log.Fatal("Could not create AWS session")
	}
	// Init our exporter.
	return &Exporter{
		client: svc,
	}, nil
}

// Describe describes all the metrics ever exported by the Aviatrix exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- scheduledEventsChecks
}

// Collect fetches the stats from configured Aviatrix location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Infoln("okay")
	var upValue float64
	var params *ec2.DescribeInstanceStatusInput
	ret, err := e.client.DescribeInstanceStatus(params)
	if err != nil {
		upValue = 0
		fmt.Println("Unable to DescribeInstanceStatus", err)
	} else {
		upValue = 1
	}
	for item := range ret.InstanceStatuses {
		instance_id := *ret.InstanceStatuses[item].InstanceId
		for _, event := range ret.InstanceStatuses[item].Events {
			// Time remaining in hours
			timeRemain := event.NotBefore.Sub(time2.Now()).Hours()
			ch <- prometheus.MustNewConstMetric(
				scheduledEventsChecks, prometheus.GaugeValue, timeRemain, instance_id,
			)
			fmt.Println("Instance: " + instance_id)
			fmt.Println(event)
		}
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, upValue,
	)
}

func init() {
	prometheus.MustRegister(version.NewCollector("aviatrix_exporter"))
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9169").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		awsRegion     = kingpin.Flag("aws.region", "AWS region to connect to.").Default("us-east-1").String()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("aviatrix_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting aws_events_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	exporter, err := NewExporter(awsRegion)
	if err != nil {
		log.Fatalln(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Aws Events Exporter</title></head>
             <body>
             <h1>Aws Events Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             <h2>Options</h2>
             </dl>
             <h2>Build</h2>
             <pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
