# Feature - App Launch

Measure tracks the cold, warm and hot app launch time automatically. No additional code is required to enable this
feature. A [method trace for cold launch](#cold-launch-method-trace) is also captured to help debug any bottlenecks.

## How it works

* [Cold launch](#cold-launch)
* [Warm launch](#warm-launch)
* [Hot launch](#hot-launch)

### Cold launch

A [cold launch](https://developer.android.com/topic/performance/vitals/launch-time#cold) refers to an app starting
from scratch. Cold launch happens in cases such as an app launching for the first time since the device booted or since
the system killed the app.

There are typically two important metrics to track for cold launch. Time to Initial Display (TTID) and Time to Full
Display (TTFD).

> [!NOTE]  
> Measuring TTFD is not possible yet, support will be added in a future version.

To measure **Time to Initial Display (TTID)**, two timestamps are required:

1. The time when the app was launched.
2. The time when the app's first frame was displayed.

_The time when app was launched_ is calculated differently for different SDK versions, we try to use the most accurate
measurement possible for the given SDK version.

* Up to API 24: the _uptime_ time when Measure content provider's attachInfo callback is invoked.
* API 24 - API 32: the process start uptime, using `Process.getStartUptimeMillis()`
* API 33 and beyond: the process start uptime, using `Process.getStartRequestedUptimeMillis()`

_The time when app's first frame was displayed_ is a bit more complex. Simplifying some of the steps, it is calculated
in the following way:

1. Get the decor view by registering
   [onContentChanged](https://developer.android.com/reference/android/app/Activity#onContentChanged()) callback on the
   first Activity.
2. Get the next draw callback by
   registering [OnDrawListener](https://developer.android.com/reference/android/view/ViewTreeObserver.OnDrawListener) on
   the decor view.
3. Post a runnable in front of the next draw callback to record the time just before the first frame was displayed. This
   is the most accurate time we can get to calculate TTID.

### Warm launch

A [warm launch](https://developer.android.com/topic/performance/vitals/launch-time#warm) refers to the re-launch of an
app causing an Activity `onCreate` to be triggered instead of just `onResume`. This requires the system to recreate
the activity from scratch and hence requires more work than a hot launch.

Warm launch is calculated by keeping track of the time when the Activity `onCreate` of the Activity being recreated is
triggered and the time when the first frame is displayed. The same method as for cold launch is used to calculate the
time when the first frame is displayed.

### Hot launch

A [hot launch](https://developer.android.com/topic/performance/vitals/launch-time#hot) refers to the re-launch of an
app causing an Activity `onResume` to be triggered. This typically requires less work than a warm launch as the system
does not need to recreate the activity from scratch. However, if there were any trim memory events leading to the
certain resources being released, the system might need to recreate those resources.

## Cold launch method trace

Measure automatically tracks a method trace for every cold launch. The method trace is captured from the time when the
Measure SDK is initialized up to the time when the first frame is displayed. This method trace can be used to identify
any bottlenecks in the app's startup process.

> [!NOTE]  
> The method trace is only available for cold launches and is currently collected for every cold launch with hardcoded
> interval. A configuration will be exposed to control this in the future. The progress can be tracked
> [here](https://github.com/measure-sh/measure/issues/550)

Any custom traces added using [Trace](https://developer.android.com/reference/kotlin/androidx/tracing/Trace) will also
be captured along with the method trace, we recommend using this API to add custom traces to the method trace to provide
a detailed view of the app's startup process.

To view the trace, any of the following tools can be used:

1. [YAMP](https://github.com/Grigory-Rylov/android-methods-profiler) - it supports opening the trace file,
   applying symbolication by providing the mapping file and an advanced search.
2. [Android Studio](https://developer.android.com/studio/profile/cpu-profiler) - drag and drop the trace file to view it
   in Android Studio.

## Data collected

Checkout the data collected by Measure
for [Cold Launch](../../../docs/api/sdk/README.md#coldlaunch), [Warm Launch](../../../docs/api/sdk/README.md#warmlaunch)
and [Hot Launch](../../../docs/api/sdk/README.md#hotlaunch) sections respectively.

### Further reading

* [Android docs on app startup](https://developer.android.com/topic/performance/vitals/launch-time#warm)
* [Py's android vitals series](https://dev.to/pyricau/series/7827)
* [Py's PAPA github project](https://github.com/square/papa)