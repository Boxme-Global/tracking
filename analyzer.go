package pirsch

import (
	"fmt"
	"time"
)

const (
	byAttributeQuery = `SELECT "%s", count(DISTINCT fingerprint) "visitors", "visitors" / (
			SELECT sum(s)
			FROM (
				SELECT count(DISTINCT fingerprint) "s"
				FROM "hit"
				WHERE %s
				GROUP BY "%s"
			)	
		) "relative_visitors"
		FROM "hit"
		WHERE %s
		GROUP BY "%s"
		ORDER BY "visitors" DESC, "%s" ASC`
)

// Analyzer provides an interface to analyze statistics.
type Analyzer struct {
	store Store
}

// NewAnalyzer returns a new Analyzer for given Store.
func NewAnalyzer(store Store) *Analyzer {
	return &Analyzer{
		store,
	}
}

// ActiveVisitors returns the active visitors per path and the total number of active visitors for given duration.
// Use time.Minute*5 for example to see the active visitors for the past 5 minutes.
// The correct date/time is not included.
func (analyzer *Analyzer) ActiveVisitors(filter *Filter, duration time.Duration) ([]Stats, int, error) {
	filter = analyzer.getFilter(filter)
	filter.Start = time.Now().UTC().Add(-duration)
	args, filterQuery := filter.query()
	query := fmt.Sprintf(`SELECT "path", count(DISTINCT fingerprint) "visitors"
		FROM "hit"
		WHERE %s
		GROUP BY "path"
		ORDER BY "visitors" DESC, "path" ASC`, filterQuery)
	visitors, err := analyzer.store.Select(query, args...)

	if err != nil {
		return nil, 0, err
	}

	query = fmt.Sprintf(`SELECT count(DISTINCT fingerprint) "visitors" FROM "hit" WHERE %s`, filterQuery)
	count, err := analyzer.store.Count(query, args...)

	if err != nil {
		return nil, 0, err
	}

	return visitors, count, nil
}

// Languages returns the visitor count per language.
func (analyzer *Analyzer) Languages(filter *Filter) ([]Stats, error) {
	return analyzer.selectByAttribute(filter, "language")
}

// Countries returns the visitor count per country.
func (analyzer *Analyzer) Countries(filter *Filter) ([]Stats, error) {
	return analyzer.selectByAttribute(filter, "country_code")
}

// Browser returns the visitor count per browser.
func (analyzer *Analyzer) Browser(filter *Filter) ([]Stats, error) {
	return analyzer.selectByAttribute(filter, "browser")
}

// OS returns the visitor count per operating system.
func (analyzer *Analyzer) OS(filter *Filter) ([]Stats, error) {
	return analyzer.selectByAttribute(filter, "os")
}

// Platform returns the visitor count per platform.
func (analyzer *Analyzer) Platform(filter *Filter) (*Stats, error) {
	filterArgs, filterQuery := analyzer.getFilter(filter).query()
	query := fmt.Sprintf(`SELECT (
			SELECT count(DISTINCT fingerprint)
			FROM "hit"
			WHERE %s
			AND desktop = 1
			AND mobile = 0
		) AS "platform_desktop",
		(
			SELECT count(DISTINCT fingerprint)
			FROM "hit"
			WHERE %s
			AND desktop = 0
			AND mobile = 1
		) AS "platform_mobile",
		(
			SELECT count(DISTINCT fingerprint)
			FROM "hit"
			WHERE %s
			AND desktop = 0
			AND mobile = 0
		) AS "platform_unknown",
		"platform_desktop" / IF("platform_desktop" + "platform_mobile" + "platform_unknown" = 0, 1, "platform_desktop" + "platform_mobile" + "platform_unknown") AS relative_platform_desktop,
		"platform_mobile" / IF("platform_desktop" + "platform_mobile" + "platform_unknown" = 0, 1, "platform_desktop" + "platform_mobile" + "platform_unknown") AS relative_platform_mobile,
		"platform_unknown" / IF("platform_desktop" + "platform_mobile" + "platform_unknown" = 0, 1, "platform_desktop" + "platform_mobile" + "platform_unknown") AS relative_platform_unknown`, filterQuery, filterQuery, filterQuery)
	args := make([]interface{}, 0, len(filterArgs)*3)
	args = append(args, filterArgs...)
	args = append(args, filterArgs...)
	args = append(args, filterArgs...)
	return analyzer.store.Get(query, args...)
}

// ScreenClass returns the visitor count per screen class.
func (analyzer *Analyzer) ScreenClass(filter *Filter) ([]Stats, error) {
	return analyzer.selectByAttribute(filter, "screen_class")
}

func (analyzer *Analyzer) selectByAttribute(filter *Filter, attr string) ([]Stats, error) {
	args, filterQuery := analyzer.getFilter(filter).query()
	query := fmt.Sprintf(byAttributeQuery, attr, filterQuery, attr, filterQuery, attr, attr)
	args = append(args, args...)
	return analyzer.store.Select(query, args...)
}

func (analyzer *Analyzer) getFilter(filter *Filter) *Filter {
	if filter == nil {
		return NewFilter(NullTenant)
	}

	filter.validate()
	return filter
}
