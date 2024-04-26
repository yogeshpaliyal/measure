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
import sh.measure.android.events.EventType
import sh.measure.android.fakes.FakeEventProcessor
import sh.measure.android.fakes.FakeTimeProvider

@RunWith(AndroidJUnit4::class)
class ComposeNavigationCollectorTest {
    private val timeProvider = FakeTimeProvider()
    private lateinit var tracker: FakeEventProcessor

    @get:Rule
    val composeRule = createComposeRule()

    @Before
    fun setup() {
        tracker = FakeEventProcessor()
        Measure.setEventProcessor(tracker)
        Measure.setTimeProvider(timeProvider)
    }

    @Test
    fun tracks_navigation_event_for_compose_navigation() {
        composeRule.setContent {
            testApp()
        }

        // initial state
        assertEquals(1, tracker.getTrackedEventsByType(EventType.NAVIGATION).size)
        assertEquals(
            "home",
            (tracker.getTrackedEventsByType(EventType.NAVIGATION)[0].data as NavigationData).route
        )

        // forward navigation
        composeRule.onNodeWithText("Checkout").performClick()
        assertEquals(2, tracker.getTrackedEventsByType(EventType.NAVIGATION).size)
        assertEquals(
            "checkout",
            (tracker.getTrackedEventsByType(EventType.NAVIGATION)[1].data as NavigationData).route
        )

        // back
        pressBack()
        assertEquals(3, tracker.getTrackedEventsByType(EventType.NAVIGATION).size)
        assertEquals(
            "home",
            (tracker.getTrackedEventsByType(EventType.NAVIGATION)[2].data as NavigationData).route
        )
    }
}
