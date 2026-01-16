package com.expense.handlers;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import java.io.IOException;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import com.google.gson.Gson;

public abstract class BaseHandler implements HttpHandler {
    protected Gson gson = new Gson();

    protected void sendResponse(HttpExchange exchange, int statusCode, Object response) throws IOException {
        String jsonResponse = gson.toJson(response);
        byte[] bytes = jsonResponse.getBytes(StandardCharsets.UTF_8);
        
        exchange.getResponseHeaders().set("Content-Type", "application/json");
        exchange.sendResponseHeaders(statusCode, bytes.length);
        
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }

    protected void sendError(HttpExchange exchange, int statusCode, String message) throws IOException {
        sendResponse(exchange, statusCode, new ErrorResponse(message));
    }
    
    // Helper classes for standard response format
    static class ErrorResponse {
        boolean success = false;
        String message;
        public ErrorResponse(String message) { this.message = message; }
    }
    
    static class SuccessResponse {
        boolean success = true;
        Object data;
        String message;
        public SuccessResponse(Object data) { this.data = data; }
        public SuccessResponse(Object data, String message) { this.data = data; this.message = message; }
    }
}
