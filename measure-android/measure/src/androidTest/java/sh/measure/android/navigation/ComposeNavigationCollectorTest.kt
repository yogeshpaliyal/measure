package sh.measure.android.navigation

import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import androidx.test.espresso.Espresso.pressBack
import androidx.test.ext.junit.runners.AndroidJUnit4
import junit.framework.TestCase.assertEquals
import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith
import sh.measure.android.Measure
import sh.measure.android.fakes.FakeEventTracker
import sh.measure.android.fakes.FakeTimeProvider
import sh.measure.android.utils.CurrentThread

@RunWith(AndroidJUnit4::class)
class ComposeNavigationCollectorTest {
    private val timeProvider = FakeTimeProvider()
    private lateinit var tracker: FakeEventTracker
    private val currentThread = CurrentThread()

    @get:Rule
    val composeRule = createComposeRule()

    @Before
    fun setup() {
        tracker = FakeEventTracker()
        Measure.setEventTracker(tracker)
        Measure.setTimeProvider(timeProvider)
        Measure.setCurrentThread(currentThread)
    }

    @Test
    fun tracks_navigation_event_for_compose_navigation() {
        composeRule.setContent {
            testApp()
        }

        // initial state
        assertEquals(1, tracker.trackedNavigationEvents.size)
        assertEquals("home", tracker.trackedNavigationEvents[0].route)

        // forward navigation
        composeRule.onNodeWithText("Checkout").performClick()
        assertEquals(2, tracker.trackedNavigationEvents.size)
        assertEquals("checkout", tracker.trackedNavigationEvents[1].route)

        // back
        pressBack()
        assertEquals(3, tracker.trackedNavigationEvents.size)
        assertEquals("home", tracker.trackedNavigationEvents[2].route)
    }
}