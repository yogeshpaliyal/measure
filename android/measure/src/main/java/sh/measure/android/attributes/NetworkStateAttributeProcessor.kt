package sh.measure.android.attributes

import sh.measure.android.networkchange.NetworkGeneration
import sh.measure.android.networkchange.NetworkProvider
import sh.measure.android.networkchange.NetworkStateProvider
import sh.measure.android.networkchange.NetworkType

/**
 * Generates the network state attributes. These attributes are expected to change during the
 * session. This class computes the attributes every time [appendAttributes] is called.
 */
internal class NetworkStateAttributeProcessor(
    private val networkStateProvider: NetworkStateProvider,
) : AttributeProcessor {
    private var networkType: String = NetworkType.UNKNOWN
    private var networkGeneration: String = NetworkGeneration.UNKNOWN
    private var networkProvider: String = NetworkProvider.UNKNOWN

    override fun appendAttributes(attributes: MutableMap<String, Any?>) {
        computeAttributes()
        attributes.apply {
            put(Attribute.NETWORK_TYPE_KEY, networkType)
            put(Attribute.NETWORK_GENERATION_KEY, networkGeneration)
            put(Attribute.NETWORK_PROVIDER_KEY, networkProvider)
        }
    }

    private fun computeAttributes() {
        val networkState = networkStateProvider.getNetworkState()
        networkType = networkState?.networkType ?: NetworkType.UNKNOWN
        networkGeneration = networkState?.networkGeneration ?: NetworkGeneration.UNKNOWN
        networkProvider = networkState?.networkProvider ?: NetworkProvider.UNKNOWN
    }
}
