package sh.measure.android.gestures

import androidx.test.core.app.ActivityScenario
import androidx.test.espresso.Espresso.onView
import androidx.test.espresso.action.ViewActions.click
import androidx.test.espresso.action.ViewActions.swipeUp
import androidx.test.espresso.matcher.ViewMatchers.withId
import androidx.test.ext.junit.runners.AndroidJUnit4
import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith
import sh.measure.android.fakes.FakeEventTracker
import sh.measure.android.fakes.FakeTimeProvider
import sh.measure.android.fakes.NoopLogger
import sh.measure.android.test.R
import sh.measure.android.utils.CurrentThread

@RunWith(AndroidJUnit4::class)
class ComposeGestureCollectorTest {
    private val logger = NoopLogger()
    private val timeProvider = FakeTimeProvider()
    private lateinit var tracker: FakeEventTracker
    private val currentThread = CurrentThread()

    @Before
    fun setup() {
        tracker = FakeEventTracker()
        GestureCollector(logger, tracker, timeProvider, currentThread).register()
    }

    @Test
    fun tracks_clicks_on_clickable_views() {
        ActivityScenario.launch(GestureTestActivity::class.java)
        onView(withId(R.id.clickable_compose_view)).perform(click())
        assertEquals(1, tracker.trackedClicks.size)
    }

    @Test
    fun tracks_clicked_view_properties() {
        ActivityScenario.launch(GestureTestActivity::class.java)
        onView(withId(R.id.clickable_compose_view)).perform(click())

        val event = tracker.trackedClicks[0]
        assertEquals("androidx.compose.ui.platform.AndroidComposeView", event.target)
        assertEquals("compose_clickable", event.target_id)
        assertEquals("main", event.thread_name)
        assertTrue(event.touch_down_time > 0)
        assertTrue(event.touch_up_time > 0)
        assertTrue(event.x > 0)
        assertTrue(event.y > 0)

        // we currently don't have ability to track compose view bounds:
        assertNull(event.width)
        assertNull(event.height)
    }

    @Test
    fun ignores_clicks_on_non_clickable_views() {
        ActivityScenario.launch(GestureTestActivity::class.java)
        onView(withId(R.id.non_clickable_compose_view)).perform(click())
        assertEquals(0, tracker.trackedClicks.size)
    }

    @Test
    fun tracks_scrolls_on_scrollable_views() {
        ActivityScenario.launch(GestureTestActivity::class.java)
        onView(withId(R.id.scrollable_compose_view)).perform(swipeUp())
        assertEquals(1, tracker.trackedScrolls.size)
    }

    @Test
    fun tracks_scrollable_view_properties() {
        ActivityScenario.launch(GestureTestActivity::class.java)
        onView(withId(R.id.scrollable_compose_view)).perform(swipeUp())

        val event = tracker.trackedScrolls[0]
        assertEquals("androidx.compose.ui.platform.AndroidComposeView", event.target)
        assertEquals("compose_scrollable", event.target_id)
        assertEquals("main", event.thread_name)
        assertTrue(event.touch_down_time > 0)
        assertTrue(event.touch_up_time > 0)
        assertTrue(event.x > 0)
        assertTrue(event.y > 0)
    }

    @Test
    fun ignores_scrolls_on_non_scrollable_views() {
        ActivityScenario.launch(GestureTestActivity::class.java)
        onView(withId(R.id.clickable_compose_view)).perform(swipeUp())
        assertEquals(0, tracker.trackedScrolls.size)
    }
}
