package omisocial

import (
	"fmt"
	"strings"
)

var (
	fieldPath = field{
		querySessions:  "path",
		queryPageViews: "path",
		queryDirection: "ASC",
		name:           "path",
	}
	fieldEntryPath = field{
		querySessions:  "entry_path",
		queryPageViews: "entry_path",
		queryDirection: "ASC",
		name:           "entry_path",
	}
	fieldEntries = field{
		querySessions:  "sum(sign)",
		queryPageViews: "uniq(visitor_id, session_id)",
		queryDirection: "DESC",
		name:           "entries",
	}
	fieldExitPath = field{
		querySessions:  "exit_path",
		queryPageViews: "exit_path",
		queryDirection: "ASC",
		name:           "exit_path",
	}
	fieldExits = field{
		querySessions:  "sum(sign)",
		queryPageViews: "uniq(visitor_id, session_id)",
		queryDirection: "DESC",
		name:           "exits",
	}
	fieldVisitors = field{
		querySessions:  "uniq(visitor_id)",
		queryPageViews: "uniq(visitor_id)",
		queryDirection: "DESC",
		name:           "visitors",
	}
	fieldRelativeVisitors = field{
		querySessions:  "visitors / greatest((SELECT uniq(visitor_id) FROM session WHERE %s), 1)",
		queryPageViews: "visitors / greatest((SELECT uniq(visitor_id) FROM session WHERE %s), 1)",
		queryDirection: "DESC",
		filterTime:     true,
		name:           "relative_visitors",
	}
	fieldCR = field{
		querySessions:  "visitors / greatest((SELECT uniq(visitor_id) FROM session WHERE %s), 1)",
		queryPageViews: "visitors / greatest((SELECT uniq(visitor_id) FROM session WHERE %s), 1)",
		queryDirection: "DESC",
		filterTime:     true,
		name:           "cr",
	}
	fieldSessions = field{
		querySessions:  "uniq(visitor_id, session_id)",
		queryPageViews: "uniq(visitor_id, session_id)",
		queryDirection: "DESC",
		name:           "sessions",
	}
	fieldViews = field{
		querySessions:  "sum(page_views*sign)",
		queryPageViews: "count(1)",
		queryDirection: "DESC",
		name:           "views",
	}
	fieldRelativeViews = field{
		querySessions:  "views / greatest((SELECT sum(page_views*sign) views FROM session WHERE %s), 1)",
		queryPageViews: "views / greatest((SELECT sum(page_views*sign) views FROM session WHERE %s), 1)",
		queryDirection: "DESC",
		filterTime:     true,
		name:           "relative_views",
	}
	fieldBounces = field{
		querySessions:  "sum(is_bounce*sign)",
		queryPageViews: "sum(is_bounce)",
		queryDirection: "DESC",
		name:           "bounces",
	}
	fieldBounceRate = field{
		querySessions:  "bounces / IF(sessions = 0, 1, sessions)",
		queryPageViews: "bounces / IF(sessions = 0, 1, sessions)",
		queryDirection: "DESC",
		name:           "bounce_rate",
	}
	fieldReferrer = field{
		querySessions:  "referrer",
		queryPageViews: "referrer",
		queryDirection: "ASC",
		name:           "referrer",
	}
	fieldAnyReferrer = field{
		querySessions:  "any(referrer)",
		queryPageViews: "any(referrer)",
		queryDirection: "ASC",
		name:           "referrer",
	}
	fieldReferrerName = field{
		querySessions:  "referrer_name",
		queryPageViews: "referrer_name",
		queryDirection: "ASC",
		name:           "referrer_name",
	}
	fieldReferrerIcon = field{
		querySessions:  "any(referrer_icon)",
		queryPageViews: "any(referrer_icon)",
		queryDirection: "ASC",
		name:           "referrer_icon",
	}
	fieldLanguage = field{
		querySessions:  "language",
		queryPageViews: "language",
		queryDirection: "ASC",
		name:           "language",
	}
	fieldCountry = field{
		querySessions:  "country_code",
		queryPageViews: "country_code",
		queryDirection: "ASC",
		name:           "country_code",
	}
	fieldCity = field{
		querySessions:  "city",
		queryPageViews: "city",
		queryDirection: "ASC",
		name:           "city",
	}
	fieldBrowser = field{
		querySessions:  "browser",
		queryPageViews: "browser",
		queryDirection: "ASC",
		name:           "browser",
	}
	fieldBrowserVersion = field{
		querySessions:  "browser_version",
		queryPageViews: "browser_version",
		queryDirection: "ASC",
		name:           "browser_verfieldDaysion",
	}
	fieldOS = field{
		querySessions:  "os",
		queryPageViews: "os",
		queryDirection: "ASC",
		name:           "os",
	}
	fieldOSVersion = field{
		querySessions:  "os_version",
		queryPageViews: "os_version",
		queryDirection: "ASC",
		name:           "os_version",
	}
	fieldScreenClass = field{
		querySessions:  "screen_class",
		queryPageViews: "screen_class",
		queryDirection: "ASC",
		name:           "screen_class",
	}
	fieldUTMSource = field{
		querySessions:  "utm_source",
		queryPageViews: "utm_source",
		queryDirection: "ASC",
		name:           "utm_source",
	}
	fieldUTMMedium = field{
		querySessions:  "utm_medium",
		queryPageViews: "utm_medium",
		queryDirection: "ASC",
		name:           "utm_medium",
	}
	fieldUTMCampaign = field{
		querySessions:  "utm_campaign",
		queryPageViews: "utm_campaign",
		queryDirection: "ASC",
		name:           "utm_campaign",
	}
	fieldUTMContent = field{
		querySessions:  "utm_content",
		queryPageViews: "utm_content",
		queryDirection: "ASC",
		name:           "utm_content",
	}
	fieldUTMTerm = field{
		querySessions:  "utm_term",
		queryPageViews: "utm_term",
		queryDirection: "ASC",
		name:           "utm_term",
	}
	fieldOTMSource = field{
		querySessions:  "otm_source",
		queryPageViews: "otm_source",
		queryDirection: "ASC",
		name:           "otm_source",
	}
	fieldTitle = field{
		querySessions:  "title",
		queryPageViews: "title",
		queryDirection: "ASC",
		name:           "title",
	}
	fieldEntryTitle = field{
		querySessions:  "entry_title",
		queryPageViews: "entry_title",
		queryDirection: "ASC",
		name:           "title",
	}
	fieldExitTitle = field{
		querySessions:  "exit_title",
		queryPageViews: "exit_title",
		queryDirection: "ASC",
		name:           "title",
	}
	fieldHour = field{
		querySessions:  "toHour(time, '%s')",
		queryPageViews: "toHour(time, '%s')",
		queryDirection: "ASC",
		queryWithFill:  "WITH FILL FROM toHour(toDate(?), '%s') TO toHour(toDate(?), '%s')",
		timezone:       true,
		name:           "hour",
	}
	fieldDay = field{
		querySessions:  "toDate(time, '%s')",
		queryPageViews: "toDate(time, '%s')",
		queryDirection: "ASC",
		withFill:       true,
		timezone:       true,
		name:           "period",
	}
	fieldWeek = field{
		querySessions:  "toWeek(time, 9, '%s')",
		queryPageViews: "toWeek(time, 9, '%s')",
		queryDirection: "ASC",
		queryWithFill:  "WITH FILL FROM toWeek(toDate(?), 9, '%s') TO toWeek(toDate(?), 9, '%s')",
		timezone:       true,
		name:           "period",
	}
	fieldMonth = field{
		querySessions:  "toMonth(time, '%s')",
		queryPageViews: "toMonth(time, '%s')",
		queryDirection: "ASC",
		queryWithFill:  "WITH FILL FROM toMonth(toDate(?), '%s') TO toMonth(toDate(?), '%s')",
		timezone:       true,
		name:           "period",
	}
	fieldEventTimeSpent = field{
		querySessions:  "ifNull(toUInt64(avg(nullIf(duration_seconds, 0))), 0)",
		queryPageViews: "ifNull(toUInt64(avg(nullIf(duration_seconds, 0))), 0)",
		name:           "average_time_spent_seconds",
	}
	fieldDesktop = field{
		querySessions:  "desktop",
		queryPageViews: "desktop",
		queryDirection: "DESC",
		name:           "desktop",
	}
	fieldMobile = field{
		querySessions:  "mobile",
		queryPageViews: "mobile",
		queryDirection: "DESC",
		name:           "mobile",
	}
)

