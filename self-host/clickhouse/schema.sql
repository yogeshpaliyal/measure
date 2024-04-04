
--
-- Database schema
--

CREATE DATABASE IF NOT EXISTS default;

CREATE TABLE default.events
(
    `id` UUID COMMENT 'unique event id',
    `type` LowCardinality(FixedString(32)) COMMENT 'type of the event',
    `session_id` UUID COMMENT 'associated session id',
    `app_id` UUID COMMENT 'associated app id',
    `inet.ipv4` Nullable(IPv4) COMMENT 'ipv4 address',
    `inet.ipv6` Nullable(IPv6) COMMENT 'ipv6 address',
    `inet.country_code` FixedString(8) COMMENT 'country code',
    `timestamp` DateTime64(9, 'UTC') COMMENT 'event timestamp',
    `attribute.installation_id` UUID COMMENT 'unique id for an installation of an app, generated by sdk',
    `attribute.app_version` FixedString(32) COMMENT 'app version identifier',
    `attribute.app_build` FixedString(32) COMMENT 'app build identifier',
    `attribute.app_unique_id` FixedString(128) COMMENT 'app bundle identifier',
    `attribute.platform` LowCardinality(FixedString(32)) COMMENT 'platform identifier',
    `attribute.measure_sdk_version` FixedString(16) COMMENT 'measure sdk version identifier',
    `attribute.thread_name` FixedString(64) COMMENT 'thread on which the event was captured',
    `attribute.user_id` FixedString(128) COMMENT 'id of the app\'s end user',
    `attribute.device_name` FixedString(32) COMMENT 'name of the device',
    `attribute.device_model` FixedString(32) COMMENT 'model of the device',
    `attribute.device_manufacturer` FixedString(32) COMMENT 'manufacturer of the device',
    `attribute.device_type` LowCardinality(FixedString(32)) COMMENT 'type of the device, like phone or tablet',
    `attribute.device_is_foldable` Bool COMMENT 'true for foldable devices',
    `attribute.device_is_physical` Bool COMMENT 'true for physical devices',
    `attribute.device_density_dpi` UInt16 COMMENT 'dpi density',
    `attribute.device_width_px` UInt16 COMMENT 'screen width',
    `attribute.device_height_px` UInt16 COMMENT 'screen height',
    `attribute.device_density` Float32 COMMENT 'device density',
    `attribute.device_locale` FixedString(64) COMMENT 'rfc 5646 locale string',
    `attribute.os_name` FixedString(32) COMMENT 'name of the operating system',
    `attribute.os_version` FixedString(32) COMMENT 'version of the operating system',
    `attribute.network_type` LowCardinality(FixedString(16)) COMMENT 'either - wifi, cellular, vpn, unknown, no_network',
    `attribute.network_generation` LowCardinality(FixedString(8)) COMMENT 'either - 2g, 3g, 4g, 5g',
    `attribute.network_provider` FixedString(64) COMMENT 'name of the network service provider',
    `anr.handled` Bool COMMENT 'anr was handled by the application code',
    `anr.fingerprint` FixedString(16) COMMENT 'fingerprint for anr similarity classification',
    `anr.exceptions` String COMMENT 'anr exception data',
    `anr.threads` String COMMENT 'anr thread data',
    `anr.foreground` Bool COMMENT 'true if the anr was perceived by end user',
    `exception.handled` Bool COMMENT 'exception was handled by application code',
    `exception.fingerprint` FixedString(16) COMMENT 'fingerprint for exception similarity classification',
    `exception.exceptions` String COMMENT 'exception data',
    `exception.threads` String COMMENT 'exception thread data',
    `exception.foreground` Bool COMMENT 'true if the exception was perceived by end user',
    `app_exit.reason` LowCardinality(FixedString(64)) COMMENT 'reason for app exit',
    `app_exit.importance` LowCardinality(FixedString(32)) COMMENT 'importance of process that it used to have before death',
    `app_exit.trace` String COMMENT 'modified trace given by ApplicationExitInfo to help debug anrs.',
    `app_exit.process_name` String COMMENT 'name of the process that died',
    `app_exit.pid` String COMMENT 'id of the process that died',
    `string.severity_text` LowCardinality(FixedString(10)) COMMENT 'log level - info, warning, error, fatal, debug',
    `string.string` String COMMENT 'log message text',
    `gesture_long_click.target` FixedString(128) COMMENT 'class or instance name of the originating view',
    `gesture_long_click.target_id` FixedString(128) COMMENT 'unique identifier of the target',
    `gesture_long_click.touch_down_time` UInt32 COMMENT 'time for touch down gesture',
    `gesture_long_click.touch_up_time` UInt32 COMMENT 'time for touch up gesture',
    `gesture_long_click.width` UInt16 COMMENT 'width of the target view in pixels',
    `gesture_long_click.height` UInt16 COMMENT 'height of the target view in pixels',
    `gesture_long_click.x` Float32 COMMENT 'x coordinate of where the gesture happened',
    `gesture_long_click.y` Float32 COMMENT 'y coordinate of where the gesture happened',
    `gesture_click.target` FixedString(128) COMMENT 'class or instance name of the originating view',
    `gesture_click.target_id` FixedString(128) COMMENT 'unique identifier of the target',
    `gesture_click.touch_down_time` UInt32 COMMENT 'time for touch down gesture',
    `gesture_click.touch_up_time` UInt32 COMMENT 'time for the touch up gesture',
    `gesture_click.width` UInt16 COMMENT 'width of the target view in pixels',
    `gesture_click.height` UInt16 COMMENT 'height of the target view in pixels',
    `gesture_click.x` Float32 COMMENT 'x coordinate of where the gesture happened',
    `gesture_click.y` Float32 COMMENT 'y coordinate of where the gesture happened',
    `gesture_scroll.target` FixedString(128) COMMENT 'class or instance name of the originating view',
    `gesture_scroll.target_id` FixedString(128) COMMENT 'unique identifier of the target',
    `gesture_scroll.touch_down_time` UInt32 COMMENT 'time for touch down gesture',
    `gesture_scroll.touch_up_time` UInt32 COMMENT 'time for touch up gesture',
    `gesture_scroll.x` Float32 COMMENT 'x coordinate of where the gesture started',
    `gesture_scroll.y` Float32 COMMENT 'y coordinate of where the gesture started',
    `gesture_scroll.end_x` Float32 COMMENT 'x coordinate of where the gesture ended',
    `gesture_scroll.end_y` Float32 COMMENT 'y coordinate of where the gesture ended',
    `gesture_scroll.direction` FixedString(8) COMMENT 'direction of the scroll',
    `lifecycle_activity.type` FixedString(32) COMMENT 'type of the lifecycle activity, either - created, resumed, paused, destroyed',
    `lifecycle_activity.class_name` FixedString(128) COMMENT 'fully qualified class name of the activity',
    `lifecycle_activity.intent` String COMMENT 'intent data serialized as string',
    `lifecycle_activity.saved_instance_state` Bool COMMENT 'represents that activity was recreated with a saved state. only available for type created.',
    `lifecycle_fragment.type` FixedString(32) COMMENT 'type of the lifecycle fragment, either - attached, resumed, paused, detached',
    `lifecycle_fragment.class_name` FixedString(128) COMMENT 'fully qualified class name of the fragment',
    `lifecycle_fragment.parent_activity` String COMMENT 'fully qualified class name of the parent activity that the fragment is attached to',
    `lifecycle_fragment.tag` String COMMENT 'optional fragment tag',
    `lifecycle_app.type` FixedString(32) COMMENT 'type of the lifecycle app, either - background, foreground',
    `cold_launch.process_start_uptime` UInt32 COMMENT 'start uptime in msec',
    `cold_launch.process_start_requested_uptime` UInt32 COMMENT 'start uptime in msec',
    `cold_launch.content_provider_attach_uptime` UInt32 COMMENT 'start uptime in msec',
    `cold_launch.on_next_draw_uptime` UInt32 COMMENT 'time at which app became visible',
    `cold_launch.launched_activity` FixedString(128) COMMENT 'activity which drew the first frame during cold launch',
    `cold_launch.has_saved_state` Bool COMMENT 'whether the launched_activity was created with a saved state bundle',
    `cold_launch.intent_data` String COMMENT 'intent data used to launch the launched_activity',
    `cold_launch.duration` UInt32 COMMENT 'computed cold launch duration',
    `warm_launch.app_visible_uptime` UInt32 COMMENT 'time since the app became visible to user, in msec',
    `warm_launch.on_next_draw_uptime` UInt32 COMMENT 'time at which app became visible to user, in msec',
    `warm_launch.launched_activity` FixedString(128) COMMENT 'activity which drew the first frame during warm launch',
    `warm_launch.has_saved_state` Bool COMMENT 'whether the launched_activity was created with a saved state bundle',
    `warm_launch.intent_data` String COMMENT 'intent data used to launch the launched_activity',
    `warm_launch.duration` UInt32 COMMENT 'computed warm launch duration',
    `hot_launch.app_visible_uptime` UInt32 COMMENT 'time elapsed since the app became visible to user, in msec',
    `hot_launch.on_next_draw_uptime` UInt32 COMMENT 'time at which app became visible to user, in msec',
    `hot_launch.launched_activity` FixedString(128) COMMENT 'activity which drew the first frame during hot launch',
    `hot_launch.has_saved_state` Bool COMMENT 'whether the launched_activity was created with a saved state bundle',
    `hot_launch.intent_data` String COMMENT 'intent data used to launch the launched_activity',
    `hot_launch.duration` UInt32 COMMENT 'computed hot launch duration',
    `network_change.network_type` LowCardinality(FixedString(16)) COMMENT 'type of the network, wifi, cellular etc',
    `network_change.previous_network_type` LowCardinality(FixedString(16)) COMMENT 'type of the previous network',
    `network_change.network_generation` LowCardinality(FixedString(8)) COMMENT '2g, 3g, 4g etc',
    `network_change.previous_network_generation` LowCardinality(FixedString(8)) COMMENT 'previous network generation',
    `network_change.network_provider` FixedString(64) COMMENT 'name of the network service provider',
    `http.url` String COMMENT 'url of the http request',
    `http.method` LowCardinality(FixedString(16)) COMMENT 'method like get, post',
    `http.status_code` UInt16 COMMENT 'http status code',
    `http.start_time` UInt64 COMMENT 'uptime at when the http call started, in msec',
    `http.end_time` UInt64 COMMENT 'uptime at when the http call ended, in msec',
    `http_request_headers` Map(String, String) COMMENT 'http request headers',
    `http_response_headers` Map(String, String) COMMENT 'http response headers',
    `http.request_body` String COMMENT 'request body',
    `http.response_body` String COMMENT 'response body',
    `http.failure_reason` String COMMENT 'reason for failure',
    `http.failure_description` String COMMENT 'description of the failure',
    `http.client` LowCardinality(FixedString(32)) COMMENT 'name of the http client',
    `memory_usage.java_max_heap` UInt64 COMMENT 'maximum size of the java heap allocated, in kb',
    `memory_usage.java_total_heap` UInt64 COMMENT 'total size of the java heap available for allocation, in KB',
    `memory_usage.java_free_heap` UInt64 COMMENT 'free memory available in the java heap, in kb',
    `memory_usage.total_pss` UInt64 COMMENT 'total proportional set size - amount of memory used by the process, including shared memory and code. in kb.',
    `memory_usage.rss` UInt64 COMMENT 'resident set size - amount of physical memory currently used, in kb',
    `memory_usage.native_total_heap` UInt64 COMMENT 'total size of the native heap (memory out of java\'s control) available for allocation, in kb',
    `memory_usage.native_free_heap` UInt64 COMMENT 'amount of free memory available in the native heap, in kb',
    `memory_usage.interval_config` UInt32 COMMENT 'interval between two consecutive readings, in msec',
    `low_memory.java_max_heap` UInt64 COMMENT 'maximum size of the java heap allocated, in kb',
    `low_memory.java_total_heap` UInt64 COMMENT 'total size of the java heap available for allocation, in kb',
    `low_memory.java_free_heap` UInt64 COMMENT 'free memory available in the java heap, in kb',
    `low_memory.total_pss` UInt64 COMMENT 'total proportional set size - amount of memory used by the process, including shared memory and code. in kb.',
    `low_memory.rss` UInt64 COMMENT 'resident set size - amount of physical memory currently used, in kb',
    `low_memory.native_total_heap` UInt64 COMMENT 'total size of the native heap (memory out of java',
    `low_memory.native_free_heap` UInt64 COMMENT 'amount of free memory available in the native heap, in kb',
    `trim_memory.level` LowCardinality(FixedString(64)) COMMENT 'one of the trim memory constants as received by component callback',
    `cpu_usage.num_cores` UInt8 COMMENT 'number of cores on the device',
    `cpu_usage.clock_speed` UInt32 COMMENT 'clock speed of the processor, in hz',
    `cpu_usage.uptime` UInt64 COMMENT 'time since the device booted, in msec',
    `cpu_usage.utime` UInt64 COMMENT 'execution time in user mode, in jiffies',
    `cpu_usage.cutime` UInt64 COMMENT 'execution time in user mode with child processes, in jiffies',
    `cpu_usage.stime` UInt64 COMMENT 'execution time in kernel mode, in jiffies',
    `cpu_usage.cstime` UInt64 COMMENT 'execution time in user mode with child processes, in jiffies',
    `cpu_usage.interval_config` UInt32 COMMENT 'interval between two consecutive readings, in msec',
    `navigation.route` FixedString(128) COMMENT 'the destination route'
)
ENGINE = MergeTree
PRIMARY KEY (id, app_id, timestamp)
ORDER BY (id, app_id, timestamp)
SETTINGS index_granularity = 8192
COMMENT 'events master table';

CREATE TABLE default.schema_migrations
(
    `version` String,
    `ts` DateTime DEFAULT now(),
    `applied` UInt8 DEFAULT 1
)
ENGINE = ReplacingMergeTree(ts)
PRIMARY KEY version
ORDER BY version
SETTINGS index_granularity = 8192;


--
-- Dbmate schema migrations
--

INSERT INTO schema_migrations (version) VALUES
    ('20231117020810');
