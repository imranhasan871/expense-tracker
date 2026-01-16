package com.expense.config;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

public class DatabaseConfig {
    // Hardcoding for "raw" simplicity, matching docker-compose
    private static final String URL = "jdbc:postgresql://localhost:5432/expense_tracker";
    private static final String USER = "admin";
    private static final String PASSWORD = "root";

    public static Connection getConnection() throws SQLException {
        return DriverManager.getConnection(URL, USER, PASSWORD);
    }
}