type field struct {
	querySessions  string
	queryPageViews string
	queryDirection string
	queryWithFill  string
	withFill       bool
	timezone       bool
	filterTime     bool
	name           string
}

func buildQuery(filter *Filter, fields, groupBy, orderBy []field) ([]interface{}, string) {
	table := filter.table()
	args := make([]interface{}, 0)
	var query strings.Builder

	if table == "event" || filter.Path != "" || filter.PathPattern != "" || fieldsContain(fields, fieldPath.name) {
		if table == "session" {
			table = "page_view"
		}

		query.WriteString(fmt.Sprintf(`SELECT %s FROM %s v `, joinPageViewFields(&args, filter, fields), table))

		if filter.EntryPath != "" ||
			filter.ExitPath != "" ||
			fieldsContain(fields, fieldBounces.name) ||
			fieldsContain(fields, fieldViews.name) ||
			fieldsContain(fields, fieldEntryPath.name) ||
			fieldsContain(fields, fieldExitPath.name) {
			path, pathPattern, eventName := filter.Path, filter.PathPattern, filter.EventName
			filter.Path, filter.PathPattern, filter.EventName = "", "", ""
			filterArgs, filterQuery := filter.query()
			filter.Path, filter.PathPattern, filter.EventName = path, pathPattern, eventName
			args = append(args, filterArgs...)

			if table == "page_view" {
				query.WriteString("INNER ")
			} else {
				query.WriteString("LEFT ")
			}

			sessionFields := make([]string, 0, 4)

			if fieldsContain(fields, fieldEntryPath.name) {
				sessionFields = append(sessionFields, fieldEntryPath.name)
			}

			if fieldsContain(fields, fieldExitPath.name) {
				sessionFields = append(sessionFields, fieldExitPath.name)
			}

			if fieldsContain(fields, fieldBounces.name) {
				sessionFields = append(sessionFields, "sum(is_bounce*sign) is_bounce")
			}

			if fieldsContain(fields, fieldViews.name) {
				sessionFields = append(sessionFields, "sum(page_views*sign) page_views")
			}

			sessionFieldsQuery := strings.Join(sessionFields, ",")

			if sessionFieldsQuery != "" {
				sessionFieldsQuery = "," + sessionFieldsQuery
			}

			query.WriteString(fmt.Sprintf(`JOIN (
				SELECT visitor_id,
				session_id
				%s
				FROM session
				WHERE %s
				GROUP BY visitor_id, session_id, entry_path, exit_path
				HAVING sum(sign) > 0
			) s
			ON s.visitor_id = v.visitor_id AND s.session_id = v.session_id `, sessionFieldsQuery, filterQuery))

			if filter.EventName != "" {
				filterArgs, filterQuery = filter.query()
				args = append(args, filterArgs...)
				query.WriteString(fmt.Sprintf(`WHERE %s `, filterQuery))
			} else if filter.Path != "" || filter.PathPattern != "" {
				filterArgs, filterQuery = filter.queryPageOrEvent()
				args = append(args, filterArgs...)
				query.WriteString(fmt.Sprintf(`WHERE %s `, filterQuery))
			}
		} else {
			filterArgs, filterQuery := filter.query()
			args = append(args, filterArgs...)
			query.WriteString(fmt.Sprintf(`WHERE %s `, filterQuery))
		}

		if len(groupBy) > 0 {
			query.WriteString(fmt.Sprintf(`GROUP BY %s `, joinGroupBy(groupBy)))
		}

		if len(orderBy) > 0 {
			query.WriteString(fmt.Sprintf(`ORDER BY %s `, joinOrderBy(&args, filter, orderBy)))
		}
	} else {
		filterArgs, filterQuery := filter.query()
		query.WriteString(fmt.Sprintf(`SELECT %s
			FROM session
			WHERE %s `, joinSessionFields(&args, filter, fields), filterQuery))
		args = append(args, filterArgs...)

		if len(groupBy) > 0 {
			query.WriteString(fmt.Sprintf(`GROUP BY %s `, joinGroupBy(groupBy)))
		}

		query.WriteString(`HAVING sum(sign) > 0 `)

		if len(orderBy) > 0 {
			query.WriteString(fmt.Sprintf(`ORDER BY %s `, joinOrderBy(&args, filter, orderBy)))
		}
	}

	query.WriteString(filter.withLimit())
	query.WriteString(filter.withOffset())
	return args, query.String()
}

