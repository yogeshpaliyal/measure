package sh.measure.android.session

import kotlinx.serialization.Serializable

/**
 * A resource represents the entity producing telemetry. This typically contains device, OS, app
 * identifiers, etc.
 */
@Serializable
data class Resource(
    // device info
    val device_name: String? = null,
    val device_model: String? = null,
    val device_manufacturer: String? = null,
    val device_type: String? = null, // tablet, phone, tv, watch, etc.
    val device_is_foldable: Boolean? = null,
    val device_is_physical: Boolean? = null,
    val device_density_dpi: Int? = null,
    val device_width_px: Int? = null,
    val device_height_px: Int? = null,
    val device_density: Float? = null,
    // os info
    val os_name: String? = null,
    val os_version: String? = null,
    val platform: String? = null,
    // app info
    val app_version: String? = null,
    val app_build: String? = null,
    val app_unique_id: String? = null, // package name,
    val measure_sdk_version: String? = null,
)