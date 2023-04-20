/*
 * This file was generated by the Gradle 'init' task.
 *
 * This generated file contains a sample Kotlin application project to get you started.
 * For more details take a look at the 'Building Java & JVM projects' chapter in the Gradle
 * User Manual available at https://docs.gradle.org/7.4.2/userguide/building_java_projects.html
 * This project uses @Incubating APIs which are subject to change.
 */
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.springframework.boot.gradle.tasks.bundling.BootJar

plugins {
  kotlin("jvm") version "1.8.10"
  id("application")
  id("idea")

  id("org.springframework.boot") version "3.0.5"
  id("io.spring.dependency-management") version "1.1.0"
}

tasks.withType<KotlinCompile> {
  kotlinOptions {
    // Enforce null safety with Spring
    freeCompilerArgs = listOf("-Xjsr305=strict")
    jvmTarget = "17"
  }
}

repositories {
  // Use Maven Central for resolving dependencies.
  mavenCentral()
}

dependencies {
  // Align versions of all Kotlin components
  implementation(platform("org.jetbrains.kotlin:kotlin-bom"))
  // Use the Kotlin JDK 8 standard library.
  implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")

  // Spring
  implementation("org.springframework.boot:spring-boot-starter-web")
//  implementation("com.fasterxml.jackson.module:jackson-module-kotlin")

  testImplementation("org.springframework.boot:spring-boot-starter-test")
  testImplementation("io.rest-assured:rest-assured")
  testImplementation(kotlin("test"))
}

tasks.test { // See 5️⃣
  useJUnitPlatform() // JUnitPlatform for tests. See 6️⃣
}

val bootJar: BootJar by tasks
bootJar.apply {
  // Set name of the generated jar
  archiveBaseName.set("app")
}

val javaMainClass = "com.example.shoppinglist.AppKt"
application {
  // Define the main class for the application.
  mainClass.set(javaMainClass)
}
