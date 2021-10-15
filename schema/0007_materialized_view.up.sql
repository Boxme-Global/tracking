CREATE MATERIALIZED VIEW events
ENGINE = AggregatingMergeTree()
PARTITION BY toYYYYMM(time)
ORDER BY (time, fingerprint, session_id, path)
POPULATE AS
SELECT client_id,
       fingerprint,
       session_id,
       event_name,
       event_meta_keys,
       event_meta_values,
       arrayJoin(paths) path,
       datetime time,
    duration_seconds,
    event_duration_seconds,
    entry_path,
    exit_path,
    page_views,
    is_bounce,
    language,
    country_code,
    city,
    referrer,
    referrer_name,
    referrer_icon,
    os,
    os_version,
    browser,
    browser_version,
    desktop,
    mobile,
    screen_width,
    screen_height,
    screen_class,
    utm_source,
    utm_medium,
    utm_campaign,
    utm_content,
    utm_term
FROM (
    SELECT client_id,
    fingerprint,
    session_id,
    event_name,
    event_meta_keys,
    event_meta_values,
    groupArray(path) paths,
    max(time) datetime,
    sum(duration_seconds) duration_seconds,
    sum(event_duration_seconds) event_duration_seconds,
    argMin(path, time) entry_path,
    argMax(path, time) exit_path,
    argMax(page_views, time) page_views,
    argMax(is_bounce, time) is_bounce,
    argMax(language, time) language,
    argMax(country_code, time) country_code,
    argMax(city, time) city,
    argMax(referrer, time) referrer,
    argMax(referrer_name, time) referrer_name,
    argMax(referrer_icon, time) referrer_icon,
    argMax(os, time) os,
    argMax(os_version, time) os_version,
    argMax(browser, time) browser,
    argMax(browser_version, time) browser_version,
    argMax(desktop, time) desktop,
    argMax(mobile, time) mobile,
    argMax(screen_width, time) screen_width,
    argMax(screen_height, time) screen_height,
    argMax(screen_class, time) screen_class,
    argMax(utm_source, time) utm_source,
    argMax(utm_medium, time) utm_medium,
    argMax(utm_campaign, time) utm_campaign,
    argMax(utm_content, time) utm_content,
    argMax(utm_term, time) utm_term
    FROM event
    GROUP BY client_id, fingerprint, session_id, event_name, event_meta_keys, event_meta_values
    )
GROUP BY client_id, fingerprint, session_id, event_name, event_meta_keys, event_meta_values, path, time,
    duration_seconds, event_duration_seconds, entry_path, exit_path, page_views,
    is_bounce, language, country_code, city, referrer, referrer_name, referrer_icon,
    os, os_version, browser, browser_version, desktop, mobile, screen_width, screen_height, screen_class,
    utm_source, utm_medium, utm_campaign, utm_content, utm_term
;