func joinPageViewFields(args *[]interface{}, filter *Filter, fields []field) string {
	if len(fields) == 0 {
		return ""
	}

	var out strings.Builder

	for i := range fields {
		if fields[i].filterTime {
			timeArgs, timeQuery := filter.queryTime()
			*args = append(*args, timeArgs...)
			out.WriteString(fmt.Sprintf(`%s %s,`, fmt.Sprintf(fields[i].queryPageViews, timeQuery), fields[i].name))
		} else if fields[i].timezone {
			out.WriteString(fmt.Sprintf(`%s %s,`, fmt.Sprintf(fields[i].queryPageViews, filter.Timezone.String()), fields[i].name))
		} else {
			out.WriteString(fmt.Sprintf(`%s %s,`, fields[i].queryPageViews, fields[i].name))
		}
	}

	str := out.String()
	return str[:len(str)-1]
}

func joinSessionFields(args *[]interface{}, filter *Filter, fields []field) string {
	if len(fields) == 0 {
		return ""
	}

	var out strings.Builder

	for i := range fields {
		if fields[i].filterTime {
			timeArgs, timeQuery := filter.queryTime()
			*args = append(*args, timeArgs...)
			out.WriteString(fmt.Sprintf(`%s %s,`, fmt.Sprintf(fields[i].queryPageViews, timeQuery), fields[i].name))
		} else if fields[i].timezone {
			out.WriteString(fmt.Sprintf(`%s %s,`, fmt.Sprintf(fields[i].querySessions, filter.Timezone.String()), fields[i].name))
		} else {
			out.WriteString(fmt.Sprintf(`%s %s,`, fields[i].querySessions, fields[i].name))
		}
	}

	str := out.String()
	return str[:len(str)-1]
}

