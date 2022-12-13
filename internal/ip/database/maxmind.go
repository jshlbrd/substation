package database

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/brexhq/substation/internal/file"
	"github.com/brexhq/substation/internal/ip"
	"github.com/oschwald/geoip2-golang"
)

// MaxMindASN provides read access to MaxMind ASN database.
type MaxMindASN struct {
	Database string `json:"database"`
	Language string `json:"language"`
	db       *geoip2.Reader
}

// IsEnabled returns true if the database is open and ready for use.
func (d *MaxMindASN) IsEnabled() bool {
	return d.db != nil
}

// Open retrieves the database and opens it for querying. The location of the database can be either a path on local disk, an HTTP(S) URL, or an AWS S3 URL. MaxMind language support is provided by calling GetMaxMindLanguage to retrieve a user-configured language.
func (d *MaxMindASN) Open(ctx context.Context) error {
	// language defaults to English
	if d.Language == "" {
		d.Language = "en"
	}

	path, err := file.Get(ctx, d.Database)
	defer os.Remove(path)

	if err != nil {
		return fmt.Errorf("database: %v", err)
	}

	if d.db, err = geoip2.Open(path); err != nil {
		return fmt.Errorf("database: %v", err)
	}

	return nil
}

// Close closes the open database.
func (d *MaxMindASN) Close() error {
	if d.IsEnabled() {
		if err := d.db.Close(); err != nil {
			return fmt.Errorf("database: %v", err)
		}
	}

	return nil
}

// Get queries the database and returns an aggregated database record containing enrichment information.
func (d *MaxMindASN) Get(addr string) (*ip.EnrichmentRecord, error) {
	paddr := net.ParseIP(addr)
	if paddr == nil {
		return nil, fmt.Errorf("database: %v", ip.ErrInvalidIPAddress)
	}

	resp, err := d.db.ASN(paddr)
	if err != nil {
		return nil, err
	}

	rec := &ip.EnrichmentRecord{
		ASN: &ip.ASN{
			Number:       resp.AutonomousSystemNumber,
			Organization: resp.AutonomousSystemOrganization,
		},
	}

	return rec, nil
}

// MaxMindCity provides read access to a MaxMind City database.
type MaxMindCity struct {
	Database string `json:"database"`
	Language string `json:"language"`
	db       *geoip2.Reader
}

// IsEnabled returns true if the database is open and ready for use.
func (d *MaxMindCity) IsEnabled() bool {
	return d.db != nil
}

// Open retrieves the database and opens it for querying. MaxMind language support is provided by calling GetMaxMindLanguage to retrieve a user-configured language.
func (d *MaxMindCity) Open(ctx context.Context) error {
	// language defaults to English
	if d.Language == "" {
		d.Language = "en"
	}

	path, err := file.Get(ctx, d.Database)
	defer os.Remove(path)

	if err != nil {
		return fmt.Errorf("database: %v", err)
	}

	if d.db, err = geoip2.Open(path); err != nil {
		return fmt.Errorf("database: %v", err)
	}

	return nil
}

// Close closes the open database.
func (d *MaxMindCity) Close() error {
	if d.IsEnabled() {
		if err := d.db.Close(); err != nil {
			return fmt.Errorf("database: %v", err)
		}
	}

	return nil
}

// Get queries the database and returns an aggregated database record containing enrichment information.
func (d *MaxMindCity) Get(addr string) (*ip.EnrichmentRecord, error) {
	paddr := net.ParseIP(addr)
	if paddr == nil {
		return nil, fmt.Errorf("database: %v", ip.ErrInvalidIPAddress)
	}

	resp, err := d.db.City(paddr)
	if err != nil {
		return nil, fmt.Errorf("database: %v", err)
	}

	rec := &ip.EnrichmentRecord{
		Location: &ip.Location{
			Coordinates: &ip.Coordinates{
				Latitude:  float32(resp.Location.Latitude),
				Longitude: float32(resp.Location.Longitude),
			},
			Continent:  resp.Continent.Names[d.Language],
			Country:    resp.Country.Names[d.Language],
			City:       resp.City.Names[d.Language],
			PostalCode: resp.Postal.Code,
			Accuracy:   float32(resp.Location.AccuracyRadius),
			TimeZone:   resp.Location.TimeZone,
		},
	}

	return rec, nil
}