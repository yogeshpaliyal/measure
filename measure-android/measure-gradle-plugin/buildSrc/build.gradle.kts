plugins {
    id("org.gradle.kotlin.kotlin-dsl") version "4.1.2"
    id("java-gradle-plugin")
    id("org.jetbrains.kotlin.jvm") version "1.9.10"
}

repositories {
    mavenCentral()
}

dependencies {
    gradleApi()
}

gradlePlugin {
    plugins {
        register("sh.measure.plugin.aar2jar") {
            id = "sh.measure.plugin.aar2jar"
            implementationClass = "sh.measure.plugin.gradle.Aar2Jar"
        }
    }
}