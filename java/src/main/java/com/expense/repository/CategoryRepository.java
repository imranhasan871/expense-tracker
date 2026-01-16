package com.expense.repository;

import com.expense.models.Category;
import java.sql.*;
import java.util.ArrayList;
import java.util.List;

public class CategoryRepository {
  private Connection conn;

  public CategoryRepository(Connection conn) {
    this.conn = conn;
  }

  public List<Category> findAll() throws SQLException {
    List<Category> list = new ArrayList<>();
    String sql = "SELECT id, name, is_active, created_at FROM categories ORDER BY id";
    try (Statement stmt = conn.createStatement(); ResultSet rs = stmt.executeQuery(sql)) {
      while (rs.next()) {
        list.add(new Category(
            rs.getInt("id"),
            rs.getString("name"),
            rs.getBoolean("is_active"),
            rs.getTimestamp("created_at")));
      }
    }
    return list;
  }

  public void create(Category c) throws SQLException {
    String sql = "INSERT INTO categories (name, is_active, created_at, updated_at) VALUES (?, ?, NOW(), NOW())";
    try (PreparedStatement pstmt = conn.prepareStatement(sql)) {
      pstmt.setString(1, c.getName());
      pstmt.setBoolean(2, c.isIs_active());
      pstmt.executeUpdate();
    }
  }
}
