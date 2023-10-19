package sh.measure.android.fakes

import sh.measure.android.logger.LogLevel
import sh.measure.android.logger.Logger

internal class NoopLogger : Logger  {
    override fun log(level: LogLevel, message: String, throwable: Throwable?) {
        // No-op
    }
}