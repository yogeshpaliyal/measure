package sh.measure.android.executors

import org.jetbrains.annotations.TestOnly
import java.util.concurrent.Executors
import java.util.concurrent.Future
import java.util.concurrent.ScheduledExecutorService
import java.util.concurrent.TimeUnit

internal interface MeasureExecutorService {
    val isClosed: Boolean
    fun submit(runnable: Runnable): Future<*>
    fun schedule(runnable: Runnable, delayMillis: Long): Future<*>
    fun scheduleAtFixedRate(
        runnable: Runnable,
        initialDelay: Long,
        delayMillis: Long,
        delayUnit: TimeUnit,
    ): Future<*>

    fun close(timeoutMillis: Long)
}

internal class MeasureExecutorServiceImpl @TestOnly constructor(private val executorService: ScheduledExecutorService) :
    MeasureExecutorService {

    constructor(customThreadFactory: CustomThreadFactory) : this(
        Executors.newSingleThreadScheduledExecutor(customThreadFactory),
    )

    override fun submit(runnable: Runnable): Future<*> {
        return executorService.submit(runnable)
    }

    override fun schedule(runnable: Runnable, delayMillis: Long): Future<*> {
        return executorService.schedule(runnable, delayMillis, TimeUnit.MILLISECONDS)
    }

    override fun scheduleAtFixedRate(
        runnable: Runnable,
        initialDelay: Long,
        delayMillis: Long,
        delayUnit: TimeUnit,
    ): Future<*> {
        return executorService.scheduleAtFixedRate(runnable, initialDelay, delayMillis, delayUnit)
    }

    override fun close(timeoutMillis: Long) {
        synchronized(executorService) {
            if (!executorService.isShutdown) {
                executorService.shutdown()
                try {
                    if (!executorService.awaitTermination(timeoutMillis, TimeUnit.MILLISECONDS)) {
                        executorService.shutdownNow()
                    }
                } catch (e: InterruptedException) {
                    executorService.shutdownNow()
                    Thread.currentThread().interrupt()
                }
            }
        }
    }

    override val isClosed: Boolean
        get() {
            synchronized(executorService) { return executorService.isShutdown }
        }
}
