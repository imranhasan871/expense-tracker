package com.expense.handlers;

import com.expense.config.DatabaseConfig;
import com.expense.models.Budget;
import com.expense.repository.BudgetRepository;
import com.expense.repository.ExpenseRepository;
import com.sun.net.httpserver.HttpExchange;
import java.io.IOException;
import java.sql.Connection;
import java.util.HashMap;
import java.util.Map;
import java.util.List;
import java.io.InputStreamReader;
import java.net.URI;
import java.io.BufferedReader;

public class BudgetHandler extends BaseHandler {

  @Override
  public void handle(HttpExchange exchange) throws IOException {
    String path = exchange.getRequestURI().getPath();
    String method = exchange.getRequestMethod();

    // Router logic inside handler (simple standard pattern)
    if (path.startsWith("/api/budgets/status")) {
      handleStatus(exchange);
      return;
    }
    if (path.endsWith("/lock") && "POST".equalsIgnoreCase(method)) {
      handleLock(exchange);
      return;
    }
    if (path.equals("/api/monitoring") && "GET".equalsIgnoreCase(method)) {
      handleMonitoring(exchange);
      return;
    }
    if (path.equals("/api/budgets")) {
      if ("GET".equalsIgnoreCase(method))
        handleGet(exchange);
      else if ("POST".equalsIgnoreCase(method))
        handlePost(exchange);
      else
        sendError(exchange, 405, "Method Not Allowed");
      return;
    }

    sendError(exchange, 404, "Not Found");
  }

  private void handleGet(HttpExchange exchange) throws IOException {
    String query = exchange.getRequestURI().getQuery();
    int year = 2026; // parse from query
    if (query != null && query.contains("year=")) {
      // quick parsing hack, robust parsing skipped for brevity
      year = Integer.parseInt(query.split("year=")[1].split("&")[0]);
    }

    try (Connection conn = DatabaseConfig.getConnection()) {
      BudgetRepository repo = new BudgetRepository(conn);
      List<Budget> budgets = repo.findAll(year);
      double totalAnnual = repo.getTotalAnnualBudget(year);

      // Construct Summary structure needed by frontend
      Map<String, Object> summary = new HashMap<>();
      summary.put("TotalAnnualBudget", totalAnnual);
      // simplified summary for now
      summary.put("HighestAllocation", 0);
      summary.put("SavingsTarget", 0);
      summary.put("RemainingBudget", 0);

      Map<String, Object> resp = new HashMap<>();
      resp.put("budgets", budgets);
      resp.put("summary", summary);

      sendResponse(exchange, 200, new SuccessResponse(resp));
    } catch (Exception e) {
      e.printStackTrace();
      sendError(exchange, 500, e.getMessage());
    }
  }

  private void handlePost(HttpExchange exchange) throws IOException {
    try (Connection conn = DatabaseConfig.getConnection()) {
      Budget req = gson.fromJson(new InputStreamReader(exchange.getRequestBody()), Budget.class);
      BudgetRepository repo = new BudgetRepository(conn);
      repo.createOrUpdate(req);
      sendResponse(exchange, 200, new SuccessResponse(req, "Budget set successfully"));
    } catch (Exception e) {
      e.printStackTrace();
      sendError(exchange, 500, e.getMessage());
    }
  }

  private void handleStatus(HttpExchange exchange) throws IOException {
    String query = exchange.getRequestURI().getQuery();
    Map<String, String> queryMap = parseQuery(query);
    int categoryId = Integer.parseInt(queryMap.get("category_id"));
    int year = Integer.parseInt(queryMap.get("year"));

    try (Connection conn = DatabaseConfig.getConnection()) {
      BudgetRepository budgetRepo = new BudgetRepository(conn);
      ExpenseRepository expenseRepo = new ExpenseRepository(conn);

      Budget b = budgetRepo.getByCategory(categoryId, year);
      double spent = expenseRepo.getYearlyTotal(categoryId, year);

      Map<String, Object> data = new HashMap<>();
      data.put("spent", spent);

      if (b != null) {
        data.put("allocated", b.getAmount());
        data.put("remaining", b.getAmount() - spent);
        double percent = (b.getAmount() > 0) ? (spent / b.getAmount() * 100) : 0;
        data.put("percent", percent);
        data.put("is_locked", b.isIs_locked());
      } else {
        data.put("allocated", 0);
        data.put("remaining", 0);
        data.put("percent", 0);
        data.put("is_locked", false);
      }

      sendResponse(exchange, 200, new SuccessResponse(data));
    } catch (Exception e) {
      e.printStackTrace();
      sendError(exchange, 500, e.getMessage());
    }
  }

  private void handleMonitoring(HttpExchange exchange) throws IOException {
    String query = exchange.getRequestURI().getQuery();
    int year = 2026;
    if (query != null && query.contains("year=")) {
      year = Integer.parseInt(query.split("year=")[1].split("&")[0]);
    }

    try (Connection conn = DatabaseConfig.getConnection()) {
      BudgetRepository repo = new BudgetRepository(conn);
      List<Budget> data = repo.getMonitoringData(year);
      sendResponse(exchange, 200, new SuccessResponse(data));
    } catch (Exception e) {
      e.printStackTrace();
      sendError(exchange, 500, "DB Error");
    }
  }

  private void handleLock(HttpExchange exchange) throws IOException {
    // Path: /api/budgets/{id}/lock
    String path = exchange.getRequestURI().getPath();
    String[] parts = path.split("/");
    int id = Integer.parseInt(parts[3]);

    try (Connection conn = DatabaseConfig.getConnection()) {
      // Read body {"is_locked": true}
      Budget req = gson.fromJson(new InputStreamReader(exchange.getRequestBody()), Budget.class);

      BudgetRepository repo = new BudgetRepository(conn);
      repo.toggleLock(id, req.isIs_locked());

      sendResponse(exchange, 200, new SuccessResponse(null, "Lock updated"));
    } catch (Exception e) {
      e.printStackTrace();
      sendError(exchange, 500, e.getMessage());
    }
  }

  private Map<String, String> parseQuery(String query) {
    Map<String, String> res = new HashMap<>();
    if (query == null)
      return res;
    for (String param : query.split("&")) {
      String[] entry = param.split("=");
      if (entry.length > 1)
        res.put(entry[0], entry[1]);
    }
    return res;
  }
}
