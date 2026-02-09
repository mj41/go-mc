///usr/bin/env java --source 21 "$0" "$@"; exit $?
// ^^^ allows running as: java ExtractAll.java (JEP 458, JDK 21+)

import java.io.*;
import java.net.*;
import java.net.http.*;
import java.nio.file.*;
import java.util.*;
import java.util.regex.*;
import java.util.zip.*;

/**
 * ExtractAll — Downloads an unobfuscated MC server jar and extracts all data as JSON.
 *
 * Runs inside a Docker/Podman container with these mounted volumes:
 *   /cache      — persistent cache for downloaded jars
 *   /jsons      — output directory for JSON files
 *   /extractors — this file + custom Java extractors (read-only)
 *
 * Environment:
 *   MC_VERSION  — Minecraft version to extract (e.g., "1.21.11")
 *
 * Steps:
 *   1. Download unobfuscated server jar (cached)
 *   2. Extract inner jar from bundler format
 *   3. Run MC --all data generator → reports/*.json
 *   4. Compile and run custom extractors (GenEntities, GenComponents, GenBlockEntities)
 *   5. Copy all JSON to /jsons/<version>/
 */
public class ExtractAll {

    static final String MANIFEST_URL = "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json";
    static final String PISTON_DATA_URL = "https://piston-data.mojang.com/v1/objects";
    static final String WIKI_API_URL = "https://minecraft.wiki/api.php";

    static final HttpClient HTTP = HttpClient.newHttpClient();

    static String version;
    static Path cacheDir  = Path.of("/cache");
    static Path jsonsDir  = Path.of("/jsons");
    static Path extractorsDir = Path.of("/extractors");

    public static void main(String[] args) throws Exception {
        version = System.getenv("MC_VERSION");
        if (version == null || version.isEmpty()) {
            if (args.length > 0) {
                version = args[0];
            } else {
                fatal("MC_VERSION environment variable or first argument required");
            }
        }

        // Allow running outside container with custom paths.
        String cacheEnv = System.getenv("MC_CACHE_DIR");
        String jsonsEnv = System.getenv("MC_JSONS_DIR");
        String extractorsEnv = System.getenv("MC_EXTRACTORS_DIR");
        if (cacheEnv != null) cacheDir = Path.of(cacheEnv);
        if (jsonsEnv != null) jsonsDir = Path.of(jsonsEnv);
        if (extractorsEnv != null) extractorsDir = Path.of(extractorsEnv);

        log("ExtractAll: version=%s", version);
        log("  cache:      %s", cacheDir);
        log("  jsons:      %s", jsonsDir);
        log("  extractors: %s", extractorsDir);

        // Step 1: Download unobfuscated server jar.
        Path serverJar = downloadServerJar();

        // Step 2: Extract inner jar from bundler format.
        Path innerJar = extractInnerJar(serverJar);

        // Step 3: Run --all data generator.
        Path reportsDir = runDataGenerator(serverJar);

        // Step 4: Copy --all reports to output.
        Path outputDir = jsonsDir.resolve(version);
        Files.createDirectories(outputDir);
        copyReports(reportsDir, outputDir);

        // Step 5: Run custom extractors (if present).
        runCustomExtractors(serverJar, innerJar, outputDir);

        // Done.
        log("");
        log("Extraction complete. Output:");
        try (var stream = Files.list(outputDir)) {
            stream.sorted().forEach(p -> {
                try {
                    long size = Files.size(p);
                    log("  %-30s %s", p.getFileName(), humanSize(size));
                } catch (IOException e) {
                    log("  %s (error: %s)", p.getFileName(), e.getMessage());
                }
            });
        }
    }

    // --- Step 1: Download ---

