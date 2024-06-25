package sh.measure.android.exceptions

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import sh.measure.android.events.EventProcessor
import sh.measure.android.events.EventType
import sh.measure.android.fakes.FakeProcessInfoProvider
import sh.measure.android.fakes.FakeTimeProvider
import sh.measure.android.fakes.NoopLogger

internal class UnhandledExceptionCollectorTest {

    private var originalDefaultHandler: Thread.UncaughtExceptionHandler? = null
    private val logger = NoopLogger()
    private val timeProvider = FakeTimeProvider()
    private val eventProcessor = mock<EventProcessor>()
    private val processInfo = FakeProcessInfoProvider()

    @Before
    fun setUp() {
        originalDefaultHandler = Thread.getDefaultUncaughtExceptionHandler()
    }

    @Test
    fun `UnhandledExceptionCollector registers itself as an uncaught exception handler`() {
        // When
        val collector = UnhandledExceptionCollector(
            logger,
            eventProcessor,
            timeProvider,
            processInfo,
        ).apply { register() }
        val currentDefaultHandler = Thread.getDefaultUncaughtExceptionHandler()

        // Then
        assertEquals(collector, currentDefaultHandler)
    }

    @Test
    fun `UnhandledExceptionCollector tracks uncaught exceptions`() {
        val collector = UnhandledExceptionCollector(
            logger,
            eventProcessor,
            timeProvider,
            processInfo,
        ).apply { register() }

        // Given
        val thread = Thread.currentThread()
        val exception = RuntimeException("Test exception")
        val expectedException = ExceptionFactory.createMeasureException(
            exception,
            handled = false,
            thread = thread,
            foreground = processInfo.isForegroundProcess(),
        )

        // When
        collector.uncaughtException(thread, exception)

        // Then
        verify(eventProcessor).trackCrash(
            timestamp = timeProvider.currentTimeSinceEpochInMillis,
            type = EventType.EXCEPTION,
            data = expectedException,
            attributes = mutableMapOf(),
            attachments = mutableListOf(),
        )
    }

    @Test
    fun `UnhandledExceptionCollector calls the original handler after capturing the exception`() {
        var originalHandlerCalled = false
        Thread.setDefaultUncaughtExceptionHandler { _, _ ->
            originalHandlerCalled = true
        }
        val collector = UnhandledExceptionCollector(
            logger,
            eventProcessor,
            timeProvider,
            processInfo,
        ).apply { register() }

        // Given
        val thread = Thread.currentThread()
        val exception = RuntimeException("Test exception")

        // When
        collector.uncaughtException(thread, exception)

        // Then
        assertTrue(originalHandlerCalled)
    }
}
