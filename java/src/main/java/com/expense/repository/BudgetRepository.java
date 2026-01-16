package com.expense.repository;

import com.expense.models.Budget;
import java.sql.*;
import java.util.ArrayList;
import java.util.List;

public class BudgetRepository {
  private Connection conn;

  public BudgetRepository(Connection conn) {
    this.conn = conn;
  }

  public List<Budget> findAll(int year) throws SQLException {
    List<Budget> list = new ArrayList<>();
    String sql = "SELECT b.id, b.category_id, b.amount, b.year, b.is_locked, c.name as category_name " +
        "FROM budgets b JOIN categories c ON b.category_id = c.id " +
        "WHERE b.year = ? ORDER BY c.name ASC";

    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setInt(1, year);
      try (ResultSet rs = pstmt.executeQuery()) {
        while (rs.next()) {
          Budget b = new Budget();
          b.setId(rs.getInt("id"));
          b.setCategory_id(rs.getInt("category_id"));
          b.setAmount(rs.getDouble("amount"));
          b.setYear(rs.getInt("year"));
          b.setIs_locked(rs.getBoolean("is_locked"));
          b.setCategory_name(rs.getString("category_name"));
          list.add(b);
        }
      }
    }
    return list;
  }

  public void createOrUpdate(Budget b) throws SQLException {
    // Upsert logic
    String sql = "INSERT INTO budgets (category_id, amount, year, created_at, updated_at) " +
        "VALUES (?, ?, ?, NOW(), NOW()) " +
        "ON CONFLICT (category_id, year) DO UPDATE " +
        "SET amount = EXCLUDED.amount, updated_at = NOW()";

    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setInt(1, b.getCategory_id());
      pstmt.setDouble(2, b.getAmount());
      pstmt.setInt(3, b.getYear());
      pstmt.executeUpdate();
    }
  }

  public Budget getByCategory(int categoryId, int year) throws SQLException {
    String sql = "SELECT * FROM budgets WHERE category_id = ? AND year = ?";
    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setInt(1, categoryId);
      pstmt.setInt(2, year);
      try (ResultSet rs = pstmt.executeQuery()) {
        if (rs.next()) {
          Budget b = new Budget();
          b.setId(rs.getInt("id"));
          b.setCategory_id(rs.getInt("category_id"));
          b.setAmount(rs.getDouble("amount"));
          b.setYear(rs.getInt("year"));
          b.setIs_locked(rs.getBoolean("is_locked"));
          return b;
        }
      }
    }
    return null;
  }

  public void toggleLock(int id, boolean isLocked) throws SQLException {
    String sql = "UPDATE budgets SET is_locked = ? WHERE id = ?";
    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setBoolean(1, isLocked);
      pstmt.setInt(2, id);
      pstmt.executeUpdate();
    }
  }

  // New: Monitoring Data (Complex Query)
  public List<Budget> getMonitoringData(int year) throws SQLException {
    List<Budget> list = new ArrayList<>();
    // Similar join query from Go implementation
    String sql = "SELECT b.id as budget_id, c.name as category_name, b.amount as budget_amount, " +
        "b.is_locked, COALESCE(SUM(e.amount), 0) as total_spent " +
        "FROM budgets b " +
        "JOIN categories c ON b.category_id = c.id " +
        "LEFT JOIN expenses e ON b.category_id = e.category_id " +
        "AND EXTRACT(YEAR FROM e.expense_date) = b.year " +
        "WHERE b.year = ? " +
        "GROUP BY b.id, c.name, b.amount, b.is_locked " +
        "ORDER BY c.name ASC";

    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setInt(1, year);
      try (ResultSet rs = pstmt.executeQuery()) {
        while (rs.next()) {
          Budget b = new Budget();
          b.setId(rs.getInt("budget_id"));
          b.setCategory_name(rs.getString("category_name"));
          b.setAmount(rs.getDouble("budget_amount"));
          b.setIs_locked(rs.getBoolean("is_locked"));
          b.setSpent(rs.getDouble("total_spent"));

          if (b.getAmount() > 0) {
            b.setPercentage((b.getSpent() / b.getAmount()) * 100);
          }
          list.add(b);
        }
      }
    }
    return list;
  }

  // Dashboard summary could go here too
  public double getTotalAnnualBudget(int year) throws SQLException {
    String sql = "SELECT COALESCE(SUM(amount), 0) FROM budgets WHERE year = ?";
    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setInt(1, year);
      try (ResultSet rs = pstmt.executeQuery()) {
        if (rs.next())
          return rs.getDouble(1);
      }
    }
    return 0;
  }
}