    static Path downloadServerJar() throws Exception {
        Path jarPath;

        if (isNativelyUnobfuscated(version)) {
            log("");
            log("Version %s is natively unobfuscated (26.1+), using standard download.", version);
            jarPath = cacheDir.resolve(version + "-server.jar");
            if (Files.exists(jarPath) && Files.size(jarPath) > 0) {
                log("  Cached: %s (%s)", jarPath, humanSize(Files.size(jarPath)));
                return jarPath;
            }
            downloadStandardJar(jarPath);
        } else {
            log("");
            log("Looking up unobfuscated server download for %s via Minecraft Wiki...", version);
            jarPath = cacheDir.resolve(version + "-unobf-server.jar");
            if (Files.exists(jarPath) && Files.size(jarPath) > 0) {
                log("  Cached: %s (%s)", jarPath, humanSize(Files.size(jarPath)));
                return jarPath;
            }

            String hash = findUnobfuscatedHash(version);
            if (hash != null) {
                String url = PISTON_DATA_URL + "/" + hash + "/server.jar";
                log("  Downloading unobfuscated jar (hash: %s...)...", hash.substring(0, 12));
                downloadFile(url, jarPath);
            } else {
                log("  WARNING: unobfuscated download not found, falling back to standard jar.");
                jarPath = cacheDir.resolve(version + "-server.jar");
                if (Files.exists(jarPath) && Files.size(jarPath) > 0) {
                    return jarPath;
                }
                downloadStandardJar(jarPath);
            }
        }
        return jarPath;
    }

    static void downloadStandardJar(Path jarPath) throws Exception {
        // Fetch manifest → find version → get download URL.
        String manifestJson = httpGet(MANIFEST_URL);
        String versionUrl = parseVersionUrl(manifestJson, version);
        if (versionUrl == null) {
            fatal("Version %s not found in Mojang manifest", version);
        }

        String detailJson = httpGet(versionUrl);
        String serverUrl = parseJsonString(detailJson, "\"server\"", "\"url\"");
        if (serverUrl == null) {
            fatal("No server download URL for version %s", version);
        }

        downloadFile(serverUrl, jarPath);
    }

    static String findUnobfuscatedHash(String version) throws Exception {
        String pageTitle = "Java Edition " + version;
        String apiUrl = WIKI_API_URL +
            "?action=parse&page=" + URLEncoder.encode(pageTitle, "UTF-8") +
            "&prop=wikitext&format=json&section=0";

        HttpRequest req = HttpRequest.newBuilder()
            .uri(URI.create(apiUrl))
            .header("User-Agent", "ExtractAll/1.0 (w42-mc-cubes; Java)")
            .GET().build();

        HttpResponse<String> resp = HTTP.send(req, HttpResponse.BodyHandlers.ofString());
        if (resp.statusCode() != 200) {
            return null;
        }

        // Look for {{dl|HASH|server|title=Unobfuscated}}.
        Pattern p = Pattern.compile("\\{\\{dl\\|([0-9a-f]{40})\\|server\\|title=Unobfuscated\\}\\}");
        Matcher m = p.matcher(resp.body());
        return m.find() ? m.group(1) : null;
    }

    static boolean isNativelyUnobfuscated(String version) {
        // 26.x+ versions are natively unobfuscated.
        String[] parts = version.split("\\.", 2);
        try {
            int major = Integer.parseInt(parts[0]);
            return major >= 26;
        } catch (NumberFormatException e) {
            return false;
        }
    }

    // --- Step 2: Extract inner jar ---

    static Path extractInnerJar(Path bundledJar) throws Exception {
        Path innerPath = cacheDir.resolve(version + "-inner.jar");
        if (Files.exists(innerPath) && Files.size(innerPath) > 0) {
            log("  Inner jar cached: %s (%s)", innerPath, humanSize(Files.size(innerPath)));
            return innerPath;
        }

        log("");
        log("Extracting inner server jar from bundler format...");

        try (ZipFile zip = new ZipFile(bundledJar.toFile())) {
            // Look for META-INF/versions/*/server-*.jar.
            ZipEntry serverEntry = null;
            for (var entries = zip.entries(); entries.hasMoreElements(); ) {
                ZipEntry e = entries.nextElement();
                if (e.getName().startsWith("META-INF/versions/") && e.getName().endsWith(".jar")) {
                    serverEntry = e;
                    break;
                }
            }

            if (serverEntry == null) {
                // Not bundled — just copy.
                log("  Not bundled, using jar directly.");
                Files.copy(bundledJar, innerPath, StandardCopyOption.REPLACE_EXISTING);
                return innerPath;
            }

            log("  Found: %s (%s)", serverEntry.getName(), humanSize(serverEntry.getSize()));
            try (InputStream in = zip.getInputStream(serverEntry)) {
                Files.copy(in, innerPath, StandardCopyOption.REPLACE_EXISTING);
            }
        }

        log("  Extracted: %s", humanSize(Files.size(innerPath)));
        return innerPath;
    }

