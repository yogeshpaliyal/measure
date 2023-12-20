package sh.measure.android.performance

import android.content.ComponentCallbacks2.TRIM_MEMORY_BACKGROUND
import android.content.ComponentCallbacks2.TRIM_MEMORY_COMPLETE
import android.content.ComponentCallbacks2.TRIM_MEMORY_MODERATE
import android.content.ComponentCallbacks2.TRIM_MEMORY_RUNNING_CRITICAL
import android.content.ComponentCallbacks2.TRIM_MEMORY_RUNNING_LOW
import android.content.ComponentCallbacks2.TRIM_MEMORY_RUNNING_MODERATE
import android.content.ComponentCallbacks2.TRIM_MEMORY_UI_HIDDEN
import org.junit.Before
import org.junit.Test
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import sh.measure.android.events.EventTracker
import sh.measure.android.fakes.FakeTimeProvider
import sh.measure.android.utils.CurrentThread

internal class ComponentCallbacksCollectorTest {
    private val eventTracker = mock<EventTracker>()
    private val timeProvider = FakeTimeProvider()
    private val currentThread = CurrentThread()
    private lateinit var componentCallbacksCollector: ComponentCallbacksCollector

    @Before
    fun setUp() {
        componentCallbacksCollector = ComponentCallbacksCollector(
            mock(), eventTracker, timeProvider, currentThread
        ).apply { register() }
    }

    @Test
    fun `ComponentCallbacksCollector tracks low memory event`() {
        componentCallbacksCollector.onLowMemory()

        verify(eventTracker).trackLowMemory(
            LowMemory(
                timestamp = timeProvider.currentTimeSinceEpochInMillis,
                thread_name = currentThread.name
            )
        )
    }

    @Test
    fun `ComponentCallbacksCollector tracks trim memory event`() {
        testTrimMemoryEvent(TRIM_MEMORY_UI_HIDDEN, "TRIM_MEMORY_UI_HIDDEN")
        testTrimMemoryEvent(TRIM_MEMORY_RUNNING_MODERATE, "TRIM_MEMORY_RUNNING_MODERATE")
        testTrimMemoryEvent(TRIM_MEMORY_RUNNING_LOW, "TRIM_MEMORY_RUNNING_LOW")
        testTrimMemoryEvent(TRIM_MEMORY_RUNNING_CRITICAL, "TRIM_MEMORY_RUNNING_CRITICAL")
        testTrimMemoryEvent(TRIM_MEMORY_BACKGROUND, "TRIM_MEMORY_BACKGROUND")
        testTrimMemoryEvent(TRIM_MEMORY_MODERATE, "TRIM_MEMORY_MODERATE")
        testTrimMemoryEvent(TRIM_MEMORY_COMPLETE, "TRIM_MEMORY_COMPLETE")
        testTrimMemoryEvent(999, "TRIM_MEMORY_UNKNOWN")
    }

    private fun testTrimMemoryEvent(trimLevel: Int, expectedLevel: String) {
        componentCallbacksCollector.onTrimMemory(trimLevel)
        verify(eventTracker).trackTrimMemory(
            TrimMemory(
                level = expectedLevel,
                timestamp = timeProvider.currentTimeSinceEpochInMillis,
                thread_name = currentThread.name
            )
        )
    }
}