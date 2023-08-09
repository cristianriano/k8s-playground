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
  implementation("org.springframework.boot:spring-boot-starter-data-jpa")
  // Kotlin reflection is needed to instantiate Hibernate components
  implementation("org.jetbrains.kotlin:kotlin-reflect")
  implementation("mysql:mysql-connector-java:8.0.32")
//  implementation("com.fasterxml.jackson.module:jackson-module-kotlin")

  testImplementation("org.springframework.boot:spring-boot-starter-test")
  testImplementation("io.rest-assured:rest-assured")
  testImplementation("com.h2database:h2:2.1.214")
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

// Custom tasks examples
tasks.register<JavaExec>("runAsJavaExec") {
  group = "Execution"
  description = "Run the main class with JavaExecTask"
  mainClass.set(javaMainClass)
  classpath = sourceSets.main.get().runtimeClasspath
}

tasks.register<Exec>("runAsExec") {
  group = "Execution"
  description = "Run the main class with ExecTask"

  dependsOn("build")
  commandLine("java", "-classpath", sourceSets.main.get().runtimeClasspath.asPath, javaMainClass)
}
