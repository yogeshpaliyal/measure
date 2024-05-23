package sh.measure.android.exporter

import androidx.concurrent.futures.ResolvableFuture
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import org.junit.Assert
import org.junit.Test
import org.junit.runner.RunWith
import org.mockito.Mockito
import org.mockito.kotlin.any
import org.mockito.kotlin.mock
import sh.measure.android.events.EventType
import sh.measure.android.fakes.FakeEventFactory
import sh.measure.android.fakes.FakeEventFactory.toEvent
import sh.measure.android.fakes.FakeIdProvider
import sh.measure.android.fakes.FakeTimeProvider
import sh.measure.android.fakes.ImmediateExecutorService
import sh.measure.android.fakes.NoopLogger
import sh.measure.android.storage.AttachmentEntity
import sh.measure.android.storage.DatabaseImpl
import sh.measure.android.storage.FileStorage

@RunWith(AndroidJUnit4::class)
class EventExporterTest {

    private val logger = NoopLogger()
    private val database =
        DatabaseImpl(InstrumentationRegistry.getInstrumentation().context, logger)
    private val networkClient = mock<NetworkClient>()
    private val idProvider = FakeIdProvider()
    private val timeProvider = FakeTimeProvider()
    private val fileStorage = mock<FileStorage>()
    private val executorService = ImmediateExecutorService(ResolvableFuture.create<Any>())

    private val eventExporter = EventExporterImpl(
        logger = logger,
        database = database,
        networkClient = networkClient,
        idProvider = idProvider,
        timeProvider = timeProvider,
        fileStorage = fileStorage,
        executorService = executorService,
    )

    @Test
    fun `given a export succeeds, deletes the events and the batch from database`() {
        Mockito.`when`(networkClient.execute(any(), any(), any())).thenReturn(true)
        val exceptionData = FakeEventFactory.getExceptionData()
        val event = exceptionData.toEvent(type = EventType.EXCEPTION)
        database.insertEvent(FakeEventFactory.fakeEventEntity(eventId = event.id))

        eventExporter.export(event)

        Assert.assertEquals(0, database.getEventsCount())
        Assert.assertEquals(0, database.getBatchesCount())
    }

    @Test
    fun `given export succeeds, deletes the events from file storage`() {
        Mockito.`when`(networkClient.execute(any(), any(), any())).thenReturn(true)
        val exceptionData = FakeEventFactory.getExceptionData()
        val event = exceptionData.toEvent(id = "event-id", type = EventType.EXCEPTION)
        database.insertEvent(
            FakeEventFactory.fakeEventEntity(
                eventId = event.id,
                attachmentEntities = listOf(
                    AttachmentEntity("attachment-id", type = "type", path = "path", name = "name"),
                ),
            ),
        )

        eventExporter.export(event)

        Mockito.verify(fileStorage).deleteEventIfExist(event.id, listOf("attachment-id"))
    }

    @Test
    fun `given export fails, event and the batch are not deleted`() {
        Mockito.`when`(networkClient.execute(any(), any(), any())).thenReturn(false)
        val exceptionData = FakeEventFactory.getExceptionData()
        val event = exceptionData.toEvent(type = EventType.EXCEPTION)
        database.insertEvent(FakeEventFactory.fakeEventEntity(eventId = event.id))

        eventExporter.export(event)

        Assert.assertEquals(1, database.getEventsCount())
        Assert.assertEquals(1, database.getBatchesCount())
    }
}