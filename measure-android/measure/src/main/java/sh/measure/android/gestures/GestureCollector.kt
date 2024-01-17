package sh.measure.android.gestures

import android.view.MotionEvent
import android.view.ViewGroup
import android.view.Window
import sh.measure.android.events.EventTracker
import sh.measure.android.logger.LogLevel
import sh.measure.android.logger.Logger
import sh.measure.android.utils.CurrentThread
import sh.measure.android.utils.TimeProvider

internal class GestureCollector(
    private val logger: Logger,
    private val tracker: EventTracker,
    private val timeProvider: TimeProvider,
    private val currentThread: CurrentThread,
) {
    fun register() {
        logger.log(LogLevel.Debug, "Registering gesture collector")
        WindowInterceptor().apply {
            init()
            registerInterceptor(object : WindowTouchInterceptor {
                override fun intercept(motionEvent: MotionEvent, window: Window) {
                    trackGesture(motionEvent, window)
                }
            })
        }
    }

    private fun trackGesture(motionEvent: MotionEvent, window: Window) {
        val gesture = GestureDetector.detect(window.context, motionEvent, timeProvider, currentThread)
        if (gesture == null || motionEvent.action != MotionEvent.ACTION_UP) {
            return
        }
        // Find the potential view on which the gesture ended on.
        val target = getTarget(gesture, window, motionEvent)
        if (target == null) {
            logger.log(
                LogLevel.Debug,
                "No target found for gesture ${gesture.javaClass.simpleName}",
            )
            return
        } else {
            logger.log(
                LogLevel.Debug,
                "Target found for gesture ${gesture.javaClass.simpleName}: ${target.className}:${target.id}",
            )
        }

        when (gesture) {
            is DetectedGesture.Click -> tracker.trackClick(
                ClickEvent.fromDetectedGesture(gesture, target),
            )

            is DetectedGesture.LongClick -> tracker.trackLongClick(
                LongClickEvent.fromDetectedGesture(gesture, target),
            )

            is DetectedGesture.Scroll -> tracker.trackScroll(
                ScrollEvent.fromDetectedGesture(gesture, target),
            )
        }
    }

    private fun getTarget(
        gesture: DetectedGesture,
        window: Window,
        motionEvent: MotionEvent,
    ): Target? {
        return when (gesture) {
            is DetectedGesture.Scroll -> {
                GestureTargetFinder.findScrollable(
                    window.decorView as ViewGroup,
                    motionEvent,
                )
            }

            else -> {
                GestureTargetFinder.findClickable(
                    window.decorView as ViewGroup,
                    motionEvent,
                )
            }
        }
    }
}