func joinGroupBy(fields []field) string {
	if len(fields) == 0 {
		return ""
	}

	var out strings.Builder

	for i := range fields {
		out.WriteString(fields[i].name + ",")
	}

	str := out.String()
	return str[:len(str)-1]
}

func joinOrderBy(args *[]interface{}, filter *Filter, fields []field) string {
	if len(fields) == 0 {
		return ""
	}

	var out strings.Builder

	for i := range fields {
		if fields[i].queryWithFill != "" {
			queryFill := fields[i].queryWithFill
			if fields[i].timezone {
				queryFill = fmt.Sprintf(fields[i].queryWithFill, filter.Timezone.String(), filter.Timezone.String())
			}
			fillArgs := []interface{}{filter.From, filter.To}
			*args = append(*args, fillArgs...)
			out.WriteString(fmt.Sprintf(`%s %s %s,`, fields[i].name, fields[i].queryDirection, queryFill))
		} else if fields[i].withFill {
			fillArgs, fillQuery := filter.withFill()
			*args = append(*args, fillArgs...)
			out.WriteString(fmt.Sprintf(`%s %s %s,`, fields[i].name, fields[i].queryDirection, fillQuery))
		} else {
			out.WriteString(fmt.Sprintf(`%s %s,`, fields[i].name, fields[i].queryDirection))
		}
	}

	str := out.String()
	return str[:len(str)-1]
}

func fieldsContain(haystack []field, needle string) bool {
	for i := range haystack {
		if haystack[i].name == needle {
			return true
		}
	}

	return false
}
