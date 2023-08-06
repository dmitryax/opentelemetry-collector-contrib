// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package sshcheckreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sshcheckreceiver"

import (
	"context"
	"errors"
	"runtime"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sshcheckreceiver/internal/configssh"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sshcheckreceiver/internal/metadata"
)

var errClientNotInit = errors.New("client not initialized")

type sshcheckScraper struct {
	*configssh.Client
	*Config
	settings component.TelemetrySettings
	mb       *metadata.MetricsBuilder
}

// start starts the scraper by creating a new SSH Client on the scraper
func (s *sshcheckScraper) start(_ context.Context, host component.Host) error {
	var err error
	if !supportedOS() {
		return errWindowsUnsupported
	}
	s.Client, err = s.Config.ToClient(host, s.settings)
	return err
}

func (s *sshcheckScraper) scrapeSSH(rmb *metadata.ResourceMetricsBuilder, now pcommon.Timestamp) error {
	var success int64

	start := time.Now()
	err := s.Client.Dial(s.Config.SSHClientSettings.Endpoint)
	if err == nil {
		success = 1
	}
	rmb.RecordSshcheckDurationDataPoint(now, time.Since(start).Milliseconds())
	rmb.RecordSshcheckStatusDataPoint(now, success)
	return err
}

func (s *sshcheckScraper) scrapeSFTP(rmb *metadata.ResourceMetricsBuilder, now pcommon.Timestamp) error {
	var success int64

	start := time.Now()
	// upgrade to SFTP and read fs
	sftpc, err := s.Client.SFTPClient()
	if err == nil {
		_, err = sftpc.ReadDir(".")
		if err == nil {
			success = 1
		}
	}
	rmb.RecordSshcheckSftpDurationDataPoint(now, time.Since(start).Milliseconds())
	rmb.RecordSshcheckSftpStatusDataPoint(now, success)
	return err
}

// timeout chooses the shorter between between a given deadline and timeout
func timeout(deadline time.Time, timeout time.Duration) time.Duration {
	timeToDeadline := time.Until(deadline)
	if timeToDeadline < timeout {
		return timeToDeadline
	}
	return timeout
}

// scrape connects to the endpoint and produces metrics based on the response. TBH the flow-of-control
// is a bit awkward here, because the SFTP checks are not enabled by default and they would panic on nil
// ref to the underlying Conn when SSH checks failed.
func (s *sshcheckScraper) scrape(ctx context.Context) (_ pmetric.Metrics, err error) {
	var (
		to time.Duration
	)
	// check cancellation
	select {
	case <-ctx.Done():
		return pmetric.NewMetrics(), ctx.Err()
	default:
	}

	cleanup := func() {
		s.Client.Close()
	}

	// if the context carries a shorter deadline then timeout that quickly
	deadline, ok := ctx.Deadline()
	if ok {
		to = timeout(deadline, s.Client.Timeout)
		s.Client.Timeout = to
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	now := pcommon.NewTimestampFromTime(time.Now())
	if s.Client == nil {
		return pmetric.NewMetrics(), errClientNotInit
	}

	rb := s.mb.NewResourceBuilder()
	rb.SetSSHEndpoint(s.Config.SSHClientSettings.Endpoint)
	rmb := s.mb.ResourceMetricsBuilder(rb.Emit())

	if err = s.scrapeSSH(rmb, now); err != nil {
		rmb.RecordSshcheckErrorDataPoint(now, int64(1), err.Error())
	} else {
		go func() {
			<-ctx.Done()
			cleanup()
		}()
	}

	if s.SFTPEnabled() {
		if err := s.scrapeSFTP(rmb, now); err != nil {
			rmb.RecordSshcheckSftpErrorDataPoint(now, int64(1), err.Error())
		}
	}

	return s.mb.Emit(), nil
}

func newScraper(conf *Config, settings receiver.CreateSettings) *sshcheckScraper {
	return &sshcheckScraper{
		Config:   conf,
		settings: settings.TelemetrySettings,
		mb:       metadata.NewMetricsBuilder(conf.MetricsBuilderConfig, settings),
	}
}

func supportedOS() bool {
	return !(runtime.GOOS == "windows")
}
