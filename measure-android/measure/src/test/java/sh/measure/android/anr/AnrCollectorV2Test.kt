package sh.measure.android.anr

import android.os.Looper
import org.junit.Assert.assertEquals
import org.junit.Test
import org.mockito.Mockito.`when`
import org.mockito.kotlin.argumentCaptor
import org.mockito.kotlin.mock
import org.mockito.kotlin.verify
import sh.measure.NativeBridge
import sh.measure.android.events.EventProcessor
import sh.measure.android.events.EventType
import sh.measure.android.exceptions.ExceptionData
import sh.measure.android.fakes.FakeProcessInfoProvider
import sh.measure.android.fakes.NoopLogger

class AnrCollectorV2Test {
    private val logger = NoopLogger()
    private val processInfo = FakeProcessInfoProvider()
    private val eventProcessor = mock<EventProcessor>()
    private val nativeBridge = mock<NativeBridge>()
    private val looper = mock<Looper>()
    private val anrCollectorV2 =
        AnrCollectorV2(logger, processInfo, eventProcessor, nativeBridge, looper)

    @Test
    fun `register enables anr reporting and registers itself as listener`() {
        anrCollectorV2.register()
        verify(nativeBridge).enableAnrReporting(anrListener = anrCollectorV2)
    }

    @Test
    fun `unregister disables anr reporting`() {
        anrCollectorV2.unregister()
        verify(nativeBridge).disableAnrReporting()
    }

    @Test
    fun `tracks ANR event when ANR is detected`() {
        val thread = Thread.currentThread()
        `when`(looper.thread).thenReturn(thread)
        val message = "ANR"
        val timestamp = 876544454L
        val expectedAnrError = AnrError(thread, timestamp, message)

        // When
        anrCollectorV2.onAnrDetected(timestamp)

        // Then
        val typeCaptor = argumentCaptor<String>()
        val timestampCaptor = argumentCaptor<Long>()
        val dataCaptor = argumentCaptor<ExceptionData>()

        // the arguments must be in the same order as the method signature, otherwise
        // argumentCaptor will not capture the correct value and verify will fail.
        verify(eventProcessor).track(
            data = dataCaptor.capture(),
            timestamp = timestampCaptor.capture(),
            type = typeCaptor.capture(),
        )

        assertEquals(EventType.ANR, typeCaptor.firstValue)
        assertEquals(expectedAnrError.timestamp, timestampCaptor.firstValue)
        assertEquals(false, dataCaptor.firstValue.handled)
        assertEquals(processInfo.isForegroundProcess(), dataCaptor.firstValue.foreground)
    }
}