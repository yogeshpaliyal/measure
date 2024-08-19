package sh.measure.android.okhttp

import sh.measure.android.events.EventProcessor
import sh.measure.android.logger.Logger
import sh.measure.android.utils.TimeProvider
import sh.measure.android.utils.isClassAvailable

/**
 * Factory for creating [HttpEventCollector] instances. This is required to
 * avoid accessing [OkHttpEventCollectorImpl] when OkHttp is not available as a runtime dependency.
 */
internal class HttpEventCollectorFactory(
    private val logger: Logger,
    private val eventProcessor: EventProcessor,
    private val timeProvider: TimeProvider,
) {
    fun create(): HttpEventCollector {
        return if (isClassAvailable("okhttp3.OkHttpClient")) {
            OkHttpEventCollectorImpl(logger, eventProcessor, timeProvider)
        } else {
            NoOpHttpEventCollector()
        }
    }
}
