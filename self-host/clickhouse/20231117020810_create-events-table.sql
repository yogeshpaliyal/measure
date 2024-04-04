-- migrate:up
create table if not exists default.events
(
  `id` UUID comment 'unique event id',
  `type` LowCardinality(FixedString(32)) comment 'type of the event',
  `session_id` UUID comment 'associated session id',
  `app_id` UUID comment 'associated app id',
  `inet.ipv4` Nullable(IPv4) comment 'ipv4 address',
  `inet.ipv6` Nullable(IPv6) comment 'ipv6 address',
  `inet.country_code` FixedString(8) comment 'country code',
  `timestamp` DateTime64(9, 'UTC') comment 'event timestamp',
  `attribute.installation_id` UUID not null comment 'unique id for an installation of an app, generated by sdk',
  `attribute.app_version` FixedString(32) not null comment 'app version identifier',
  `attribute.app_build` FixedString(32) not null comment 'app build identifier',
  `attribute.app_unique_id` FixedString(128) not null comment 'app bundle identifier',
  `attribute.platform` LowCardinality(FixedString(32)) not null comment 'platform identifier',
  `attribute.measure_sdk_version` FixedString(16) not null comment 'measure sdk version identifier',
  `attribute.thread_name` FixedString(64) comment 'thread on which the event was captured',
  `attribute.user_id` FixedString(128) comment 'id of the app''s end user',
  `attribute.device_name` FixedString(32) comment 'name of the device',
  `attribute.device_model` FixedString(32) comment 'model of the device',
  `attribute.device_manufacturer` FixedString(32) comment 'manufacturer of the device',
  `attribute.device_type` LowCardinality(FixedString(32)) comment 'type of the device, like phone or tablet',
  `attribute.device_is_foldable` Bool comment 'true for foldable devices',
  `attribute.device_is_physical` Bool comment 'true for physical devices',
  `attribute.device_density_dpi` UInt16 comment 'dpi density',
  `attribute.device_width_px` UInt16 comment 'screen width',
  `attribute.device_height_px` UInt16 comment 'screen height',
  `attribute.device_density` Float32 comment 'device density',
  `attribute.device_locale` FixedString(64) comment 'rfc 5646 locale string',
  `attribute.os_name` FixedString(32) comment 'name of the operating system',
  `attribute.os_version` FixedString(32) comment 'version of the operating system',
  `attribute.network_type` LowCardinality(FixedString(16)) comment 'either - wifi, cellular, vpn, unknown, no_network',
  `attribute.network_generation` LowCardinality(FixedString(8)) comment 'either - 2g, 3g, 4g, 5g',
  `attribute.network_provider` FixedString(64) comment 'name of the network service provider',
  `anr.handled` Bool comment 'anr was handled by the application code',
  `anr.fingerprint` FixedString(16) comment 'fingerprint for anr similarity classification',
  `anr.exceptions` String comment 'anr exception data',
  `anr.threads` String comment 'anr thread data',
  `anr.foreground` Bool comment 'true if the anr was perceived by end user',
  `exception.handled` Bool comment 'exception was handled by application code',
  `exception.fingerprint` FixedString(16) comment 'fingerprint for exception similarity classification',
  `exception.exceptions` String comment 'exception data',
  `exception.threads` String comment 'exception thread data',
  `exception.foreground` Bool comment 'true if the exception was perceived by end user',
  `app_exit.reason` LowCardinality(FixedString(64)) comment 'reason for app exit',
  `app_exit.importance` LowCardinality(FixedString(32)) comment 'importance of process that it used to have before death',
  `app_exit.trace` String comment 'modified trace given by ApplicationExitInfo to help debug anrs.',
  `app_exit.process_name` String comment 'name of the process that died',
  `app_exit.pid` String comment 'id of the process that died',
  `string.severity_text` LowCardinality(FixedString(10)) comment 'log level - info, warning, error, fatal, debug',
  `string.string` String comment 'log message text',
  `gesture_long_click.target` FixedString(128) comment 'class or instance name of the originating view',
  `gesture_long_click.target_id` FixedString(128) comment 'unique identifier of the target',
  `gesture_long_click.touch_down_time` UInt32 comment 'time for touch down gesture',
  `gesture_long_click.touch_up_time` UInt32 comment 'time for touch up gesture',
  `gesture_long_click.width` UInt16 comment 'width of the target view in pixels',
  `gesture_long_click.height` UInt16 comment 'height of the target view in pixels',
  `gesture_long_click.x` Float32 comment 'x coordinate of where the gesture happened',
  `gesture_long_click.y` Float32 comment 'y coordinate of where the gesture happened',
  `gesture_click.target` FixedString(128) comment 'class or instance name of the originating view',
  `gesture_click.target_id` FixedString(128) comment 'unique identifier of the target',
  `gesture_click.touch_down_time` UInt32 comment 'time for touch down gesture',
  `gesture_click.touch_up_time` UInt32 comment 'time for the touch up gesture',
  `gesture_click.width` UInt16 comment 'width of the target view in pixels',
  `gesture_click.height` UInt16 comment 'height of the target view in pixels',
  `gesture_click.x` Float32 comment 'x coordinate of where the gesture happened',
  `gesture_click.y` Float32 comment 'y coordinate of where the gesture happened',
  `gesture_scroll.target` FixedString(128) comment 'class or instance name of the originating view',
  `gesture_scroll.target_id` FixedString(128) comment 'unique identifier of the target',
  `gesture_scroll.touch_down_time` UInt32 comment 'time for touch down gesture',
  `gesture_scroll.touch_up_time` UInt32 comment 'time for touch up gesture',
  `gesture_scroll.x` Float32 comment 'x coordinate of where the gesture started',
  `gesture_scroll.y` Float32 comment 'y coordinate of where the gesture started',
  `gesture_scroll.end_x` Float32 comment 'x coordinate of where the gesture ended',
  `gesture_scroll.end_y` Float32 comment 'y coordinate of where the gesture ended',
  `gesture_scroll.direction` FixedString(8) comment 'direction of the scroll',
  `lifecycle_activity.type` FixedString(32) comment 'type of the lifecycle activity, either - created, resumed, paused, destroyed',
  `lifecycle_activity.class_name` FixedString(128) comment 'fully qualified class name of the activity',
  `lifecycle_activity.intent` String comment 'intent data serialized as string',
  `lifecycle_activity.saved_instance_state` Bool comment 'represents that activity was recreated with a saved state. only available for type created.',
  `lifecycle_fragment.type` FixedString(32) comment 'type of the lifecycle fragment, either - attached, resumed, paused, detached',
  `lifecycle_fragment.class_name` FixedString(128) comment 'fully qualified class name of the fragment',
  `lifecycle_fragment.parent_activity` String comment 'fully qualified class name of the parent activity that the fragment is attached to',
  `lifecycle_fragment.tag` String comment 'optional fragment tag',
  `lifecycle_app.type` FixedString(32) comment 'type of the lifecycle app, either - background, foreground',
  `cold_launch.process_start_uptime` UInt32 comment 'start uptime in msec',
  `cold_launch.process_start_requested_uptime` UInt32 comment 'start uptime in msec',
  `cold_launch.content_provider_attach_uptime` UInt32 comment 'start uptime in msec',
  `cold_launch.on_next_draw_uptime` UInt32 comment 'time at which app became visible',
  `cold_launch.launched_activity` FixedString(128) comment 'activity which drew the first frame during cold launch',
  `cold_launch.has_saved_state` Bool comment 'whether the launched_activity was created with a saved state bundle',
  `cold_launch.intent_data` String comment 'intent data used to launch the launched_activity',
  `cold_launch.duration` UInt32 comment 'computed cold launch duration',
  `warm_launch.app_visible_uptime` UInt32 comment 'time since the app became visible to user, in msec',
  `warm_launch.on_next_draw_uptime` UInt32 comment 'time at which app became visible to user, in msec',
  `warm_launch.launched_activity` FixedString(128) comment 'activity which drew the first frame during warm launch',
  `warm_launch.has_saved_state` Bool comment 'whether the launched_activity was created with a saved state bundle',
  `warm_launch.intent_data` String comment 'intent data used to launch the launched_activity',
  `warm_launch.duration` UInt32 comment 'computed warm launch duration',
  `hot_launch.app_visible_uptime` UInt32 comment 'time elapsed since the app became visible to user, in msec',
  `hot_launch.on_next_draw_uptime` UInt32 comment 'time at which app became visible to user, in msec',
  `hot_launch.launched_activity` FixedString(128) comment 'activity which drew the first frame during hot launch',
  `hot_launch.has_saved_state` Bool comment 'whether the launched_activity was created with a saved state bundle',
  `hot_launch.intent_data` String comment 'intent data used to launch the launched_activity',
  `hot_launch.duration` UInt32 comment 'computed hot launch duration',
  `network_change.network_type` LowCardinality(FixedString(16)) comment 'type of the network, wifi, cellular etc',
  `network_change.previous_network_type` LowCardinality(FixedString(16)) comment 'type of the previous network',
  `network_change.network_generation` LowCardinality(FixedString(8)) comment '2g, 3g, 4g etc',
  `network_change.previous_network_generation` LowCardinality(FixedString(8)) comment 'previous network generation',
  `network_change.network_provider` FixedString(64) comment 'name of the network service provider',
  `http.url` String comment 'url of the http request',
  `http.method` LowCardinality(FixedString(16)) comment 'method like get, post',
  `http.status_code` UInt16 comment 'http status code',
  `http.start_time` UInt64 comment 'uptime at when the http call started, in msec',
  `http.end_time` UInt64 comment 'uptime at when the http call ended, in msec',
  `http_request_headers` Map(String, String) comment 'http request headers',
  `http_response_headers` Map(String, String) comment 'http response headers',
  `http.request_body` String comment 'request body',
  `http.response_body` String comment 'response body',
  `http.failure_reason` String comment 'reason for failure',
  `http.failure_description` String comment 'description of the failure',
  `http.client` LowCardinality(FixedString(32)) comment 'name of the http client',
  `memory_usage.java_max_heap` UInt64 comment 'maximum size of the java heap allocated, in kb',
  `memory_usage.java_total_heap` UInt64 comment 'total size of the java heap available for allocation, in KB',
  `memory_usage.java_free_heap` UInt64 comment 'free memory available in the java heap, in kb',
  `memory_usage.total_pss` UInt64 comment 'total proportional set size - amount of memory used by the process, including shared memory and code. in kb.',
  `memory_usage.rss` UInt64 comment 'resident set size - amount of physical memory currently used, in kb',
  `memory_usage.native_total_heap` UInt64 comment 'total size of the native heap (memory out of java''s control) available for allocation, in kb',
  `memory_usage.native_free_heap` UInt64 comment 'amount of free memory available in the native heap, in kb',
  `memory_usage.interval_config` UInt32 comment 'interval between two consecutive readings, in msec',
  `low_memory.java_max_heap` UInt64 comment 'maximum size of the java heap allocated, in kb',
  `low_memory.java_total_heap` UInt64 comment 'total size of the java heap available for allocation, in kb',
  `low_memory.java_free_heap` UInt64 comment 'free memory available in the java heap, in kb',
  `low_memory.total_pss` UInt64 comment 'total proportional set size - amount of memory used by the process, including shared memory and code. in kb.',
  `low_memory.rss` UInt64 comment 'resident set size - amount of physical memory currently used, in kb',
  `low_memory.native_total_heap` UInt64 comment 'total size of the native heap (memory out of java',
  `low_memory.native_free_heap` UInt64 comment 'amount of free memory available in the native heap, in kb',
  `trim_memory.level` LowCardinality(FixedString(64)) comment 'one of the trim memory constants as received by component callback',
  `cpu_usage.num_cores` UInt8 comment 'number of cores on the device',
  `cpu_usage.clock_speed` UInt32 comment 'clock speed of the processor, in hz',
  `cpu_usage.uptime` UInt64 comment 'time since the device booted, in msec',
  `cpu_usage.utime` UInt64 comment 'execution time in user mode, in jiffies',
  `cpu_usage.cutime` UInt64 comment 'execution time in user mode with child processes, in jiffies',
  `cpu_usage.stime` UInt64 comment 'execution time in kernel mode, in jiffies',
  `cpu_usage.cstime` UInt64 comment 'execution time in user mode with child processes, in jiffies',
  `cpu_usage.interval_config` UInt32 comment 'interval between two consecutive readings, in msec',
  `navigation.route` FixedString(128) comment 'the destination route',
)
engine = MergeTree
primary key (id, app_id, timestamp)
order by (id, app_id, timestamp)
comment 'events master table';

-- migrate:down
drop table if exists default.events;