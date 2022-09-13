// Code generated by mdatagen. DO NOT EDIT.

package supervisordreceiver

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// MetricSettings provides common settings for a particular metric.
type MetricSettings struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsSettings provides settings for hypervisord receiver metrics.
type MetricsSettings struct {
	SupervisordProcessCount  MetricSettings `mapstructure:"supervisord.process.count"`
	SupervisordProcessUptime MetricSettings `mapstructure:"supervisord.process.uptime"`
}

func DefaultMetricsSettings() MetricsSettings {
	return MetricsSettings{
		SupervisordProcessCount: MetricSettings{
			Enabled: true,
		},
		SupervisordProcessUptime: MetricSettings{
			Enabled: true,
		},
	}
}

// AttributeName specifies the a value name attribute.
type AttributeName int

const (
	_ AttributeName = iota
	AttributeNameNameless
	AttributeNameName
)

// String returns the string representation of the AttributeName.
func (av AttributeName) String() string {
	switch av {
	case AttributeNameNameless:
		return "nameless"
	case AttributeNameName:
		return "name"
	}
	return ""
}

// MapAttributeName is a helper map of string to AttributeName attribute value.
var MapAttributeName = map[string]AttributeName{
	"nameless": AttributeNameNameless,
	"name":     AttributeNameName,
}

type metricSupervisordProcessCount struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills supervisord.process.count metric with initial data.
func (m *metricSupervisordProcessCount) init() {
	m.data.SetName("supervisord.process.count")
	m.data.SetDescription("The number of supervisord monitored processes")
	m.data.SetUnit("{processe}")
	m.data.SetDataType(pmetric.MetricDataTypeGauge)
}

