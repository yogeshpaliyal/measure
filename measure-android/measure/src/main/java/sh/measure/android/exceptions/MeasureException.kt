package sh.measure.android.exceptions

import kotlinx.serialization.Serializable
import kotlinx.serialization.Transient
import kotlinx.serialization.json.Json
import sh.measure.android.session.iso8601Timestamp
import sh.measure.android.storage.Signal
import sh.measure.android.tracker.EventType
import sh.measure.android.tracker.SignalType

/**
 * Represents an exception in Measure. This is used to track handled and unhandled exceptions.
 */
@Serializable
internal data class MeasureException(
    /**
     * The timestamp at which the exception occurred.
     */
    @Transient
    val timestamp: Long = 0L,

    /**
     * The name of the thread that threw the exception.
     */
    val thread_name: String,

    /**
     * A list of exceptions that were thrown. Multiple exceptions represent "chained" exceptions.
     */
    val exceptions: List<ExceptionUnit>,

    /**
     * The stacktrace of all the threads at the time of the exception.
     */
    val threads: List<MeasureThread>,

    /**
     * Whether the exception was handled or not.
     */
    val handled: Boolean
) {
    fun toSignal(sessionId: String): Signal {
        return Signal(
            sessionId = sessionId,
            timestamp = timestamp.iso8601Timestamp(),
            signalType = SignalType.EVENT,
            dataType = EventType.EXCEPTION,
            data = Json.encodeToString(serializer(), this)
        )
    }
}

@Serializable
internal data class MeasureThread(
    val name: String, val frames: List<Frame>
)

/**
 * Represents a stacktrace in Measure.
 */
@Serializable
internal data class ExceptionUnit(
    /**
     * The fully qualified type of the exception. For example, java.lang.Exception.
     */
    val type: String,

    /**
     * A message which describes the exception.
     */
    val message: String? = null,

    /**
     * A list of stack frames for the exception.
     */
    val frames: List<Frame>
)

/**
 * Represents a stackframe in Measure.
 */
@Serializable
internal data class Frame(
    /**
     * The fully qualified class name.
     */
    val class_name: String? = null,

    /**
     * The name of the method in the stacktrace.
     */
    val method_name: String? = null,

    /**
     * The name of the source file in the stacktrace.
     */
    val file_name: String? = null,

    /**
     * The line number of the method called.
     */
    val line_num: Int? = null,

    val module_name: String? = null,

    val col_num: Int? = null,
)