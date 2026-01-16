package com.expense;

import com.expense.config.DatabaseConfig;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.net.InetSocketAddress;
import java.util.concurrent.Executors;

public class Main {
    public static void main(String[] args) throws IOException {
        int port = 8080;
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        
        System.out.println("Starting server on port " + port);

        // Serve Static Files
        server.createContext("/", new com.expense.handlers.StaticFileHandler());
        
        // API Routes
        server.createContext("/api/categories", new com.expense.handlers.CategoryHandler());
        // Add other context handlers here

        server.setExecutor(Executors.newCachedThreadPool());
        server.start();
        System.out.println("Server started.");
    }
}
