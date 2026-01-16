package com.expense.models;

import java.sql.Timestamp;

public class Budget {
  private int id;
  private int category_id;
  private double amount;
  private int year;
  private boolean is_locked; // Added for circuit breaker
  private Timestamp created_at;
  private Timestamp updated_at;

  // Joined fields
  private String category_name;

  // For monitoring stats
  private double percentage;
  private double spent;

  public Budget() {
  }

  // Getters and Setters
  public int getId() {
    return id;
  }

  public void setId(int id) {
    this.id = id;
  }

  public int getCategory_id() {
    return category_id;
  }

  public void setCategory_id(int category_id) {
    this.category_id = category_id;
  }

  public double getAmount() {
    return amount;
  }

  public void setAmount(double amount) {
    this.amount = amount;
  }

  public int getYear() {
    return year;
  }

  public void setYear(int year) {
    this.year = year;
  }

  public boolean isIs_locked() {
    return is_locked;
  }

  public void setIs_locked(boolean is_locked) {
    this.is_locked = is_locked;
  }

  public String getCategory_name() {
    return category_name;
  }

  public void setCategory_name(String category_name) {
    this.category_name = category_name;
  }

  public double getPercentage() {
    return percentage;
  }

  public void setPercentage(double percentage) {
    this.percentage = percentage;
  }

  public double getSpent() {
    return spent;
  }

  public void setSpent(double spent) {
    this.spent = spent;
  }
}
