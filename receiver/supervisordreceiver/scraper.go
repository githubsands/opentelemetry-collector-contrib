package supervisordreceiver

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// hypervisorDScraper handles scraping of hypervisorD metrics
type hypervisorDScraper struct {
	logger          *zap.Logger
	aggregationRate time.Duration
	cfg             *Config
	settings        component.TelemetrySettings
	client          *svClient
	mb              *MetricsBuilder
}

// newScraper creates a new scraper
func newHypervisorDScraper(logger *zap.Logger, cfg *Config, settings component.ReceiverCreateSettings) *hypervisorDScraper {
	return &hypervisorDScraper{
		logger:   logger,
		cfg:      cfg,
		settings: settings.TelemetrySettings,
		mb:       NewMetricsBuilder(cfg.MetricsSettings, settings.BuildInfo),
	}
}

// start starts the scraper by creating a new HTTP Client on the scraper
func (r *hypervisorDScraper) start(ctx context.Context, host component.Host) (err error) {
	r.client, err = newSVClient(r.cfg, r.logger)
	if err != nil {
		return errors.New("no client defined")
	}
	return nil
}

func (hvs *hypervisorDScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	// Validate we don't attempt to scrape without initializing the client
	if hvs.client == nil {
		return pmetric.NewMetrics(), errors.New("client not initialized")
	}
	var stats []*Statistic
	stats = hvs.client.getStats(stats)
	metrics := hvs.collectStatistics(stats)
	return metrics, nil
}

func (hvs *hypervisorDScraper) collectStatistics(stats []*Statistic) pmetric.Metrics {
	now := pcommon.NewTimestampFromTime(time.Now())
	hvs.recordProcessCount(now, int64(len(stats)))
	hvs.recordUptimePerProcess(now, stats)
	return hvs.mb.Emit()
}

// TODO: mdatagen should be able to export uint8s
func (hvs *hypervisorDScraper) recordProcessCount(now pcommon.Timestamp, count int64) {
	hvs.mb.RecordSupervisordProcessCountDataPoint(now, count)
}

func (hvs *hypervisorDScraper) recordUptimePerProcess(now pcommon.Timestamp, stats []*Statistic) {
	var start, stop int
	var upTime int64
	var err error
	for _, v := range stats {
		start, err = strconv.Atoi(v.Start)
		if err != nil {
			hvs.logInvalid("integer", err.Error(), "")
		}
		stop, err = strconv.Atoi(v.Stop)
		if err != nil {
			hvs.logInvalid("integer", err.Error(), "")
		}
		upTime = int64(stop - start)
		hvs.mb.RecordSupervisordProcessUptimeDataPoint(now, upTime, AttributeNameName)
	}
}

func (r *hypervisorDScraper) logInvalid(expectedType, key, value string) {
	r.logger.Info(
		"invalid value",
		zap.String("expectedType", expectedType),
		zap.String("key", key),
		zap.String("value", value),
	)
}
