package sh.measure.android.exporter

import sh.measure.android.logger.LogLevel
import sh.measure.android.logger.Logger
import sh.measure.android.storage.FileStorage
import java.net.URL

internal interface NetworkClient {
    fun init(baseUrl: String, apiKey: String)
    fun execute(
        batchId: String,
        eventPackets: List<EventPacket>,
        attachmentPackets: List<AttachmentPacket>,
        spanPackets: List<SpanPacket>,
    ): HttpResponse
}

internal class NetworkClientImpl(
    private val logger: Logger,
    private val fileStorage: FileStorage,
    private val httpClient: HttpClient = HttpUrlConnectionClient(logger),
    private val multipartDataFactory: MultipartDataFactory = MultipartDataFactoryImpl(
        logger,
        fileStorage,
    ),
) : NetworkClient {
    private var baseUrl: URL? = null
    private var eventsUrl: URL? = null
    private var apiKey: String? = null

    companion object {
        private const val PATH_EVENTS = "/events"
    }

    override fun init(baseUrl: String, apiKey: String) {
        this.baseUrl = parseUrl(baseUrl)
        this.apiKey = apiKey
        this.eventsUrl = this.baseUrl?.let { createEventsUrl(it) }
    }

    override fun execute(
        batchId: String,
        eventPackets: List<EventPacket>,
        attachmentPackets: List<AttachmentPacket>,
        spanPackets: List<SpanPacket>,
    ): HttpResponse {
        if (!isInitialized()) {
            // Handling this case as a HTTP response to make error handling consistent
            // with other network errors. This can only happen if the API_URL or API_KEY are
            // not correctly set.
            return HttpResponse.Error.UnknownError(UninitializedPropertyAccessException("Unable to initialize network client, please check the API_KEY and API_URL"))
        }

        val headers = createHeaders(batchId)
        val multipartData = prepareMultipartData(eventPackets, attachmentPackets, spanPackets)

        return try {
            val response =
                httpClient.sendMultipartRequest(eventsUrl.toString(), "PUT", headers, multipartData)
            handleResponse(response)
        } catch (e: Exception) {
            logger.log(LogLevel.Error, "Failed to send request", e)
            HttpResponse.Error.UnknownError(e)
        }
    }

    private fun parseUrl(url: String): URL? {
        return try {
            URL(url)
        } catch (e: Exception) {
            logger.log(LogLevel.Error, "Invalid API_URL", e)
            null
        }
    }

    private fun createEventsUrl(baseUrl: URL): URL? {
        return try {
            baseUrl.toURI().resolve(PATH_EVENTS).toURL()
        } catch (e: Exception) {
            logger.log(LogLevel.Error, "Invalid API_URL", e)
            null
        }
    }

    private fun isInitialized(): Boolean {
        return !(baseUrl == null || eventsUrl == null || apiKey == null)
    }

    private fun createHeaders(batchId: String): Map<String, String> {
        return mapOf(
            "msr-req-id" to batchId,
            "Authorization" to "Bearer $apiKey",
        )
    }

    private fun prepareMultipartData(
        eventPackets: List<EventPacket>,
        attachmentPackets: List<AttachmentPacket>,
        spanPackets: List<SpanPacket>,
    ): List<MultipartData> {
        val events = eventPackets.mapNotNull {
            multipartDataFactory.createFromEventPacket(it)
        }
        val attachments = attachmentPackets.mapNotNull {
            multipartDataFactory.createFromAttachmentPacket(it)
        }
        val spans = spanPackets.map {
            multipartDataFactory.createFromSpanPacket(it)
        }
        return events + attachments + spans
    }

    private fun handleResponse(response: HttpResponse): HttpResponse {
        return when (response) {
            is HttpResponse.Success -> {
                logger.log(LogLevel.Debug, "Request successful")
                response
            }

            is HttpResponse.Error.RateLimitError -> {
                logger.log(LogLevel.Debug, "Request rate limited, will retry later")
                response
            }

            is HttpResponse.Error.ClientError -> {
                logger.log(LogLevel.Error, "Unable to process request: ${response.code}")
                response
            }

            is HttpResponse.Error.ServerError -> {
                logger.log(LogLevel.Error, "Request failed with code: ${response.code}")
                response
            }

            is HttpResponse.Error.UnknownError -> {
                logger.log(LogLevel.Error, "Request failed with unknown error")
                response
            }
        }
    }
}
