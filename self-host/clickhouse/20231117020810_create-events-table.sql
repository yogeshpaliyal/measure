-- migrate:up
create table if not exists default.events
(
    `id` UUID,
    `type` LowCardinality(FixedString(32)),
    `session_id` UUID,
    `app_id` UUID,
    `inet.ipv4` Nullable(IPv4),
    `inet.ipv6` Nullable(IPv6),
    `inet.country_code` FixedString(8),
    `timestamp` DateTime64(9, 'UTC'),
    `thread_name` FixedString(64),
    `resource.device_name` FixedString(32),
    `resource.device_model` FixedString(32),
    `resource.device_manufacturer` FixedString(32),
    `resource.device_type` LowCardinality(FixedString(32)),
    `resource.device_is_foldable` Bool,
    `resource.device_is_physical` Bool,
    `resource.device_density_dpi` UInt16,
    `resource.device_width_px` UInt16,
    `resource.device_height_px` UInt16,
    `resource.device_density` Float32,
    `resource.os_name` FixedString(32),
    `resource.os_version` FixedString(32),
    `resource.platform` LowCardinality(FixedString(32)),
    `resource.app_version` FixedString(32),
    `resource.app_build` FixedString(32),
    `resource.app_unique_id` FixedString(128),
    `resource.measure_sdk_version` FixedString(16),
    `anr.thread_name` FixedString(64),
    `anr.handled` Bool,
    `anr.fingerprint` FixedString(16),
    `anr_exceptions` Array(Tuple(LowCardinality(String), LowCardinality(String), Array(Tuple(Int32, Int32, LowCardinality(String), LowCardinality(String), LowCardinality(String), LowCardinality(String))))),
    `anr_threads` Array(Tuple(LowCardinality(String), Array(Tuple(Int32, Int32, LowCardinality(String), LowCardinality(String), LowCardinality(String), LowCardinality(String))))),
    `exception.thread_name` FixedString(64),
    `exception.handled` Bool,
    `exception.fingerprint` FixedString(16),
    `exception_exceptions` Array(Tuple(LowCardinality(String), LowCardinality(String), Array(Tuple(Int32, Int32, LowCardinality(String), LowCardinality(String), LowCardinality(String), LowCardinality(String))))),
    `exception_threads` Array(Tuple(LowCardinality(String), Array(Tuple(Int32, Int32, LowCardinality(String), LowCardinality(String), LowCardinality(String), LowCardinality(String))))),
    `app_exit.reason` LowCardinality(FixedString(64)),
    `app_exit.importance` LowCardinality(FixedString(32)),
    `app_exit.trace` String,
    `app_exit.process_name` String,
    `app_exit.pid` String,
    `app_exit.timestamp` DateTime64(9, 'UTC'),
    `string.severity_text` LowCardinality(FixedString(10)),
    `string.string` String,
    `gesture_long_click.target` FixedString(128),
    `gesture_long_click.target_id` FixedString(128),
    `gesture_long_click.touch_down_time` UInt32,
    `gesture_long_click.touch_up_time` UInt32,
    `gesture_long_click.width` UInt16,
    `gesture_long_click.height` UInt16,
    `gesture_long_click.x` Float32,
    `gesture_long_click.y` Float32,
    `gesture_click.target` FixedString(128),
    `gesture_click.target_id` FixedString(128),
    `gesture_click.touch_down_time` UInt32,
    `gesture_click.touch_up_time` UInt32,
    `gesture_click.width` UInt16,
    `gesture_click.height` UInt16,
    `gesture_click.x` Float32,
    `gesture_click.y` Float32,
    `gesture_scroll.target` FixedString(128),
    `gesture_scroll.target_id` FixedString(128),
    `gesture_scroll.touch_down_time` UInt32,
    `gesture_scroll.touch_up_time` UInt32,
    `gesture_scroll.x` Float32,
    `gesture_scroll.y` Float32,
    `gesture_scroll.end_x` Float32,
    `gesture_scroll.end_y` Float32,
    `gesture_scroll.direction` FixedString(8),
    `lifecycle_activity.type` FixedString(32),
    `lifecycle_activity.class_name` FixedString(128),
    `lifecycle_activity.intent` String,
    `lifecycle_activity.saved_instance_state` Bool,
    `lifecycle_fragment.type` FixedString(32),
    `lifecycle_fragment.class_name` FixedString(128),
    `lifecycle_fragment.parent_activity` String,
    `lifecycle_fragment.tag` String,
    `lifecycle_app.type` FixedString(32),
    `cold_launch.process_start_uptime` UInt32,
    `cold_launch.process_start_requested_uptime` UInt32,
    `cold_launch.content_provider_attach_uptime` UInt32,
    `cold_launch.on_next_draw_uptime` UInt32,
    `cold_launch.launched_activity` FixedString(128),
    `cold_launch.has_saved_state` Bool,
    `cold_launch.intent_data` String,
    `warm_launch.app_visible_uptime` UInt32,
    `warm_launch.on_next_draw_uptime` UInt32,
    `warm_launch.launched_activity` FixedString(128),
    `warm_launch.has_saved_state` Bool,
    `warm_launch.intent_data` String,
    `hot_launch.app_visible_uptime` UInt32,
    `hot_launch.on_next_draw_uptime` UInt32,
    `hot_launch.launched_activity` FixedString(128),
    `hot_launch.has_saved_state` Bool,
    `hot_launch.intent_data` String,
    `attributes` Map(String, String),
    `network_change.network_type` LowCardinality(FixedString(16)),
    `network_change.previous_network_type` LowCardinality(FixedString(16)),
    `network_change.network_generation` LowCardinality(FixedString(8)),
    `network_change.previous_network_generation` LowCardinality(FixedString(8)),
    `network_change.network_provider` FixedString(64),
    `anr.network_type` LowCardinality(FixedString(16)),
    `anr.network_generation` LowCardinality(FixedString(8)),
    `anr.network_provider` FixedString(64),
    `exception.network_type` LowCardinality(FixedString(16)),
    `exception.network_generation` LowCardinality(FixedString(8)),
    `exception.network_provider` FixedString(64),
    `resource.network_type` LowCardinality(FixedString(16)),
    `resource.network_generation` LowCardinality(FixedString(8)),
    `resource.network_provider` FixedString(64),
    `resource.device_locale` FixedString(64),
    `anr.device_locale` FixedString(64),
    `exception.device_locale` FixedString(64),
    `http.url` String,
    `http.method` LowCardinality(FixedString(16)),
    `http.status_code` UInt16,
    `http.request_body_size` UInt64,
    `http.response_body_size` UInt64,
    `http.request_timestamp` DateTime64(9, 'UTC'),
    `http.response_timestamp` DateTime64(9, 'UTC'),
    `http.start_time` UInt64,
    `http.end_time` UInt64,
    `http.dns_start` UInt64,
    `http.dns_end` UInt64,
    `http.connect_start` UInt64,
    `http.connect_end` UInt64,
    `http.request_start` UInt64,
    `http.request_end` UInt64,
    `http.request_headers_start` UInt64,
    `http.request_headers_end` UInt64,
    `http.request_body_start` UInt64,
    `http.request_body_end` UInt64,
    `http.response_start` UInt64,
    `http.response_end` UInt64,
    `http.response_headers_start` UInt64,
    `http.response_headers_end` UInt64,
    `http.response_body_start` UInt64,
    `http.response_body_end` UInt64,
    `http.request_headers_size` UInt64,
    `http.response_headers_size` UInt64,
    `http.failure_reason` String,
    `http.failure_description` String,
    `http_request_headers` Map(String, String),
    `http_response_headers` Map(String, String),
    `http.client` LowCardinality(FixedString(32)),
    `memory_usage.java_max_heap` UInt64,
    `memory_usage.java_total_heap` UInt64,
    `memory_usage.java_free_heap` UInt64,
    `memory_usage.total_pss` UInt64,
    `memory_usage.rss` UInt64,
    `memory_usage.native_total_heap` UInt64,
    `memory_usage.native_free_heap` UInt64,
    `memory_usage.interval_config` UInt32,
    `low_memory` Bool,
    `trim_memory.level` LowCardinality(FixedString(64)),
    `cpu_usage.num_cores` UInt8,
    `cpu_usage.clock_speed` UInt32,
    `cpu_usage.start_time` UInt64,
    `cpu_usage.uptime` UInt64,
    `cpu_usage.utime` UInt64,
    `cpu_usage.cutime` UInt64,
    `cpu_usage.stime` UInt64,
    `cpu_usage.cstime` UInt64,
    `cpu_usage.interval_config` UInt32,
)
engine = MergeTree
primary key (id, app_id, timestamp)
order by (id, app_id, timestamp);

-- migrate:down
drop table if exists default.events;