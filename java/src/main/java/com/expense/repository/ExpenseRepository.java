package com.expense.repository;

import com.expense.models.Expense;
import java.sql.*;
import java.util.ArrayList;
import java.util.List;

public class ExpenseRepository {
  private Connection conn;

  public ExpenseRepository(Connection conn) {
    this.conn = conn;
  }

  public List<Expense> findAll(String search, String startDate, String endDate, Integer categoryId, Double minAmount,
      Double maxAmount) throws SQLException {
    List<Expense> list = new ArrayList<>();
    StringBuilder sql = new StringBuilder(
        "SELECT e.id, e.amount, e.expense_date, e.category_id, e.remarks, e.created_at, c.name as category_name " +
            "FROM expenses e JOIN categories c ON e.category_id = c.id WHERE 1=1");

    List<Object> params = new ArrayList<>();

    if (search != null && !search.isEmpty()) {
      sql.append(" AND (e.remarks ILIKE ? OR c.name ILIKE ?)");
      params.add("%" + search + "%");
      params.add("%" + search + "%");
    }
    if (startDate != null && !startDate.isEmpty()) {
      sql.append(" AND e.expense_date >= ?");
      params.add(Date.valueOf(startDate));
    }
    if (endDate != null && !endDate.isEmpty()) {
      sql.append(" AND e.expense_date <= ?");
      params.add(Date.valueOf(endDate));
    }
    if (categoryId != null && categoryId > 0) {
      sql.append(" AND e.category_id = ?");
      params.add(categoryId);
    }
    if (minAmount != null) {
      sql.append(" AND e.amount >= ?");
      params.add(minAmount);
    }
    if (maxAmount != null) {
      sql.append(" AND e.amount <= ?");
      params.add(maxAmount);
    }

    sql.append(" ORDER BY e.expense_date DESC, e.created_at DESC");

    try (PreparedStatement pstmt = conn.prepareStatement(sql.toString())) {
      for (int i = 0; i < params.size(); i++) {
        Object p = params.get(i);
        if (p instanceof String)
          pstmt.setString(i + 1, (String) p);
        else if (p instanceof Integer)
          pstmt.setInt(i + 1, (Integer) p);
        else if (p instanceof Double)
          pstmt.setDouble(i + 1, (Double) p);
        else if (p instanceof Date)
          pstmt.setDate(i + 1, (Date) p);
      }

      try (ResultSet rs = pstmt.executeQuery()) {
        while (rs.next()) {
          Expense e = new Expense();
          e.setId(rs.getInt("id"));
          e.setAmount(rs.getDouble("amount"));
          e.setExpense_date(rs.getDate("expense_date").toString());
          e.setCategory_id(rs.getInt("category_id"));
          e.setRemarks(rs.getString("remarks"));
          e.setCategory_name(rs.getString("category_name"));
          list.add(e);
        }
      }
    }
    return list;
  }

  public void create(Expense e) throws SQLException {
    String sql = "INSERT INTO expenses (amount, expense_date, category_id, remarks, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())";
    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setDouble(1, e.getAmount());
      pstmt.setDate(2, Date.valueOf(e.getExpense_date()));
      pstmt.setInt(3, e.getCategory_id());
      pstmt.setString(4, e.getRemarks());
      pstmt.executeUpdate();
    }
  }

  public void delete(int id) throws SQLException {
    String sql = "DELETE FROM expenses WHERE id = ?";
    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setInt(1, id);
      pstmt.executeUpdate();
    }
  }

  public double getYearlyTotal(int categoryId, int year) throws SQLException {
    String sql = "SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category_id = ? AND EXTRACT(YEAR FROM expense_date) = ?";
    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setInt(1, categoryId);
      pstmt.setInt(2, year);
      try (ResultSet rs = pstmt.executeQuery()) {
        if (rs.next())
          return rs.getDouble(1);
      }
    }
    return 0;
  }
}