func (m *metricSupervisordProcessCount) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSupervisordProcessCount) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSupervisordProcessCount) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSupervisordProcessCount(settings MetricSettings) metricSupervisordProcessCount {
	m := metricSupervisordProcessCount{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricSupervisordProcessUptime struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills supervisord.process.uptime metric with initial data.
func (m *metricSupervisordProcessUptime) init() {
	m.data.SetName("supervisord.process.uptime")
	m.data.SetDescription("The process uptime")
	m.data.SetUnit("{count}")
	m.data.SetDataType(pmetric.MetricDataTypeGauge)
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSupervisordProcessUptime) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, nameAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().InsertString("name", nameAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSupervisordProcessUptime) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSupervisordProcessUptime) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSupervisordProcessUptime(settings MetricSettings) metricSupervisordProcessUptime {
	m := metricSupervisordProcessUptime{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user settings.
type MetricsBuilder struct {
	startTime                      pcommon.Timestamp   // start time that will be applied to all recorded data points.
	metricsCapacity                int                 // maximum observed number of metrics per resource.
	resourceCapacity               int                 // maximum observed number of resource attributes.
	metricsBuffer                  pmetric.Metrics     // accumulates metrics data before emitting.
	buildInfo                      component.BuildInfo // contains version information
	metricSupervisordProcessCount  metricSupervisordProcessCount
	metricSupervisordProcessUptime metricSupervisordProcessUptime
}

// metricBuilderOption applies changes to default metrics builder.
type metricBuilderOption func(*MetricsBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pcommon.Timestamp) metricBuilderOption {
	return func(mb *MetricsBuilder) {
		mb.startTime = startTime
	}
}

func NewMetricsBuilder(settings MetricsSettings, buildInfo component.BuildInfo, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		startTime:                      pcommon.NewTimestampFromTime(time.Now()),
		metricsBuffer:                  pmetric.NewMetrics(),
		buildInfo:                      buildInfo,
		metricSupervisordProcessCount:  newMetricSupervisordProcessCount(settings.SupervisordProcessCount),
		metricSupervisordProcessUptime: newMetricSupervisordProcessUptime(settings.SupervisordProcessUptime),
	}
	for _, op := range options {
		op(mb)
	}
	return mb
}

// updateCapacity updates max length of metrics and resource attributes that will be used for the slice capacity.
func (mb *MetricsBuilder) updateCapacity(rm pmetric.ResourceMetrics) {
	if mb.metricsCapacity < rm.ScopeMetrics().At(0).Metrics().Len() {
		mb.metricsCapacity = rm.ScopeMetrics().At(0).Metrics().Len()
	}
	if mb.resourceCapacity < rm.Resource().Attributes().Len() {
		mb.resourceCapacity = rm.Resource().Attributes().Len()
	}
}

// ResourceMetricsOption applies changes to provided resource metrics.
type ResourceMetricsOption func(pmetric.ResourceMetrics)

// WithStartTimeOverride overrides start time for all the resource metrics data points.
// This option should be only used if different start time has to be set on metrics coming from different resources.
func WithStartTimeOverride(start pcommon.Timestamp) ResourceMetricsOption {
	return func(rm pmetric.ResourceMetrics) {
		var dps pmetric.NumberDataPointSlice
		metrics := rm.ScopeMetrics().At(0).Metrics()
		for i := 0; i < metrics.Len(); i++ {
			switch metrics.At(i).DataType() {
			case pmetric.MetricDataTypeGauge:
				dps = metrics.At(i).Gauge().DataPoints()
			case pmetric.MetricDataTypeSum:
				dps = metrics.At(i).Sum().DataPoints()
			}
			for j := 0; j < dps.Len(); j++ {
				dps.At(j).SetStartTimestamp(start)
			}
		}
	}
}

// EmitForResource saves all the generated metrics under a new resource and updates the internal state to be ready for
// recording another set of data points as part of another resource. This function can be helpful when one scraper
// needs to emit metrics from several resources. Otherwise calling this function is not required,
// just `Emit` function can be called instead.
// Resource attributes should be provided as ResourceMetricsOption arguments.
func (mb *MetricsBuilder) EmitForResource(rmo ...ResourceMetricsOption) {
	rm := pmetric.NewResourceMetrics()
	rm.Resource().Attributes().EnsureCapacity(mb.resourceCapacity)
	ils := rm.ScopeMetrics().AppendEmpty()
	ils.Scope().SetName("otelcol/hypervisord receiver")
	ils.Scope().SetVersion(mb.buildInfo.Version)
	ils.Metrics().EnsureCapacity(mb.metricsCapacity)
	mb.metricSupervisordProcessCount.emit(ils.Metrics())
	mb.metricSupervisordProcessUptime.emit(ils.Metrics())
	for _, op := range rmo {
		op(rm)
	}
	if ils.Metrics().Len() > 0 {
		mb.updateCapacity(rm)
		rm.MoveTo(mb.metricsBuffer.ResourceMetrics().AppendEmpty())
	}
}

// Emit returns all the metrics accumulated by the metrics builder and updates the internal state to be ready for
// recording another set of metrics. This function will be responsible for applying all the transformations required to
// produce metric representation defined in metadata and user settings, e.g. delta or cumulative.
func (mb *MetricsBuilder) Emit(rmo ...ResourceMetricsOption) pmetric.Metrics {
	mb.EmitForResource(rmo...)
	metrics := pmetric.NewMetrics()
	mb.metricsBuffer.MoveTo(metrics)
	return metrics
}

// RecordSupervisordProcessCountDataPoint adds a data point to supervisord.process.count metric.
func (mb *MetricsBuilder) RecordSupervisordProcessCountDataPoint(ts pcommon.Timestamp, val int64) {
	mb.metricSupervisordProcessCount.recordDataPoint(mb.startTime, ts, val)
}

// RecordSupervisordProcessUptimeDataPoint adds a data point to supervisord.process.uptime metric.
func (mb *MetricsBuilder) RecordSupervisordProcessUptimeDataPoint(ts pcommon.Timestamp, val int64, nameAttributeValue AttributeName) {
	mb.metricSupervisordProcessUptime.recordDataPoint(mb.startTime, ts, val, nameAttributeValue.String())
}

// Reset resets metrics builder to its initial state. It should be used when external metrics source is restarted,
// and metrics builder should update its startTime and reset it's internal state accordingly.
func (mb *MetricsBuilder) Reset(options ...metricBuilderOption) {
	mb.startTime = pcommon.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
}
