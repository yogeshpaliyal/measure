package sh.measure.android.utils

import android.os.Debug

internal interface DebugProvider {
    fun getNativeHeapSize(): Long
    fun getNativeHeapFreeSize(): Long
    fun getMemoryInfo(memoryInfo: Debug.MemoryInfo)
}

internal class DefaultDebugProvider : DebugProvider {
    override fun getNativeHeapSize(): Long {
        return Debug.getNativeHeapSize()
    }

    override fun getNativeHeapFreeSize(): Long {
        return Debug.getNativeHeapFreeSize()
    }

    override fun getMemoryInfo(memoryInfo: Debug.MemoryInfo) {
        Debug.getMemoryInfo(memoryInfo)
    }
}