    // --- Step 3: Run --all data generator ---

    static Path runDataGenerator(Path serverJar) throws Exception {
        // Create a temp working dir for the generator (it creates files in cwd).
        Path workDir = cacheDir.resolve(version + "-datagen");
        Files.createDirectories(workDir);

        Path reportsDir = workDir.resolve("generated").resolve("reports");
        if (Files.exists(reportsDir.resolve("registries.json"))) {
            log("");
            log("Data generator output cached, skipping --all.");
            return reportsDir;
        }

        log("");
        log("Running MC --all data generator...");

        ProcessBuilder pb = new ProcessBuilder(
            "java",
            "-DbundlerMainClass=net.minecraft.data.Main",
            "-jar", serverJar.toAbsolutePath().toString(),
            "--all"
        );
        pb.directory(workDir.toFile());
        pb.inheritIO();

        Process proc = pb.start();
        int exit = proc.waitFor();
        if (exit != 0) {
            fatal("Data generator failed with exit code %d", exit);
        }

        if (!Files.exists(reportsDir.resolve("registries.json"))) {
            fatal("Data generator did not produce registries.json");
        }

        log("  Data generator complete.");
        return reportsDir;
    }

    // --- Step 4: Copy reports ---

    static void copyReports(Path reportsDir, Path outputDir) throws Exception {
        log("");
        log("Copying reports to %s ...", outputDir);

        String[] files = {"registries.json", "blocks.json", "items.json", "packets.json",
                          "commands.json", "datapack.json"};

        for (String name : files) {
            Path src = reportsDir.resolve(name);
            if (Files.exists(src)) {
                Files.copy(src, outputDir.resolve(name), StandardCopyOption.REPLACE_EXISTING);
                log("  %-25s %s", name, humanSize(Files.size(src)));
            }
        }

        // Also copy json-rpc-api-schema.json if present (1.21.11+).
        Path schema = reportsDir.resolve("json-rpc-api-schema.json");
        if (Files.exists(schema)) {
            Files.copy(schema, outputDir.resolve("json-rpc-api-schema.json"),
                       StandardCopyOption.REPLACE_EXISTING);
        }
    }

    // --- Step 5: Custom extractors ---

    static void runCustomExtractors(Path serverJar, Path innerJar, Path outputDir) throws Exception {
        String[] extractors = {"GenEntities", "GenComponents", "GenBlockEntities", "GenBlockProperties"};
        List<String> found = new ArrayList<>();

        for (String name : extractors) {
            Path src = extractorsDir.resolve(name + ".java");
            if (Files.exists(src)) {
                found.add(name);
            }
        }

        if (found.isEmpty()) {
            log("");
            log("No custom extractors found in %s (optional, skipping).", extractorsDir);
            return;
        }

        // Extract library jars from the bundled server jar.
        String classpath = buildExtractorClasspath(serverJar, innerJar);

        log("");
        log("Compiling custom extractors: %s", String.join(", ", found));

        // Compile all found extractors.
        List<String> javacArgs = new ArrayList<>(List.of(
            "javac", "-cp", classpath, "-proc:none",
            "-d", outputDir.toAbsolutePath().toString()
        ));
        for (String name : found) {
            javacArgs.add(extractorsDir.resolve(name + ".java").toAbsolutePath().toString());
        }

        int compileExit = new ProcessBuilder(javacArgs)
            .inheritIO().start().waitFor();
        if (compileExit != 0) {
            log("  WARNING: extractor compilation failed (exit %d). Skipping.", compileExit);
            return;
        }

        // Run each extractor.
        String runCp = classpath + ":" + outputDir.toAbsolutePath();
        for (String name : found) {
            log("  Running %s...", name);
            ProcessBuilder pb = new ProcessBuilder(
                "java", "-cp", runCp, name
            );
            pb.directory(outputDir.toFile());
            pb.inheritIO();

            int exit = pb.start().waitFor();
            if (exit != 0) {
                log("  WARNING: %s failed with exit code %d", name, exit);
            }
        }
    }

