package com.expense;

import com.expense.config.DatabaseConfig;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.net.InetSocketAddress;
import java.util.concurrent.Executors;

public class Main {
    public static void main(String[] args) throws IOException {
        int port = 3000;
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);

        System.out.println("Starting server on port " + port);

        // Serve Static Files
        server.createContext("/", new com.expense.handlers.StaticFileHandler());

        // API Routes
        server.createContext("/api/categories", new com.expense.handlers.CategoryHandler());
        server.createContext("/api/budgets", new com.expense.handlers.BudgetHandler());
        server.createContext("/api/monitoring", new com.expense.handlers.BudgetHandler()); // For direct monitoring call
                                                                                           // if needed, handled in
                                                                                           // BudgetHandler
        server.createContext("/api/expenses", new com.expense.handlers.ExpenseHandler());

        server.setExecutor(Executors.newCachedThreadPool());
        server.start();
        System.out.println("Server started.");
    }
}
