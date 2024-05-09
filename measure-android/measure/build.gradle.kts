@file:Suppress("UnstableApiUsage")

import com.diffplug.gradle.spotless.SpotlessExtension

plugins {
    alias(libs.plugins.android.library)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.kotlin.serialization)
    alias(libs.plugins.kotlinx.binary.compatibility.validator)
    alias(libs.plugins.diffplug.spotless)
    id("maven-publish")
}

val measureSdkVersion = properties["MEASURE_VERSION_NAME"] as String
publishing {
    publications {
        create<MavenPublication>("maven") {
            groupId = properties["GROUP"] as String
            artifactId = properties["MEASURE_ARTIFACT_ID"] as String
            version = measureSdkVersion

            afterEvaluate {
                from(components["release"])
            }
        }
    }

    repositories {
        maven {
            name = "GitHubPackages"
            url = uri("https://maven.pkg.github.com/measure-sh/measure")
            credentials {
                username = System.getenv("GITHUB_ACTOR")
                password = System.getenv("GITHUB_TOKEN")
            }
        }
    }
}

android {
    namespace = "sh.measure.android"
    compileSdk = 34

    defaultConfig {
        minSdk = 21

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        consumerProguardFiles("consumer-rules.pro")
    }

    buildTypes {
        defaultConfig {
            manifestPlaceholders["measure_url"] = properties["measure_url"]?.toString() ?: ""
            buildConfigField("String", "MEASURE_SDK_VERSION", "\"$measureSdkVersion\"")
        }
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro",
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    kotlinOptions {
        jvmTarget = "11"
    }
    testOptions {
        unitTests {
            isIncludeAndroidResources = true
            isReturnDefaultValues = true
        }
    }
    buildFeatures {
        buildConfig = true
        compose = true
    }
    composeOptions {
        kotlinCompilerExtensionVersion = "1.5.7"
    }
}

extensions.configure<SpotlessExtension>("spotless") {
    plugins.withId("org.jetbrains.kotlin.jvm") {
        configureSpotlessKotlin(this@configure)
    }
    plugins.withId("org.jetbrains.kotlin.android") {
        configureSpotlessKotlin(this@configure)
    }
    kotlinGradle {
        ktlint()
    }
    format("misc") {
        target(
            ".gitignore",
            ".gitattributes",
            ".gitconfig",
            ".editorconfig",
            "*.md",
            "src/**/*.md",
            "docs/**/*.md",
            "src/**/*.properties",
        )
        indentWithSpaces()
        trimTrailingWhitespace()
        endWithNewline()
    }
}

fun configureSpotlessKotlin(spotlessExtension: SpotlessExtension) {
    spotlessExtension.kotlin {
        ktlint().apply {
            editorConfigOverride(
                mapOf("max_line_length" to 2147483647),
            )
        }
        target("src/**/*.kt")
    }
}

dependencies {
    // Compile only, as we don't want to include the fragment dependency in the final artifact.
    compileOnly(libs.androidx.fragment.ktx)
    compileOnly(libs.androidx.compose.runtime.android)
    compileOnly(libs.androidx.compose.ui)
    compileOnly(libs.androidx.navigation.compose)

    implementation(libs.kotlinx.serialization.json)

    implementation(libs.androidx.annotation)
    implementation(libs.squareup.okio)
    implementation(libs.squareup.okhttp)
    implementation(libs.squareup.okhttp.logging)
    implementation(libs.squareup.curtains)
    implementation(project(":measure-ndk"))

    testImplementation(libs.mockito.kotlin)
    testImplementation(libs.junit)
    testImplementation(libs.androidx.junit.ktx)
    testImplementation(libs.robolectric)
    testImplementation(libs.androidx.fragment.testing)
    testImplementation(libs.androidx.rules)
    testImplementation(libs.androidx.compose.runtime.android)
    testImplementation(libs.squareup.okhttp.mockwebserver)

    debugImplementation("androidx.compose.ui:ui-test-manifest:1.4.3") {
        isTransitive = false
    }
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(libs.androidx.compose.runtime.android)
    androidTestImplementation(libs.androidx.compose.ui)
    androidTestImplementation(libs.androidx.compose.ui.test.junit4)
    androidTestImplementation(libs.androidx.material3)
    androidTestImplementation(libs.androidx.lifecycle.process)
    androidTestImplementation(libs.androidx.lifecycle.common)
    androidTestImplementation(libs.androidx.activity.compose)
    androidTestImplementation(libs.androidx.navigation.compose)
    androidTestImplementation(libs.androidx.rules)
}
