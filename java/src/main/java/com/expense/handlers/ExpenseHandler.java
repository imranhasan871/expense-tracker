package com.expense.handlers;

import com.expense.config.DatabaseConfig;
import com.expense.models.Expense;
import com.expense.models.Budget;
import com.expense.repository.ExpenseRepository;
import com.expense.repository.BudgetRepository;
import com.sun.net.httpserver.HttpExchange;
import java.io.IOException;
import java.sql.Connection;
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import java.io.InputStreamReader;
import java.time.LocalDate;

public class ExpenseHandler extends BaseHandler {

  @Override
  public void handle(HttpExchange exchange) throws IOException {
    String method = exchange.getRequestMethod();

    try (Connection conn = DatabaseConfig.getConnection()) {
      ExpenseRepository repo = new ExpenseRepository(conn);
      String path = exchange.getRequestURI().getPath();

      // DELETE /api/expenses/{id}
      if ("DELETE".equalsIgnoreCase(method)) {
        String[] parts = path.split("/");
        if (parts.length == 4) {
          int id = Integer.parseInt(parts[3]);
          repo.delete(id);
          sendResponse(exchange, 200, new SuccessResponse(null, "Expense deleted"));
          return;
        }
      }

      if ("GET".equalsIgnoreCase(method)) {
        String query = exchange.getRequestURI().getQuery();
        // Parse filters
        Map<String, String> q = parseQuery(query);

        String search = q.get("search");
        String startDate = q.get("start_date");
        String endDate = q.get("end_date");
        Integer catId = q.containsKey("category_id") ? Integer.parseInt(q.get("category_id")) : null;
        Double min = q.containsKey("min_amount") && !q.get("min_amount").isEmpty()
            ? Double.parseDouble(q.get("min_amount"))
            : null;
        Double max = q.containsKey("max_amount") && !q.get("max_amount").isEmpty()
            ? Double.parseDouble(q.get("max_amount"))
            : null;

        List<Expense> expenses = repo.findAll(search, startDate, endDate, catId, min, max);
        sendResponse(exchange, 200, new SuccessResponse(expenses));
      } else if ("POST".equalsIgnoreCase(method)) {
        Expense req = gson.fromJson(new InputStreamReader(exchange.getRequestBody()), Expense.class);

        // CIRCUIT BREAKER CHECK
        BudgetRepository budgetRepo = new BudgetRepository(conn);
        // Assume date is in YYYY-MM-DD
        int year = LocalDate.parse(req.getExpense_date()).getYear();
        Budget budget = budgetRepo.getByCategory(req.getCategory_id(), year);

        if (budget != null && budget.isIs_locked()) {
          sendError(exchange, 403, "Budget for this category is LOCKED. Cannot add expense.");
          return;
        }

        repo.create(req);
        sendResponse(exchange, 201, new SuccessResponse(null, "Expense recorded"));
      } else {
        sendError(exchange, 405, "Method Not Allowed");
      }
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
