package sh.measure.android.config

/**
 * Defines all the configuration options for the Measure SDK.
 */
internal interface IMeasureConfig {
    /**
     * Whether to capture a screenshot of the app when it crashes due to an unhandled exception or
     * ANR. Defaults to `true`.
     */
    val trackScreenshotOnCrash: Boolean

    /**
     * The level of masking to apply to the screenshot. Defaults to [ScreenshotMaskLevel.AllTextAndMedia].
     */
    val screenshotMaskLevel: ScreenshotMaskLevel

    /**
     * The color of the mask to apply to the screenshot. The value should be a hex color string.
     * For example, "#222222".
     */
    val screenshotMaskHexColor: String

    /**
     * The compression quality of the screenshot. Must be between 0 and 100, where 0 is lowest quality
     * and smallest size while 100 is highest quality and largest size.
     */
    val screenshotCompressionQuality: Int

    /**
     * Whether to capture http headers of a network request and response. Defaults to `false`.
     */
    val trackHttpHeaders: Boolean

    /**
     * Whether to capture http body of a network request and response. Defaults to `false`.
     */
    val trackHttpBody: Boolean

    /**
     * List of HTTP headers to not capture for network request and response. Defaults to an empty
     * list.
     *
     * Internally, this list is combined with [defaultHttpHeadersBlocklist] to form the final
     * blocklist.
     */
    val httpHeadersBlocklist: List<String>

    /**
     * Allows disabling collection of `http` events for certain URLs. This is useful to setup if you do not
     * want to collect data for certain endpoints or third party domains. By default, Measure endpoints
     * are always disabled.
     *
     * The check is made in order of the list and uses a simple `contains` check to see if the URL
     * contains any of the strings in the list.
     *
     * Internally, this list is combined with [defaultHttpUrlBlocklist] to form the final blocklist.
     *
     * Example:
     *
     * ```kotlin
     * MeasureConfig(
     *     httpUrlBlocklist = listOf(
     *         "example.com", // disables a domain
     *         "api.example.com", // disable a subdomain
     *         "example.com/order" // disable a particular path
     *     )
     * )
     * ```
     */
    val httpUrlBlocklist: List<String>

    /**
     * Whether to capture lifecycle activity intent data. Defaults to `false`.
     */
    val trackActivityIntentData: Boolean

    /**
     * The maximum size of attachments allowed in a single batch. Defaults to 3MB
     */
    val maxEventsAttachmentSizeInBatchBytes: Int

    /**
     * The interval at which to create a batch for export.
     */
    val eventsBatchingIntervalMs: Long

    /**
     * The maximum number of events to export in /events API. Defaults to 500.
     */
    val maxEventsInBatch: Int

    /**
     * When `httpBodyCapture` is enabled, this determines whether to capture the body or not based
     * on the content type of the request/response. Defaults to `application/json`.
     */
    val httpContentTypeAllowlist: List<String>

    /**
     * Default list of HTTP headers to not capture for network request and response.
     */
    val defaultHttpHeadersBlocklist: List<String>

    /**
     * Default list of HTTP URLs to not capture for network request and response.
     */
    val defaultHttpUrlBlocklist: List<String>
}

class MeasureConfig(
    override val trackScreenshotOnCrash: Boolean = true,
    override val screenshotMaskLevel: ScreenshotMaskLevel = ScreenshotMaskLevel.AllTextAndMedia,
    override val trackHttpHeaders: Boolean = false,
    override val trackHttpBody: Boolean = false,
    override val httpHeadersBlocklist: List<String> = emptyList(),
    override val httpUrlBlocklist: List<String> = emptyList(),
    override val trackActivityIntentData: Boolean = false,
) : IMeasureConfig {
    override val screenshotMaskHexColor: String = "#222222"
    override val screenshotCompressionQuality: Int = 25
    override val maxEventsAttachmentSizeInBatchBytes: Int = 3
    override val eventsBatchingIntervalMs: Long = 30_000 // 30 seconds
    override val maxEventsInBatch: Int = 500
    override val httpContentTypeAllowlist: List<String> = listOf("application/json")
    override val defaultHttpHeadersBlocklist: List<String> = listOf(
        "Authorization",
        "Cookie",
        "Set-Cookie",
        "Proxy-Authorization",
        "WWW-Authenticate",
        "X-Api-Key",
    )
    override val defaultHttpUrlBlocklist: List<String> = listOf(
        // TODO(abhay): review this list to block all measure API endpoints.
        "api.measure.sh",
        "10.0.2.2:8080/events",
    )
}
