// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package version

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/log"
)

// defaultReportingInterval is the default interval at which the version is
// reported.
const defaultReportingInterval = 5 * time.Minute

// ReportingService is a service that periodically logs the running chain
// version.
type ReportingService struct {
	// logger is used to log information about the running chain version.
	logger log.Logger[any]
	// version represents the current version of the running chain.
	version string
	// reportingInterval is the interval at which the version is reported.
	reportingInterval time.Duration
	// metrics contains the metrics for the version service.
	metrics *versionMetrics
}

// NewReportingService creates a new VersionReporterService.
func NewReportingService(
	logger log.Logger[any],
	telemetrySink TelemetrySink,
	version string,
) *ReportingService {
	return &ReportingService{
		logger:            logger,
		version:           version,
		reportingInterval: defaultReportingInterval,
		metrics:           newVersionMetrics(logger, telemetrySink),
	}
}

// Name returns the name of the service.
func (*ReportingService) Name() string {
	return "reporting"
}

// Start begins the periodic logging of the chain version.
func (v *ReportingService) Start(ctx context.Context) error {
	ticker := time.NewTicker(v.reportingInterval)
	v.metrics.reportVersion(v.version)
	go func() {
		for {
			select {
			case <-ticker.C:
				v.metrics.reportVersion(v.version)
				continue
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

// Status returns nil if the service is healthy.
func (*ReportingService) Status() error {
	return nil
}
