package com.expense.models;

import java.sql.Timestamp;

public class Category {
  private int id;
  private String name;
  private boolean is_active;
  private Timestamp created_at;
  private Timestamp updated_at;

  public Category() {
  }

  public Category(int id, String name, boolean is_active, Timestamp created_at) {
    this.id = id;
    this.name = name;
    this.is_active = is_active;
    this.created_at = created_at;
  }

  public int getId() {
    return id;
  }

  public String getName() {
    return name;
  }

  public boolean isIs_active() {
    return is_active;
  } // Gson convention might require manual hook or specific naming

  public void setName(String name) {
    this.name = name;
  }

  public void setIs_active(boolean is_active) {
    this.is_active = is_active;
  }
}
