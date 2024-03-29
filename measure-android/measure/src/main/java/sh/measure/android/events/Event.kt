package sh.measure.android.events

import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonElement
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.encodeToStream
import okio.BufferedSink
import sh.measure.android.appexit.AppExit
import sh.measure.android.applaunch.ColdLaunchEvent
import sh.measure.android.applaunch.HotLaunchEvent
import sh.measure.android.applaunch.WarmLaunchEvent
import sh.measure.android.exceptions.MeasureException
import sh.measure.android.gestures.ClickEvent
import sh.measure.android.gestures.LongClickEvent
import sh.measure.android.gestures.ScrollEvent
import sh.measure.android.lifecycle.ActivityLifecycleEvent
import sh.measure.android.lifecycle.ApplicationLifecycleEvent
import sh.measure.android.lifecycle.FragmentLifecycleEvent
import sh.measure.android.navigation.NavigationEvent
import sh.measure.android.networkchange.NetworkChangeEvent
import sh.measure.android.okhttp.HttpEvent
import sh.measure.android.performance.CpuUsage
import sh.measure.android.performance.LowMemory
import sh.measure.android.performance.MemoryUsage
import sh.measure.android.performance.TrimMemory
import sh.measure.android.utils.iso8601Timestamp
import sh.measure.android.utils.toJsonElement

internal data class Event(
    val timestamp: String,
    val type: String,
    val data: JsonElement,
    val attributes: JsonObject,
) {
    fun toJson(): String {
        val serializedData = Json.encodeToString(JsonElement.serializer(), data)
        return "{\"timestamp\":\"$timestamp\",\"type\":\"$type\",\"$type\":$serializedData,\"attributes\":$attributes}"
    }

    @OptIn(ExperimentalSerializationApi::class)
    fun write(sink: BufferedSink) {
        sink.writeUtf8("{")
        sink.writeUtf8("\"timestamp\":\"${timestamp}\",")
        sink.writeUtf8("\"type\":\"${type}\",")
        sink.writeUtf8("\"${type}\":")
        Json.encodeToStream(JsonElement.serializer(), data, sink.outputStream())
        sink.writeUtf8(",")
        sink.writeUtf8("\"attributes\":")
        Json.encodeToStream(JsonObject.serializer(), attributes, sink.outputStream())
        sink.writeUtf8("}")
    }
}

internal fun MeasureException.toEvent(): Event {
    return Event(
        type = if (isAnr) EventType.ANR else EventType.EXCEPTION,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(MeasureException.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun AppExit.toEvent(): Event {
    return Event(
        type = EventType.APP_EXIT,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(AppExit.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun ClickEvent.toEvent(): Event {
    return Event(
        timestamp = timestamp.iso8601Timestamp(),
        type = EventType.CLICK,
        data = Json.encodeToJsonElement(ClickEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun LongClickEvent.toEvent(): Event {
    return Event(
        timestamp = timestamp.iso8601Timestamp(),
        type = EventType.LONG_CLICK,
        data = Json.encodeToJsonElement(LongClickEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun ScrollEvent.toEvent(): Event {
    return Event(
        timestamp = timestamp.iso8601Timestamp(),
        type = EventType.SCROLL,
        data = Json.encodeToJsonElement(ScrollEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun ApplicationLifecycleEvent.toEvent(): Event {
    return Event(
        type = EventType.LIFECYCLE_APP,
        timestamp = timestamp,
        data = Json.encodeToJsonElement(ApplicationLifecycleEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun ActivityLifecycleEvent.toEvent(): Event {
    return Event(
        type = EventType.LIFECYCLE_ACTIVITY,
        timestamp = timestamp,
        data = Json.encodeToJsonElement(ActivityLifecycleEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun FragmentLifecycleEvent.toEvent(): Event {
    return Event(
        type = EventType.LIFECYCLE_FRAGMENT,
        timestamp = timestamp,
        data = Json.encodeToJsonElement(FragmentLifecycleEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun ColdLaunchEvent.toEvent(): Event {
    return Event(
        type = EventType.COLD_LAUNCH,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(ColdLaunchEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun WarmLaunchEvent.toEvent(): Event {
    return Event(
        type = EventType.WARM_LAUNCH,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(WarmLaunchEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun HotLaunchEvent.toEvent(): Event {
    return Event(
        type = EventType.HOT_LAUNCH,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(HotLaunchEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun NetworkChangeEvent.toEvent(): Event {
    return Event(
        type = EventType.NETWORK_CHANGE,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(NetworkChangeEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun HttpEvent.toEvent(): Event {
    return Event(
        type = EventType.HTTP,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(HttpEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun MemoryUsage.toEvent(): Event {
    return Event(
        type = EventType.MEMORY_USAGE,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(MemoryUsage.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun LowMemory.toEvent(): Event {
    return Event(
        type = EventType.LOW_MEMORY,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(LowMemory.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun TrimMemory.toEvent(): Event {
    return Event(
        type = EventType.TRIM_MEMORY,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(TrimMemory.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun CpuUsage.toEvent(): Event {
    return Event(
        type = EventType.CPU_USAGE,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(CpuUsage.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}

internal fun NavigationEvent.toEvent(): Event {
    return Event(
        type = EventType.NAVIGATION,
        timestamp = timestamp.iso8601Timestamp(),
        data = Json.encodeToJsonElement(NavigationEvent.serializer(), this),
        attributes = attributes.toJsonElement(),
    )
}