    /**
     * Extract library jars from the bundled server jar and build a classpath string.
     * The bundled jar contains META-INF/libraries/*.jar and META-INF/classpath-joined.
     */
    static String buildExtractorClasspath(Path serverJar, Path innerJar) throws Exception {
        Path libsDir = cacheDir.resolve(version + "-libs");
        Path doneMarker = libsDir.resolve(".extracted");

        if (!Files.exists(doneMarker)) {
            log("Extracting library jars from bundled server...");
            Files.createDirectories(libsDir);

            try (ZipFile zip = new ZipFile(serverJar.toFile())) {
                for (var entries = zip.entries(); entries.hasMoreElements(); ) {
                    ZipEntry e = entries.nextElement();
                    if (e.getName().startsWith("META-INF/libraries/") && e.getName().endsWith(".jar")) {
                        Path dest = libsDir.resolve(e.getName());
                        Files.createDirectories(dest.getParent());
                        try (InputStream in = zip.getInputStream(e)) {
                            Files.copy(in, dest, StandardCopyOption.REPLACE_EXISTING);
                        }
                    }
                }
            }
            Files.writeString(doneMarker, "ok");
            log("  Library jars extracted.");
        }

        // Build classpath: inner jar + all library jars.
        StringBuilder cp = new StringBuilder();
        cp.append(innerJar.toAbsolutePath());

        try (var stream = Files.walk(libsDir)) {
            stream.filter(p -> p.toString().endsWith(".jar")).forEach(p -> {
                cp.append(':').append(p.toAbsolutePath());
            });
        }

        return cp.toString();
    }

    // --- HTTP helpers ---

    static String httpGet(String url) throws Exception {
        HttpRequest req = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .header("User-Agent", "ExtractAll/1.0")
            .GET().build();
        HttpResponse<String> resp = HTTP.send(req, HttpResponse.BodyHandlers.ofString());
        if (resp.statusCode() != 200) {
            fatal("HTTP %d from %s", resp.statusCode(), url);
        }
        return resp.body();
    }

    static void downloadFile(String url, Path dest) throws Exception {
        Files.createDirectories(dest.getParent());
        HttpRequest req = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .header("User-Agent", "ExtractAll/1.0")
            .GET().build();

        long start = System.currentTimeMillis();
        HttpResponse<Path> resp = HTTP.send(req,
            HttpResponse.BodyHandlers.ofFile(dest));

        if (resp.statusCode() != 200) {
            Files.deleteIfExists(dest);
            fatal("HTTP %d downloading %s", resp.statusCode(), url);
        }

        long elapsed = System.currentTimeMillis() - start;
        log("  Downloaded %s in %dms → %s", humanSize(Files.size(dest)), elapsed, dest);
    }

    // --- JSON parsing (minimal, no external deps) ---

    static String parseVersionUrl(String manifestJson, String version) {
        // Find {"id":"<version>", ... "url":"<url>"}
        String needle = "\"id\":\"" + version + "\"";
        int idx = manifestJson.indexOf(needle);
        if (idx < 0) return null;

        // Find the "url" field in the same object.
        int urlIdx = manifestJson.indexOf("\"url\":\"", idx);
        if (urlIdx < 0) return null;
        int start = urlIdx + 7;
        int end = manifestJson.indexOf("\"", start);
        return manifestJson.substring(start, end);
    }

    static String parseJsonString(String json, String... keys) {
        int pos = 0;
        for (String key : keys) {
            pos = json.indexOf(key, pos);
            if (pos < 0) return null;
            pos += key.length();
        }
        // Find the next quoted string value.
        int q1 = json.indexOf("\"", pos);
        if (q1 < 0) return null;
        int q2 = json.indexOf("\"", q1 + 1);
        if (q2 < 0) return null;
        return json.substring(q1 + 1, q2);
    }

    // --- Utility ---

    static String humanSize(long bytes) {
        if (bytes >= 1024 * 1024) return String.format("%.1f MB", bytes / (1024.0 * 1024.0));
        if (bytes >= 1024) return String.format("%.1f KB", bytes / 1024.0);
        return bytes + " B";
    }

    static void log(String fmt, Object... args) {
        System.out.printf(fmt + "%n", args);
    }

    static void fatal(String fmt, Object... args) {
        System.err.printf("ExtractAll: " + fmt + "%n", args);
        System.exit(1);
    }
}
