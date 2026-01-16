package com.expense.handlers;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.net.URI;
import java.nio.file.Files;

public class StaticFileHandler implements HttpHandler {
    private static final String BASE_DIR = "src/main/resources/web";

    @Override
    public void handle(HttpExchange exchange) throws IOException {
        String path = exchange.getRequestURI().getPath();
        
        // Default to index.html
        if (path.equals("/")) {
            path = "/templates/index.html";
        } else if (!path.startsWith("/static/") && !path.startsWith("/templates/")) {
             // Handle "virtual" routes by serving the matching template if it exists
             // e.g. /categories -> /templates/categories.html
             String potentialTemplate = "/templates" + path + ".html";
             File f = new File(BASE_DIR + potentialTemplate);
             if (f.exists()) {
                 path = potentialTemplate;
             }
        }

        File file = new File(BASE_DIR + path);
        if (!file.exists()) {
            exchange.sendResponseHeaders(404, 0);
            try (OutputStream os = exchange.getResponseBody()) {
                os.write("404 Not Found".getBytes());
            }
            return;
        }

        String mimeType = Files.probeContentType(file.toPath());
        if (path.endsWith(".css")) mimeType = "text/css";
        if (path.endsWith(".js")) mimeType = "application/javascript";
        
        if (mimeType != null) {
            exchange.getResponseHeaders().set("Content-Type", mimeType);
        }

        exchange.sendResponseHeaders(200, file.length());
        try (OutputStream os = exchange.getResponseBody(); FileInputStream fs = new FileInputStream(file)) {
            final byte[] buffer = new byte[0x10000];
            int count = 0;
            while ((count = fs.read(buffer)) >= 0) {
                os.write(buffer, 0, count);
            }
        }
    }
}
