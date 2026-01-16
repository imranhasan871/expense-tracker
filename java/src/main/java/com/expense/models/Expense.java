package com.expense.models;

import java.sql.Timestamp;

public class Expense {
  private int id;
  private double amount;
  private String expense_date; // Keeping as string to match existing JS format "YYYY-MM-DD", or use Date
  private int category_id;
  private String remarks;
  private Timestamp created_at;

  // Joined
  private String category_name;

  public Expense() {
  }

  public int getId() {
    return id;
  }

  public void setId(int id) {
    this.id = id;
  }

  public double getAmount() {
    return amount;
  }

  public void setAmount(double amount) {
    this.amount = amount;
  }

  public String getExpense_date() {
    return expense_date;
  }

  public void setExpense_date(String expense_date) {
    this.expense_date = expense_date;
  }

  public int getCategory_id() {
    return category_id;
  }

  public void setCategory_id(int category_id) {
    this.category_id = category_id;
  }

  public String getRemarks() {
    return remarks;
  }

  public void setRemarks(String remarks) {
    this.remarks = remarks;
  }

  public String getCategory_name() {
    return category_name;
  }

  public void setCategory_name(String category_name) {
    this.category_name = category_name;
  }
}
